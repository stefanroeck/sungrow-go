package main

import (
	"encoding/json"
	"log"
	"maps"
	"os"
	"time"

	"github.com/sroeck/sungrow-go/mqtt"
	"github.com/sroeck/sungrow-go/params"
	"github.com/sroeck/sungrow-go/transport"
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
		fetchDataFromInverterAndSendToMqtt(cfg, mqttClient)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fetchDataFromInverterAndSendToMqtt(cfg, mqttClient)
			}
		}
	}()

	// block forever
	select {}
}

func fetchDataFromInverterAndSendToMqtt(cfg *params.Config, mqttClient *mqtt.MqttClient) {
	client, err := transport.New(transport.Kind(cfg.Transport), transport.Params{
		Inverter:  cfg.InverterParams,
		WebSocket: cfg.WebSocketParams,
		Modbus:    cfg.ModbusParams,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	var receivedValues map[string]any = make(map[string]any)

	log.Println("Fetching pv data")
	pvValues, err := client.Pv()
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
