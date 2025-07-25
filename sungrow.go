package main

import (
	"encoding/json"
	"flag"
	"log"
	"maps"
	"time"

	"github.com/sroeck/sungrow-go/mqtt"
	"github.com/sroeck/sungrow-go/ws"
)

func main() {
	inverterParams, mqttParams, sleepTimeInSeconds := flags()
	log.Printf("Polling the inverter every %d seconds\n", sleepTimeInSeconds)

	ticker := time.NewTicker(time.Duration(sleepTimeInSeconds) * time.Second)
	done := make(chan bool)
	defer ticker.Stop()

	mqttClient := mqtt.NewMqttClient(mqttParams)
	defer mqttClient.Close()

	go func() {
		// run immediately and then at the configured interval
		fetchDataFromInverterAndSendToMqtt(inverterParams, mqttClient)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fetchDataFromInverterAndSendToMqtt(inverterParams, mqttClient)
			}
		}
	}()

	// block forever
	select {}
}

func fetchDataFromInverterAndSendToMqtt(inverterParams *ws.InverterParams, mqttClient *mqtt.MqttClient) {
	webSocket := openWebSocket(inverterParams)
	defer webSocket.Close()

	var receivedValues map[string]any = make(map[string]any)

	log.Println("Fetching pv data")
	pvValues, err := webSocket.Pv()
	processWsResult(receivedValues, pvValues, err)

	if len(receivedValues) == 0 {
		log.Println("Skip sending MQTT data as no data have been returned from inverter")
		return
	}

	mqttClient.Send(receivedValues)
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

	//prettyPrintMap(resultMap)
	maps.Copy(targetMap, resultMap)
}

func prettyPrintMap(resultMap map[string]any) {
	b, _ := json.MarshalIndent(resultMap, "", "  ")
	log.Print("Received the following values: ", string(b))
}

func openWebSocket(inverterParams *ws.InverterParams) *ws.WS {
	webSocket := ws.NewWS(inverterParams)
	if err := webSocket.Connect(); err != nil {
		log.Fatalln(err)
	}
	return webSocket
}

// flags defines, parses and validates command-line flags from os.Args[1:]
func flags() (*ws.InverterParams, *mqtt.MqttParams, int) {
	protocol := flag.String("protocol", "ws", "WebSocket protocol to be used, either \"ws\" or \"wss\"")
	host := flag.String("host", "", "Hostname/IP address of the Sungrow inverter")
	port := flag.Int("port", 8082, "WebSocket port of the Sungrow inverter")
	user := flag.String("user", "", "Username for the Sungrow inverter web ui login, e.g. admin")
	password := flag.String("password", "", "Password the Sungrow inverter web ui login")
	path := flag.String("path", "/ws/home/overview", "Server path from where data is requested")
	mqttServer := flag.String("mqtt.server", "", "mqtt server incl. protocol, e.g. mqtt://localhost:1883. For TLS use ssl scheme, e.g. ssl://localhost:8883")
	mqttUser := flag.String("mqtt.user", "", "mqtt user name")
	mqttPassword := flag.String("mqtt.password", "", "mqtt password")
	mqttClientId := flag.String("mqtt.clientId", "", "mqtt clientId that is used for publishing")
	mqttTopic := flag.String("mqtt.topic", "topic", "mqtt topic to which the data are published")
	sleepBetweenCalls := flag.Int("sleep", 10, "sleep time in seconds between inverter calls.")
	flag.Parse()

	inverterParams := &ws.InverterParams{Protocol: *protocol, Host: *host, Port: *port, User: *user, Password: *password, Path: *path}
	mqttParams := &mqtt.MqttParams{Server: *mqttServer, ClientId: *mqttClientId, Topic: *mqttTopic, User: *mqttUser, Password: *mqttPassword}

	// Validate flags
	validateInverterFlags(inverterParams)
	validateMqttFlags(mqttParams)

	return inverterParams, mqttParams, *sleepBetweenCalls
}

// validateInverterFlags validates all flags
func validateInverterFlags(inverterParams *ws.InverterParams) {
	if inverterParams.Host == "" {
		log.Fatalln("Required parameter 'host' not set!\n'sungrow-go -help' lists available parameters.")
	}
}

func validateMqttFlags(params *mqtt.MqttParams) {
	if params.Server == "" {
		log.Fatalln("Missing parameter mqtt.server")
	}
}
