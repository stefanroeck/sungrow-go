package main

import (
	"encoding/json"
	"log"
	"maps"
	"os"
	"time"

	"github.com/sroeck/sungrow-go/inverter"
	"github.com/sroeck/sungrow-go/mqtt"
	"github.com/sroeck/sungrow-go/params"
	"github.com/sroeck/sungrow-go/transport"
)

func main() {
	cfg := params.ParseFlags(os.Args[1:])
	log.Printf("Using inverter transport: %s", cfg.Transport)
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
		log.Printf("Failed to create inverter client; will retry next poll: %v", err)
		return
	}
	defer client.Close()

	receivedValues := make(map[inverter.DataKey]any)

	log.Println("Fetching pv data")
	pvValues, err := client.Pv()
	if err != nil {
		log.Printf("Failed to fetch inverter data; will retry next poll: %v", err)
		return
	}
	processPvResult(receivedValues, pvValues)

	if len(receivedValues) == 0 {
		log.Println("Skip sending MQTT data as no data have been returned from inverter")
		return
	}

	mqttClient.Send(receivedValues)
}

func processPvResult(targetMap map[inverter.DataKey]any, resultMap map[inverter.DataKey]any) {
	if len(resultMap) == 0 {
		log.Println("Warning: No data returned")
		return
	}

	//prettyPrintMap(resultMap)
	maps.Copy(targetMap, resultMap)
}

func prettyPrintMap(resultMap map[inverter.DataKey]any) {
	b, _ := json.MarshalIndent(resultMap, "", "  ")
	log.Print("Received the following values: ", string(b))
}
