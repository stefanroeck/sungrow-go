# stefanroeck/sungrow

Tiny application to access real-time data from Sungrow inverters with WiNet-S dongle over WebSocket and publish data to a mqtt server.

Written in GoLang, based on https://github.com/nItroTools/sungrow-go.

# Usage

## Docker

- Pull image from DockerHub manually or via Docker-Compose file: https://hub.docker.com/r/stefanroeck/sungrow

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
$ go install .
$ go build -o ./bin/
```

### Build for Rasperry PI

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

# Supported inverters

Tested Sungrow inverters with WiNet-S dongle:

- SH10RT (by https://github.com/nItroTools/sungrow-go)
- SG15RT
