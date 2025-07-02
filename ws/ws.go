package ws

import (
	"crypto/tls"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type WS struct {
	protocol string // "ws" or "wss"
	host     string
	port     int
	user     string
	password string
	path     string
	conn     *websocket.Conn
	token    string
	uid      string
}

// NewWS returns a new WS instance.
func NewWS(inverterParams *InverterParams) *WS {
	ws := &WS{
		protocol: inverterParams.Protocol,
		host:     inverterParams.Host,
		port:     inverterParams.Port,
		user:     inverterParams.User,
		password: inverterParams.Password,
		path:     inverterParams.Path,
	}

	return ws
}

// Connect connects to the inverter using the WebSocket protocol.
func (ws *WS) Connect() (err error) {
	// Connect to WebSocket
	origin := fmt.Sprintf("https://%s", ws.host)
	url := fmt.Sprintf("wss://%s:%d%s", ws.host, ws.port, ws.path)
	log.Println("Connecting to", url, "with origin", origin)

	config, err := websocket.NewConfig(url, origin)
	if err != nil {
		return err
	}

	config.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	ws.conn, err = websocket.DialConfig(config)
	if err != nil {
		return err
	}

	// Connect to service
	ws.token = uuid.New().String()
	req := RequestConnect{"de_de", ws.token, "connect"}
	if err := websocket.JSON.Send(ws.conn, &req); err != nil {
		return err
	}
	res := ResponseConnect{}
	if err := websocket.JSON.Receive(ws.conn, &res); err != nil {
		return err
	}
	ws.token = res.ResultData.Token
	ws.uid = strconv.Itoa(res.ResultData.Uid)

	if res.ResultMsg != "success" {
		ws.Close()
		return fmt.Errorf("connected but connection request failed")
	}

	if ws.user != "" && ws.password != "" {
		// Login
		reqLogin := RequestLogin{"de_de", ws.token, "login", ws.user, ws.password}
		if err := websocket.JSON.Send(ws.conn, &reqLogin); err != nil {
			return err
		}
		resLogin := ResponseLogin{}
		if err := websocket.JSON.Receive(ws.conn, &resLogin); err != nil {
			return err
		}

		if resLogin.ResultMsg != "success" {
			ws.Close()
			return fmt.Errorf("login request failed. Please check user/password params: %s", resLogin.ResultMsg)
		}
		log.Printf("User %s successfully logged in", ws.user)

		ws.token = resLogin.ResultData.Token
	}

	return err
}

// Close closes the connection.
func (ws *WS) Close() {
	if ws.conn != nil {
		_ = ws.conn.Close()
	}
}

// Pv fetches pv data from the inverter.
func (ws *WS) Pv() (res map[string]any, err error) {
	return ws.fetch("real", pvKeys)
}

// fetch fetches data from the inverter.
func (ws *WS) fetch(service string, keyList Keys) (res map[string]any, err error) {
	req := RequestReal{"de_de", ws.token, ws.uid, service, time.Now().UnixMilli()}
	if err := websocket.JSON.Send(ws.conn, &req); err != nil {
		return nil, err
	}
	resp := ResponseReal{}
	if err := websocket.JSON.Receive(ws.conn, &resp); err != nil {
		return nil, err
	}
	if resp.ResultMsg != "success" {
		ws.Close()
		return nil, fmt.Errorf("request for service '%s' failed: %s", service, resp.ResultMsg)
	}

	// Output values
	result := make(map[string]any)
	for _, row := range resp.ResultData.List {
		//log.Println(row)
		if param, exists := keyList[row.DataName]; exists {
			var val any
			if param.KeyType == KeyTypes.Number {
				if row.DataValue == "--" { // API return "--" for empty values
					val = 0
				} else {
					val, err = strconv.ParseFloat(row.DataValue, 64)
					if err != nil {
						log.Printf("WARN: Cannot convert %s to number. Using as string.\n", row.DataName)
						val = row.DataValue
					}
				}
			} else if mappedValue, exists := valueMapping[row.DataValue]; exists {
				val = mappedValue
			} else {
				val = row.DataValue
			}
			result[param.Name] = val
		}
	}

	return result, nil
}
