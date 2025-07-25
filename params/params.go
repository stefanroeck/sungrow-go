package params

import (
	"flag"
	"log"

	"github.com/sroeck/sungrow-go/mqtt"
	"github.com/sroeck/sungrow-go/ws"
)

// Config holds all parsed parameters
type Config struct {
	InverterParams        *ws.InverterParams
	MqttParams            *mqtt.MqttParams
	SleepBetweenCallsSecs int
}

func ParseFlags(args []string) *Config {
	fs := flag.NewFlagSet("sungrow-go", flag.ContinueOnError)

	protocol := fs.String("protocol", "ws", "WebSocket protocol to be used, either \"ws\" or \"wss\"")
	host := fs.String("host", "", "Hostname/IP address of the Sungrow inverter")
	port := fs.Int("port", 8082, "WebSocket port of the Sungrow inverter")
	user := fs.String("user", "", "Username for the Sungrow inverter web ui login, e.g. admin")
	password := fs.String("password", "", "Password the Sungrow inverter web ui login")
	path := fs.String("path", "/ws/home/overview", "Server path from where data is requested")
	mqttServer := fs.String("mqtt.server", "", "mqtt server incl. protocol, e.g. mqtt://localhost:1883. For TLS use ssl scheme, e.g. ssl://localhost:8883")
	mqttUser := fs.String("mqtt.user", "", "mqtt user name")
	mqttPassword := fs.String("mqtt.password", "", "mqtt password")
	mqttClientId := fs.String("mqtt.clientId", "", "mqtt clientId that is used for publishing")
	mqttTopic := fs.String("mqtt.topic", "topic", "mqtt topic to which the data are published")
	mqttSkipSSLVerify := fs.Bool("mqtt.skipSSLVerify", false, "Skip SSL verification for MQTT connection")
	sleepBetweenCalls := fs.Int("sleep", 10, "sleep time in seconds between inverter calls.")

	err := fs.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	inverterParams := &ws.InverterParams{Protocol: *protocol, Host: *host, Port: *port, User: *user, Password: *password, Path: *path}
	mqttParams := &mqtt.MqttParams{Server: *mqttServer, ClientId: *mqttClientId, Topic: *mqttTopic, User: *mqttUser, Password: *mqttPassword, SkipSSLVerify: *mqttSkipSSLVerify}

	// You can validate here or return errors if you want, skipping for brevity
	validateInverterFlags(inverterParams)
	validateMqttFlags(mqttParams)

	return &Config{
		InverterParams:        inverterParams,
		MqttParams:            mqttParams,
		SleepBetweenCallsSecs: *sleepBetweenCalls,
	}
}

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
