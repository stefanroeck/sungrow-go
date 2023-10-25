package ws

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type WS struct {
	ip    net.IP
	port  int
	path  string
	conn  *websocket.Conn
	token string
	uid   string
}

// NewWS returns a new WS instance.
func NewWS(ip net.IP, port int, path string) *WS {
	ws := &WS{
		ip:   ip,
		port: port,
		path: path,
	}

	return ws
}

// Connect connects to the inverter using the WebSocket protocol.
func (ws *WS) Connect() (err error) {
	// Connect to WebSocket
	origin := fmt.Sprintf("http://%s", ws.ip.String())
	url := fmt.Sprintf("ws://%s:%d%s", ws.ip.String(), ws.port, ws.path)
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
func (ws *WS) Pv(keyList Keys, separator string) (err error, res map[string]float64) {
	return ws.fetch("real", keyList, separator)
}

// Battery fetches battery data from the inverter.
func (ws *WS) Battery(keyList Keys, separator string) (err error, res map[string]float64) {
	return ws.fetch("real_battery", keyList, separator)
}

// fetch fetches data from the inverter.
func (ws *WS) fetch(service string, keyList Keys, separator string) (err error, res map[string]float64) {
	req := RequestReal{"de_de", ws.token, ws.uid, service, time.Now().UnixMilli()}
	if err := websocket.JSON.Send(ws.conn, &req); err != nil {
		return err, nil
	}
	resp := ResponseReal{}
	if err := websocket.JSON.Receive(ws.conn, &resp); err != nil {
		return err, nil
	}

	// Output values
	result := make(map[string]float64)
	for _, row := range resp.ResultData.List {
		if _, exists := keyList[row.DataName]; exists {
			val, _ := strconv.ParseFloat(row.DataValue, 64)
			result[keyList[row.DataName]] = val
		}
	}

	return nil, result
}
