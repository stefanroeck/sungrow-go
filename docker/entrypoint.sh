#!/bin/sh

exec /sungrow/sungrow-go-docker \
  -transport="${SUNGROW_TRANSPORT:-websocket}" \
  -protocol="$SUNGROW_PROTOCOL" \
  -host="$SUNGROW_HOST" \
  -port="$SUNGROW_PORT" \
  -user="$SUNGROW_USER" \
  -password="$SUNGROW_PASSWORD" \
  -modbus.port="${SUNGROW_MODBUS_PORT:-502}" \
  -modbus.unitId="${SUNGROW_MODBUS_UNIT_ID:-1}" \
  -mqtt.server="$MQTT_URL" \
  -sleep="$SLEEP" \
  -mqtt.topic="$MQTT_TOPIC" \
  -mqtt.user="$MQTT_USER" \
  -mqtt.password="$MQTT_PASSWORD" \
  -mqtt.clientId="$MQTT_CLIENTID" \
  -mqtt.skipSSLVerify="$MQTT_SKIP_SSL_VERIFY"
