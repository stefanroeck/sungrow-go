# ---------- Build stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/sungrow-go-docker

# ---------- Final stage ----------

FROM alpine
LABEL MAINTAINER="stefan.roeck@gmail.com"

ENV SUNGROW_PROTOCOL=wss
ENV SUNGROW_HOST=
ENV SUNGROW_PORT=443
ENV SUNGROW_USER=
ENV SUNGROW_PASSWORD=
ENV MQTT_URL=
ENV MQTT_TOPIC="sungrow/stats"
ENV MQTT_CLIENTID=
ENV MQTT_USER=
ENV MQTT_PASSWORD=
ENV SLEEP=10

RUN mkdir /sungrow
COPY --from=builder /app/bin/sungrow-go-docker /sungrow/sungrow-go-docker
WORKDIR /sungrow

# Define entrypoint
CMD /sungrow/sungrow-go-docker \
    -protocol ${SUNGROW_PROTOCOL} \
    -host ${SUNGROW_HOST} \
    -port ${SUNGROW_PORT} \
    -user=${SUNGROW_USER} \
    -password=${SUNGROW_PASSWORD} \
    -mqtt.server ${MQTT_URL} \
    -sleep ${SLEEP} \
    -mqtt.topic ${MQTT_TOPIC} \
    -mqtt.user ${MQTT_USER} \
    -mqtt.password ${MQTT_PASSWORD} \
    -mqtt.clientId ${MQTT_CLIENTID}
