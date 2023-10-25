package mqtt

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttParams struct {
	Server string
	ClientId string
	User   string
	Password    string
	Topic  string
}

func Send(params MqttParams, data map[string]float64) {
	broker := mqtt.NewClientOptions().AddBroker(params.Server)
	broker.SetClientID(params.ClientId)
	broker.SetUsername(params.User).SetPassword(params.Password)

	client := mqtt.NewClient(broker)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connected to", params.Server)

	json, err := json.Marshal(data)
	if (err != nil) {
		panic(err)
	}
	token := client.Publish(params.Topic, 0, false, string(json))
	token.Wait()
	fmt.Println("Published to topic", params.Topic)
	
	client.Disconnect(250)
	fmt.Println("Disconnected from", params.Server)
}