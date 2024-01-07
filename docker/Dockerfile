FROM alpine
LABEL MAINTAINER="stefan.roeck@gmail.com"

ENV SUNGROW_HOST=
ENV SUNGROW_PORT=8082
ENV SUNGROW_DATA="pv,battery"
ENV MQTT_URL=
ENV MQTT_TOPIC="sungrow/stats"
ENV MQTT_CLIENTID=
ENV MQTT_USER=
ENV MQTT_PASSWORD=
ENV SLEEP=10

RUN mkdir /sungrow
COPY ./bin/sungrow-go-docker /sungrow
WORKDIR /sungrow

CMD /sungrow/sungrow-go-docker -host ${SUNGROW_HOST} -port ${SUNGROW_PORT} -data ${SUNGROW_DATA} -mqtt.server ${MQTT_URL} -sleep ${SLEEP} -mqtt.topic ${MQTT_TOPIC} -mqtt.user ${MQTT_USER} -mqtt.password ${MQTT_PASSWORD} -mqtt.clientId ${MQTT_CLIENTID}
