package ws

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type WS struct {
	host  string
	port  int
	path  string
	conn  *websocket.Conn
	token string
	uid   string
}

// NewWS returns a new WS instance.
func NewWS(host string, port int, path string) *WS {
	ws := &WS{
		host: host,
		port: port,
		path: path,
	}

	return ws
}

// Connect connects to the inverter using the WebSocket protocol.
func (ws *WS) Connect() (err error) {
	// Connect to WebSocket
	origin := fmt.Sprintf("http://%s", ws.host)
	url := fmt.Sprintf("ws://%s:%d%s", ws.host, ws.port, ws.path)
	fmt.Printf("Connecting to %s\n", url)

	ws.conn, err = websocket.Dial(url, "", origin)
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

// Battery fetches battery data from the inverter.
func (ws *WS) Battery() (res map[string]any, err error) {
	return ws.fetch("real_battery", batteryKeys)
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

	// Output values
	result := make(map[string]any)
	for _, row := range resp.ResultData.List {
		//fmt.Printf("%s\n", row)
		if param, exists := keyList[row.DataName]; exists {
			var val any
			if param.KeyType == KeyTypes.Number {
				if row.DataValue == "--" { // API return "--" for empty values
					val = 0
				} else {
					val, err = strconv.ParseFloat(row.DataValue, 64)
					if err != nil {
						fmt.Printf("Cannot convert %s to number. Using as string.\n", row.DataName)
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
