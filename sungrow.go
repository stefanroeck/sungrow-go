package main

import (
	"encoding/json"
	"log"
	"maps"
	"os"
	"time"

	"github.com/sroeck/sungrow-go/mqtt"
	"github.com/sroeck/sungrow-go/params"
	"github.com/sroeck/sungrow-go/ws"
)

func main() {
	cfg := params.ParseFlags(os.Args[1:])
	log.Printf("Polling the inverter every %d seconds\n", cfg.SleepBetweenCallsSecs)

	ticker := time.NewTicker(time.Duration(cfg.SleepBetweenCallsSecs) * time.Second)
	done := make(chan bool)
	defer ticker.Stop()

	mqttClient := mqtt.NewMqttClient(cfg.MqttParams)
	defer mqttClient.Close()

	go func() {
		// run immediately and then at the configured interval
		fetchDataFromInverterAndSendToMqtt(cfg.InverterParams, mqttClient)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fetchDataFromInverterAndSendToMqtt(cfg.InverterParams, mqttClient)
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
