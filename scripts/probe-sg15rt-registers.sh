#!/usr/bin/env bash

# Read-only probe for SG15RT FC04 input registers.
# Each reference is read separately so an unsupported address does not hide
# valid values elsewhere in the range.
set -uo pipefail

HOST="${1:-192.168.178.24}"
SLAVE="${SLAVE:-1}"
MBPOLL="${MBPOLL:-mbpoll}"
OUT="${OUT:-sg15rt-register-probe-$(date +%Y%m%d-%H%M%S).log}"

if ! command -v "$MBPOLL" >/dev/null 2>&1; then
  printf 'mbpoll not found (looked for %s)\n' "$MBPOLL" >&2
  exit 127
fi

{
  printf '# SG15RT input-register probe\n'
  printf '# host=%s unit_id=%s function=FC04 start=4990 end=5150\n' "$HOST" "$SLAVE"
  printf '# started=%s\n' "$(date -Is)"
  printf '# Each read is one register; no Modbus writes are sent.\n\n'
} | tee "$OUT"

failed=0
for ref in $(seq 4990 5150); do
  printf '===== reference %s =====\n' "$ref" | tee -a "$OUT"
  if ! "$MBPOLL" -m tcp -a "$SLAVE" -t 3:hex -r "$ref" -c 1 -1 -o 1 "$HOST" 2>&1 | tee -a "$OUT"; then
    failed=1
  fi
done

# Alternative references appearing in published protocol/map versions.
for ref in 13001 13002 13007 13008; do
  printf '===== alternate reference %s =====\n' "$ref" | tee -a "$OUT"
  if ! "$MBPOLL" -m tcp -a "$SLAVE" -t 3:hex -r "$ref" -c 1 -1 -o 1 "$HOST" 2>&1 | tee -a "$OUT"; then
    failed=1
  fi
done

printf '\n# finished=%s\n# log=%s\n' "$(date -Is)" "$OUT" | tee -a "$OUT"
exit "$failed"
