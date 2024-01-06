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
    image: stefanroeck/sungrow
    restart: unless-stopped
    environment:
      - SUNGROW_HOST=${SUNGROW_HOST}
      - MQTT_URL=${MQTT_URL}
      - MQTT_USER=${MQTT_USER}
      - MQTT_PASSWORD=${MQTT_PASSWORD}
      - SLEEP=30
```

## Install & Build locally

- Checkout Repo

```bash
$ go install .
$ go build ./bin
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
$ sungrow-go -host 192.168.2.100 -mqtt.server mqtt://test.mosquitto.org:1883 -mqtt.topic honk/demo
```

# MQTT

Sample Message:

```json
{
  "activePower": 0,
  "busVoltage": 0,
  "inverterTemp": 28.5,
  "status": "standby",
  "sunPower": 0,
  "todayEnergy": 2.9,
  "totalEnergy": 6582.7,
  "totalRunningTime": 1945
}
```

# Supported inverters

Tested Sungrow inverters with WiNet-S dongle:

- SH10RT (by https://github.com/nItroTools/sungrow-go)
- SG15RT
