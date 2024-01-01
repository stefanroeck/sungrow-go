package main

import (
	"flag"
	"fmt"
	"log"
	"maps"
	"strings"
	"time"

	"github.com/sroeck/sungrow-go/mqtt"
	"github.com/sroeck/sungrow-go/ws"
)

type InverterParams struct {
	host      string
	port      int
	path      string
	data      string
	types     []string
}

func main() {
	inv, mqttParams, sleepTimeInSeconds := flags()

	ticker := time.NewTicker(time.Duration(sleepTimeInSeconds) * time.Second)
    done := make(chan bool)
	defer ticker.Stop()

	go func() {
		// run immediately and then at the configured interval
		fetchDataFromInverterAndSendToMqtt(inv, mqttParams) 
		for {
			select {
            case <-done:
                return
            case <-ticker.C:
				fetchDataFromInverterAndSendToMqtt(inv, mqttParams)
            }
		}
	}()
	
	// block forever
	select {}
}

func fetchDataFromInverterAndSendToMqtt(inverterParams InverterParams, mqttParams mqtt.MqttParams) {
	webSocket := openWebSocket(inverterParams)	
	defer webSocket.Close()

	var receivedValues map[string]float64 = make(map[string]float64)

	for _, t := range inverterParams.types {
		switch t {
		case "pv":
			err, pvValues := webSocket.Pv(pvKeys)
			if err != nil {
				log.Fatalln(err)
			}
			maps.Copy(receivedValues, pvValues)
		case "battery":
			err, batteryValues := webSocket.Battery(batteryKeys)
			if err != nil {
				log.Fatalln(err)
			}
			maps.Copy(receivedValues, batteryValues)
		}
	}

	fmt.Println("Received the following values:")
	for k, v := range receivedValues {
		fmt.Println(k, "=", v)
	}
	mqtt.Send(mqttParams, receivedValues)
}

func openWebSocket(inverterParams InverterParams) (*ws.WS) {
	webSocket := ws.NewWS(inverterParams.host, inverterParams.port, inverterParams.path)
	if err := webSocket.Connect(); err != nil {
		log.Fatalln(err)
	}
	return webSocket
}

// flags defines, parses and validates command-line flags from os.Args[1:]
func flags() (InverterParams, mqtt.MqttParams, int) {

	host := flag.String("host", "", "Hostname/IP address of the Sungrow inverter")
	port := flag.Int("port", 8082, "WebSocket port of the Sungrow inverter")
	path := flag.String("path", "/ws/home/overview", "Server path from where data is requested")
	data := flag.String("data", "pv,battery", "Select the data to be requested comma separated.\nPossible values are \"pv\" and \"battery\"")
	mqttServer := flag.String("mqtt.server", "", "mqtt server incl. protocol, e.g. mqtt://localhost:1883. For TLS use ssl scheme, e.g. ssl://localhost:8883")
	mqttUser := flag.String("mqtt.user", "", "mqtt user name")
	mqttPassword := flag.String("mqtt.password", "", "mqtt password")
	mqttClientId := flag.String("mqtt.clientId", "", "mqtt clientId that is used for publishing")
	mqttTopic := flag.String("mqtt.topic", "topic", "mqtt topic to which the data are published")
	sleepBetweenCalls := flag.Int("sleep", 10, "sleep time in seconds between inverter calls.")
	flag.Parse()

	inv := &InverterParams{host: *host, port: *port, path: *path, data: *data}
	
	mqttParams := &mqtt.MqttParams{Server: *mqttServer, ClientId: *mqttClientId, Topic: *mqttTopic, User: *mqttUser, Password: *mqttPassword}

	// Validate flags
	flagsValidate(inv)
	flagsValidateMqtt(*mqttParams)

	return *inv, *mqttParams, *sleepBetweenCalls
}

// flagsValidate validates all flags
func flagsValidate(inv *InverterParams) {
	if len(inv.host) == 0 {
		log.Fatalln("Required parameter 'host' not set!\n'sungrow-go -help' lists available parameters.")
	}

	inv.types = strings.Split(inv.data, ",")
	if len(inv.types) < 1 {
		log.Fatalln("Required parameter 'data' not set or invalid value!\n'sungrow-go -help' lists available parameters and values.")
	}
	for _, t := range inv.types {
		switch t {
		case "pv":
		case "battery":
			break
		default:
			log.Fatalf("Invalid value \"%s\" for parameter 'data'!\n'sungrow-go -help' lists available parameters and values.\n", t)
		}
	}
}

func flagsValidateMqtt(params mqtt.MqttParams) {
	if (params.Server == "") {
		log.Fatalln("Missing parameter mqtt.server")
	}
}

