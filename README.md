# sungrow-go

GoLang implementation for accessing real-time data from Sungrow inverters with WiNet-S dongle using WebSocket.
Sends data to a mqtt server for further processing.

## Install & Build

```bash
$ go install .
$ go build .
```

## Usage

List available and required parameters

```bash
$ sungrow-go -help
```

Basic usage with ip address of your inverter (e.g. `192.168.2.100`)

```bash
$ sungrow-go -ip 192.168.2.100 -mqtt.server mqtt://test.mosquitto.org:1883 -mqtt.topic honk/demo
```

Sample Message:

```json
{
  "activePower": 0,
  "inverterTemp": 35.3,
  "sunPower": 0,
  "todayEnergy": 13.9,
  "totalEnergy": 5850.2,
  "totalRunningTime": 1349
}
```

## Supported inverters

Tested Sungrow inverters with WiNet-S dongle:

- SH10RT (by https://github.com/nItroTools/sungrow-go)
- SG15RT
