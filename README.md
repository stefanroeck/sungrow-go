# stefanroeck/sungrow

Tiny application to access real-time data from Sungrow inverters with WiNet-S dongle over WebSocket and publish data to a mqtt server.

Written in GoLang, based on https://github.com/nItroTools/sungrow-go.

# Usage

## Docker

- Pull image from DockerHub manually or via Docker-Compose file: https://hub.docker.com/r/stefanroeck/sungrow

### Build and run this repository with Docker Compose

The compose file builds the application from the repository and reads its configuration from `docker/.env`. Set the inverter and MQTT values in that file, then run these commands from the repository root:

```bash
docker compose -f docker/docker-compose.yml up --build
```

**docker-compose.yml**:

```yaml
version: "3.9"
services:
  sungrow-go:
    image: stefanroeck/sungrow:latest
    restart: unless-stopped
    environment:
      - SUNGROW_PROTOCOL=wss
      - SUNGROW_HOST=${SUNGROW_HOST}
      - SUNGROW_PORT=443
      - SUNGROW_USER=admin
      - SUNGROW_PASSWORD=secret
      - MQTT_URL=${MQTT_URL}
      - MQTT_USER=${MQTT_USER}
      - MQTT_PASSWORD=${MQTT_PASSWORD}
      - MQTT_TOPIC="sungrow/stats"
      - MQTT_CLIENTID=sungrow
      - MQTT_SKIP_SSL_VERIFY=false
      - SLEEP=30
```

## Install & Build locally

- Checkout Repo

```bash
go install .
go build -o ./bin/
```

### Build for Raspberry Pi

```bash
env GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/sungrow-go-raspi
```

### Build for Docker

```bash
env CGO_ENABLED=0 GOOS=linux go build -o ./bin/sungrow-go-docker
```

## Run locally

List available and required parameters

```bash
$ sungrow-go -help
```

Basic usage with ip address of your inverter (e.g. `192.168.2.100`)

```bash
$ sungrow-go -protocol wss -host 192.168.2.100 -user admin -password secret -mqtt.server mqtt://test.mosquitto.org:1883 -mqtt.topic honk/demo
```
For newer firmware versions (around 2024/2025): 
- Authentication using user/password is required
- Protocol needs to be wss instead of ws (TLS)

# MQTT

Sample Message:

```json
{
  "activePower": 0.99,
  "arrayInsulationResistance": 3000,
  "busVoltage": 693.6,
  "currentPhaseA": 1.5,
  "currentPhaseB": 1.4,
  "currentPhaseC": 1.4,
  "gridFrequency": 49.97,
  "inverterTemp": 34.7,
  "status": "running",
  "sunPower": 0.99,
  "todayEnergy": 2.2,
  "totalEnergy": 6584.9,
  "totalRunningTime": 1948,
  "voltagePhaseA": 232,
  "voltagePhaseB": 235.2,
  "voltagePhaseC": 236
}
```

## SG15RT Modbus register reference

The following read-only FC04 (`Read Input Registers`) mappings were checked against an SG15RT using a WiNet-S and compared with the existing web-interface payload. Addresses in the `mbpoll` column are the references used for successful reads. For a Modbus library that takes zero-based PDU addresses, use `mbpoll reference - 1`; keep this conversion in one place and verify the library's addressing convention.

The protocol sources can differ by version and sometimes show one-based references. The device identity read returned `0x2435` at reference `5000`, identifying the SG15RT. Treat the observed references below as the device-verified baseline for this unit and firmware.

| MQTT key | FC04 `mbpoll` reference(s) | Type | Raw scale / unit | MQTT/web value conversion | Notes |
|---|---:|---|---|---|---|
| `todayEnergy` | `5003` | U16 | 0.1 kWh per count | raw ÷ 10 → kWh | Observed raw `20` → `2.0 kWh`. |
| `totalEnergy` | `5004–5005` | U32 | 1 kWh per count | raw → kWh | Read both registers together. Observed register result was `56555 kWh`; web payload showed `56555.9 kWh`, so the register read appears to lack the web value's fractional precision. |
| `totalRunningTime` | `5006–5007` | U32 | 1 h per count | raw → h | Read both registers together; observed `14122 h`, matching the web payload. |
| `inverterTemp` | `5008` | S16 | 0.1 °C per count | signed raw ÷ 10 → °C | Observed `36.1 °C`, matching the web payload. |
| `sunPower` | `5017–5018` | U32 | 1 W per count | raw W ÷ 1000 → kW; web payload displays 2 decimals | This is the protocol's total DC power. Observed `2053 W` → `2.053 kW`, displayed as `2.05`, matching the web payload. |
| `voltagePhaseA` | `5019` | U16 | 0.1 V per count | raw ÷ 10 → V | Observed `230.4 V`, matching. |
| `voltagePhaseB` | `5020` | U16 | 0.1 V per count | raw ÷ 10 → V | Observed `234.4 V`; web payload showed `233.6 V` (0.8 V difference to recheck). |
| `voltagePhaseC` | `5021` | U16 | 0.1 V per count | raw ÷ 10 → V | Observed `234.4 V`, matching. |
| `currentPhaseA` | `5022` | U16 | 0.1 A per count | raw ÷ 10 → A | Observed `2.6 A`, matching. |
| `currentPhaseB` | `5023` | U16 | 0.1 A per count | raw ÷ 10 → A | Observed `2.6 A`, matching. |
| `currentPhaseC` | `5024` | U16 | 0.1 A per count | raw ÷ 10 → A | Observed `2.6 A`, matching. |
| `activePower` | `5031–5032` | U32 | 1 W per count | raw W ÷ 1000 → kW; web payload displays 2 decimals | Read both registers together. Observed `2051 W` → `2.051 kW`, displayed as `2.05`, matching the web payload. |
| `gridFrequency` | `5148` | U16 | 0.01 Hz per count | raw ÷ 100 → Hz | Higher precision read: `50.00 Hz`; web payload showed `49.99 Hz`. Reference `5036` is also readable at 0.1 Hz/count and returned `50.0 Hz`. |
| `status` | `5038` | U16 enum | Integer state code | Map code to a stable application string | Code `0` means normal grid-connected run; the web payload calls this `"on grid operation"`. Do not treat the English wording as a register value. |
| `arrayInsulationResistance` | `5071` | U16 | 1 kΩ per count | raw → kΩ | Observed `2150 kΩ`, matching. |
| `busVoltage` | `5147` | U16 | 0.1 V per count | raw ÷ 10 → V | Observed `689.8 V`; web payload showed `689.6 V` (0.2 V difference to recheck). |

Useful identification/format field: FC04 reference `5000` is U16 device type code (`0x2435` for this SG15RT); reference `5002` is output type (`1` in the observed unit, three-phase four-wire). Other MPPT voltage/current registers are `5011–5016`, each U16 at 0.1 V or 0.1 A per count. The observed third MPPT values were zero; confirm that the inverter configuration uses those inputs before interpreting zero as a fault or unavailable value.

### 32-bit decoding, units, and rounding

- Sungrow encodes U32/S32 double-word values **low word first**. For words `low` then `high`, reconstruct the unsigned integer as `(high << 16) | low`. Read each pair in one FC04 request (`count=2`); isolated reads of the second word returned `Illegal data address` in our probe.
- U16/S16 values are big-endian within the register. Apply signed interpretation for S16 temperature if the raw value has its sign bit set.
- Power registers `5017–5018` (DC) and `5031–5032` (AC active) are in **watts**. The existing MQTT/web payload expresses these power fields in **kilowatts**, so divide by 1000 before publishing. The observed `2053 W` and `2051 W` become `2.053 kW` and `2.051 kW`; the web payload rounded both to two decimal places (`2.05`). Preserve full precision internally and round only at the output boundary if matching that display.
- Energy units differ by register: daily energy is 0.1 kWh/count, while the verified total-energy pair is whole kWh/count. Do not apply one scale to both. The total-energy register's observed whole-kWh precision does not reproduce the web payload's decimal digit.
- Voltage, current, temperature, and frequency scales are shown per field above. The higher-precision frequency register `5148` is in 0.01 Hz/count; register `5036` is 0.1 Hz/count.

These mappings are read-only and specific to the tested SG15RT/WiNet-S setup. Validate them against the exact firmware and compare readings under stable conditions before using them in production. Public protocol reference: [Sungrow residential/commercial PV grid-connected inverter communication protocol, V1.1.53](https://prep.shop.kpenergy.fi/globalassets/product-sheets/inverters/sungrow/130675/ti_20230117_communication-protocol-of--residential-and-commerical-pv-grid-connected-inverter_v1.1.53_en.pdf).

# Supported inverters

Tested Sungrow inverters with WiNet-S dongle:

- SH10RT (by https://github.com/nItroTools/sungrow-go)
- SG15RT
