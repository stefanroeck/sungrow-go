// Package transport constructs inverter clients for the selected transport.
package transport

import (
	"fmt"
	"log"

	"github.com/sroeck/sungrow-go/inverter"
	"github.com/sroeck/sungrow-go/ws"
)

// Kind identifies an inverter communication transport.
type Kind string

const (
	WebSocket Kind = "websocket"
)

// Params groups settings for constructing a transport client.
type Params struct {
	Inverter  *inverter.Params
	WebSocket *ws.WebSocketParams
}

// New constructs and connects a client for the selected transport.
func New(kind Kind, params Params) (inverter.Client, error) {
	log.Printf("Constructing inverter client: %s", kind)
	switch kind {
	case WebSocket:
		if params.Inverter == nil {
			return nil, fmt.Errorf("inverter parameters are required for %q transport", kind)
		}
		if params.WebSocket == nil {
			return nil, fmt.Errorf("WebSocket parameters are required for %q transport", kind)
		}
		client, err := ws.Open(params.Inverter, params.WebSocket)
		if err != nil {
			return nil, err
		}
		log.Printf("Constructed inverter client: transport=%s, type=%T", kind, client)
		return client, nil
	default:
		return nil, fmt.Errorf("unsupported inverter transport %q", kind)
	}
}
