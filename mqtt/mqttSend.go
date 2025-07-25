package mqtt

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttParams struct {
	Server   string
	ClientId string
	User     string
	Password string
	Topic    string
}

type MqttClient struct {
	client mqtt.Client
	topic  string
}

func NewMqttClient(params *MqttParams) *MqttClient {
	log.Println("Connecting to MQTT broker at", params.Server)

	opts := mqtt.NewClientOptions().AddBroker(params.Server).
		SetClientID(params.ClientId).
		SetUsername(params.User).
		SetPassword(params.Password).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second)

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Println("MQTT connection lost:", err)
	}

	opts.OnConnect = func(client mqtt.Client) {
		log.Println("MQTT connection (re-)established")
	}

	if strings.HasPrefix(params.Server, "ssl") {
		opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &MqttClient{
		client: client,
		topic:  params.Topic,
	}
}

func (m *MqttClient) Send(data map[string]any) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	token := m.client.Publish(m.topic, 0, false, jsonBytes)
	token.Wait()
	if token.Error() != nil {
		panic(token.Error())
	}

	log.Println("Successfully sent MQTT message")
}

func (m *MqttClient) Close() {
	log.Println("Closing MQTT broker.")

	m.client.Disconnect(250)
}
