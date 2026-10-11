#!/usr/bin/env bash

# Read confirmed SG15RT FC04 values and print a comparable JSON object.
set -uo pipefail

HOST="${1:-192.168.178.24}"
SLAVE="${SLAVE:-1}"
MBPOLL="${MBPOLL:-mbpoll}"
OUT="${OUT:-sg15rt-confirmed-values-$(date +%Y%m%d-%H%M%S).log}"
VALUES_TSV="$(mktemp "${TMPDIR:-/tmp}/sg15rt-values.XXXXXX")"
trap 'rm -f "$VALUES_TSV"' EXIT

if ! command -v "$MBPOLL" >/dev/null 2>&1; then
  printf 'mbpoll not found (looked for %s)\n' "$MBPOLL" >&2
  exit 127
fi

failed=0

record_value() {
  local key="$1" raw="$2" scale="$3" unit="$4" formatted
  if (( scale == 1 )); then
    formatted="$raw"
  else
    formatted=$(awk -v raw="$raw" -v scale="$scale" 'BEGIN { value = sprintf("%.2f", raw / scale); sub(/0+$/, "", value); sub(/\\.$/, "", value); print value }')
  fi
  printf '%s\t%s\t%s\n' "$key" "$formatted" "$unit" >> "$VALUES_TSV"
}

read_u16() {
  local ref="$1" key="$2" scale="$3" unit="$4" output raw value
  printf '\n## %s (FC04 reference %s; raw / %s %s)\n' "$key" "$ref" "$scale" "$unit" | tee -a "$OUT"
  if ! output=$("$MBPOLL" -m tcp -a "$SLAVE" -t 3:hex -r "$ref" -c 1 -1 -o 1 "$HOST" 2>&1); then
    printf '%s\n' "$output" | tee -a "$OUT"
    failed=1
    return 0
  fi
  printf '%s\n' "$output" | tee -a "$OUT"
  raw=$(printf '%s\n' "$output" | sed -nE "s/.*\\[$ref\\]:[[:space:]]*0x([[:xdigit:]]+).*/\\1/p" | tail -n 1)
  if [[ ! "$raw" =~ ^[[:xdigit:]]{1,4}$ ]]; then
    printf 'Could not parse register %s for key %s\n' "$ref" "$key" | tee -a "$OUT"
    failed=1
    return 0
  fi
  value=$((16#$raw))
  record_value "$key" "$value" "$scale" "$unit"
}

read_u32() {
  local ref="$1" key="$2" scale="$3" unit="$4" output low high value
  printf '\n## %s (FC04 references %s-%s; low word first)\n' "$key" "$ref" "$((ref + 1))" | tee -a "$OUT"
  # Sungrow specifies little-endian word order for U32. Both words are
  # requested together. Do not use mbpoll -B, which reverses word order.
  if ! output=$("$MBPOLL" -m tcp -a "$SLAVE" -t 3:hex -r "$ref" -c 2 -1 -o 1 "$HOST" 2>&1); then
    printf '%s\n' "$output" | tee -a "$OUT"
    failed=1
    return 0
  fi
  printf '%s\n' "$output" | tee -a "$OUT"
  low=$(printf '%s\n' "$output" | sed -nE "s/.*\\[$ref\\]:[[:space:]]*0x([[:xdigit:]]+).*/\\1/p" | tail -n 1)
  high=$(printf '%s\n' "$output" | sed -nE "s/.*\\[$((ref + 1))\\]:[[:space:]]*0x([[:xdigit:]]+).*/\\1/p" | tail -n 1)
  if [[ ! "$low" =~ ^[[:xdigit:]]{1,4}$ || ! "$high" =~ ^[[:xdigit:]]{1,4}$ ]]; then
    printf 'Could not parse both words at references %s-%s for key %s\n' "$ref" "$((ref + 1))" "$key" | tee -a "$OUT"
    failed=1
    return 0
  fi
  low=$((16#$low))
  high=$((16#$high))
  value=$((high * 65536 + low))
  record_value "$key" "$value" "$scale" "$unit"
}

{
  printf '# Confirmed SG15RT input-register reads\n'
  printf '# host=%s unit_id=%s function=FC04\n' "$HOST" "$SLAVE"
  printf '# started=%s\n' "$(date '+%Y-%m-%dT%H:%M:%S%z')"
  printf '# Read-only; U32 values are requested as two consecutive registers.\n'
} | tee "$OUT"

run_read() {
  if ! "$@"; then
    failed=1
  fi
}

# Device identity and output type validate model/addressing and phase format.
run_read read_u16 5000 'deviceTypeCode' 1 'code'
run_read read_u16 5002 'outputType' 1 'enum'

# Fields supported by the SG-series protocol and observed probe responses.
run_read read_u16 5003 'todayEnergy' 10 'kWh'
run_read read_u32 5004 'totalEnergy' 1 'kWh'
run_read read_u32 5006 'totalRunningTime' 1 'h'
run_read read_u16 5008 'inverterTemp' 10 'degC'
run_read read_u16 5011 'mppt1Voltage' 10 'V'
run_read read_u16 5012 'mppt1Current' 10 'A'
run_read read_u16 5013 'mppt2Voltage' 10 'V'
run_read read_u16 5014 'mppt2Current' 10 'A'
run_read read_u16 5015 'mppt3Voltage' 10 'V'
run_read read_u16 5016 'mppt3Current' 10 'A'
run_read read_u32 5017 'totalDCPower' 1 'W'
run_read read_u16 5019 'voltagePhaseA' 10 'V'
run_read read_u16 5020 'voltagePhaseB' 10 'V'
run_read read_u16 5021 'voltagePhaseC' 10 'V'
run_read read_u16 5022 'currentPhaseA' 10 'A'
run_read read_u16 5023 'currentPhaseB' 10 'A'
run_read read_u16 5024 'currentPhaseC' 10 'A'
run_read read_u32 5031 'activePower' 1 'W'
run_read read_u16 5036 'gridFrequency' 10 'Hz'
run_read  read_u16 5038 'workState' 1 'enum'
run_read read_u16 5071 'arrayInsulationResistance' 1 'kOhm'
run_read read_u16 5147 'busVoltage' 10 'V'
run_read read_u16 5148 'gridFrequencyHighPrecision' 100 'Hz'

printf '\n## Comparison values (JSON)\n' | tee -a "$OUT"
printf '{\n' | tee -a "$OUT"
value_count=$(wc -l < "$VALUES_TSV" | tr -d ' ')
line_number=0
while IFS=$'\t' read -r key value unit; do
  [[ -n "$key" ]] || continue
  line_number=$((line_number + 1))
  comma=,
  (( line_number == value_count )) && comma=
  printf '  "%s": %s%s\n' "$key" "$value" "$comma" | tee -a "$OUT"
done < "$VALUES_TSV"
printf '}\n' | tee -a "$OUT"

printf '\n# finished=%s\n# log=%s\n' "$(date '+%Y-%m-%dT%H:%M:%S%z')" "$OUT" | tee -a "$OUT"
exit "$failed"
