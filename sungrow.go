package main

import (
	"encoding/json"
	"flag"
	"log"
	"maps"
	"strings"
	"time"

	"github.com/sroeck/sungrow-go/mqtt"
	"github.com/sroeck/sungrow-go/ws"
)

type InverterParams struct {
	host  string
	port  int
	path  string
	data  string
	types []string
}

func main() {
	inverterParams, mqttParams, sleepTimeInSeconds := flags()
	log.Printf("Polling the inverter every %d seconds\n", sleepTimeInSeconds)

	ticker := time.NewTicker(time.Duration(sleepTimeInSeconds) * time.Second)
	done := make(chan bool)
	defer ticker.Stop()

	go func() {
		// run immediately and then at the configured interval
		fetchDataFromInverterAndSendToMqtt(inverterParams, mqttParams)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fetchDataFromInverterAndSendToMqtt(inverterParams, mqttParams)
			}
		}
	}()

	// block forever
	select {}
}

func fetchDataFromInverterAndSendToMqtt(inverterParams *InverterParams, mqttParams *mqtt.MqttParams) {
	webSocket := openWebSocket(inverterParams)
	defer webSocket.Close()

	var receivedValues map[string]any = make(map[string]any)

	for _, t := range inverterParams.types {
		switch t {
		case "pv":
			log.Println("Fetching pv data")
			pvValues, err := webSocket.Pv()
			processWsResult(receivedValues, pvValues, err)
		case "battery":
			log.Println("Fetching battery data")
			batteryValues, err := webSocket.Battery()
			processWsResult(receivedValues, batteryValues, err)
		}
	}

	if len(receivedValues) == 0 {
		log.Println("Skip sending MQTT data as no data have been returned from inverter")
		return
	}

	mqtt.Send(mqttParams, receivedValues)
}

func processWsResult(targetMap map[string]any, resultMap map[string]any, err error) {
	if err != nil {
		log.Fatalln(err)
		return
	}
	if len(resultMap) == 0 {
		log.Println("Warning: No data returned")
		return
	}

	prettyPrintMap(resultMap)
	maps.Copy(targetMap, resultMap)
}

func prettyPrintMap(resultMap map[string]any) {
	b, _ := json.MarshalIndent(resultMap, "", "  ")
	log.Print("Received the following values: ", string(b))
}

func openWebSocket(inverterParams *InverterParams) *ws.WS {
	webSocket := ws.NewWS(inverterParams.host, inverterParams.port, inverterParams.path)
	if err := webSocket.Connect(); err != nil {
		log.Fatalln(err)
	}
	return webSocket
}

// flags defines, parses and validates command-line flags from os.Args[1:]
func flags() (*InverterParams, *mqtt.MqttParams, int) {

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

	inverterParams := &InverterParams{host: *host, port: *port, path: *path, data: *data}
	mqttParams := &mqtt.MqttParams{Server: *mqttServer, ClientId: *mqttClientId, Topic: *mqttTopic, User: *mqttUser, Password: *mqttPassword}

	// Validate flags
	validateInverterFlags(inverterParams)
	validateMqttFlags(mqttParams)

	return inverterParams, mqttParams, *sleepBetweenCalls
}

// validateInverterFlags validates all flags
func validateInverterFlags(inverterParams *InverterParams) {
	if inverterParams.host == "" {
		log.Fatalln("Required parameter 'host' not set!\n'sungrow-go -help' lists available parameters.")
	}

	inverterParams.types = strings.Split(inverterParams.data, ",")
	if len(inverterParams.types) < 1 {
		log.Fatalln("Required parameter 'data' not set or invalid value!\n'sungrow-go -help' lists available parameters and values.")
	}
	for _, t := range inverterParams.types {
		switch t {
		case "pv":
		case "battery":
		default:
			log.Fatalf("Invalid value \"%s\" for parameter 'data'!\n'sungrow-go -help' lists available parameters and values.\n", t)
		}
	}
}

func validateMqttFlags(params *mqtt.MqttParams) {
	if params.Server == "" {
		log.Fatalln("Missing parameter mqtt.server")
	}
}
