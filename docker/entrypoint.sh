#!/bin/sh

exec /sungrow/sungrow-go-docker \
  -protocol="$SUNGROW_PROTOCOL" \
  -host="$SUNGROW_HOST" \
  -port="$SUNGROW_PORT" \
  -user="$SUNGROW_USER" \
  -password="$SUNGROW_PASSWORD" \
  -mqtt.server="$MQTT_URL" \
  -sleep="$SLEEP" \
  -mqtt.topic="$MQTT_TOPIC" \
  -mqtt.user="$MQTT_USER" \
  -mqtt.password="$MQTT_PASSWORD" \
  -mqtt.clientId="$MQTT_CLIENTID" \
  -mqtt.skipSSLVerify="$MQTT_SKIP_SSL_VERIFY"
