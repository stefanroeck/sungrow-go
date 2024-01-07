package mqtt

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttParams struct {
	Server   string
	ClientId string
	User     string
	Password string
	Topic    string
}

func Send(params *MqttParams, data map[string]any) {
	log.Println("Sending mqtt message to", params.Server, "( clientId:", params.ClientId, ", topic:", params.Topic, ", user:", params.User, ")")
	broker := mqtt.NewClientOptions().AddBroker(params.Server)
	broker.SetClientID(params.ClientId)
	broker.SetUsername(params.User)
	broker.SetPassword(params.Password)

	if strings.HasPrefix(params.Server, "ssl") {
		broker.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	client := mqtt.NewClient(broker)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	token := client.Publish(params.Topic, 0, false, string(json))
	token.Wait()
	if token.Error() != nil {
		panic(token.Error())
	}

	client.Disconnect(250)
	log.Println("Successfully sent mqtt message")
}
