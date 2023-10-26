package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/sroeck/sungrow-go/mqtt"
	"github.com/sroeck/sungrow-go/ws"
)

type inverter struct {
	ws        *ws.WS
	ipS       string
	ip        net.IP
	port      int
	path      string
	data      string
	separator string
	types     []string
}

func main() {

	// Flags
	inv, mqttParams := flags()
	// Connect to inverter
	inv.ws = ws.NewWS(inv.ip, inv.port, inv.path)
	if err := inv.ws.Connect(); err != nil {
		log.Fatalln(err)
	}
	defer inv.ws.Close()

	// Output timestamp row
	fmt.Printf("%s%s%s%s%s\n", "time", inv.separator, time.Now().Format(time.RFC3339), inv.separator, "RFC3339")

	// Fetch values from inverter
	for _, t := range inv.types {
		switch t {
		case "pv":
			fetchAndProcessValues(inv, mqttParams)
			break
		case "battery":
			_, _ = inv.ws.Battery(batteryKeys, inv.separator)
			break
		}
	}
}

func fetchAndProcessValues(inv inverter, mqttParams mqtt.MqttParams) {
	_, values := inv.ws.Pv(pvKeys, inv.separator)
	for k, v := range values {
		fmt.Println(k, "=", v)
	}
	mqtt.Send(mqttParams, values)
}

// flags defines, parses and validates command-line flags from os.Args[1:]
func flags() (inverter, mqtt.MqttParams) {

	ipS := flag.String("ip", "", "IP address of the Sungrow inverter")
	port := flag.Int("port", 8082, "WebSocket port of the Sungrow inverter")
	path := flag.String("path", "/ws/home/overview", "Server path from where data is requested")
	data := flag.String("data", "pv,battery", "Select the data to be requested comma separated.\nPossible values are \"pv\" and \"battery\"")
	separator := flag.String("separator", ",", "Output data separator")
	mqttServer := flag.String("mqtt.server", "", "mqtt server incl. protocol, e.g. mqtt://localhost:1883. For TLS use ssl scheme, e.g. ssl://localhost:8883")
	mqttUser := flag.String("mqtt.user", "", "mqtt user name")
	mqttPassword := flag.String("mqtt.password", "", "mqtt password")
	mqttClientId := flag.String("mqtt.clientId", "", "mqtt clientId that is used for publishing")
	mqttTopic := flag.String("mqtt.topic", "topic", "mqtt topic to which the data are published")
	flag.Parse()

	inv := &inverter{ipS: *ipS, port: *port, path: *path, data: *data, separator: *separator}
	
	mqttParams := &mqtt.MqttParams{Server: *mqttServer, ClientId: *mqttClientId, Topic: *mqttTopic, User: *mqttUser, Password: *mqttPassword}

	// Validate flags
	flagsValidate(inv)
	flagsValidateMqtt(*mqttParams)

	return *inv, *mqttParams
}

// flagsValidate validates all flags
func flagsValidate(inv *inverter) {
	if inv.ip = net.ParseIP(inv.ipS); inv.ip == nil {
		log.Fatalln("Required parameter 'ip' not set or invalid ip address!\n'sungrow-go -help' lists available parameters.")
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

