// Package modbus implements the SG15RT Modbus TCP inverter transport.
package modbus

import (
	"context"
	"fmt"
	"log"
	"time"

	modbusclient "github.com/aldas/go-modbus-client"
	"github.com/aldas/go-modbus-client/packet"
	"github.com/sroeck/sungrow-go/inverter"
)

// Reference addresses in the README use mbpoll's one-based notation. The
// client library uses zero-based PDU addresses, so convert at request creation.
const firstRegisterReference uint16 = 5003
const (
	maxReadRegisterCount uint16 = 125
	registerCount        uint16 = 147 // through reference 5149, inclusive
)

type Client struct {
	client *modbusclient.Client
	params Params
}

// Open creates and connects an SG15RT Modbus TCP client.
func Open(params *Params) (inverter.Client, error) {
	if params == nil {
		return nil, fmt.Errorf("Modbus parameters are required")
	}
	if params.Host == "" {
		return nil, fmt.Errorf("Modbus host is required")
	}
	if params.Port < 1 || params.Port > 65535 {
		return nil, fmt.Errorf("Modbus port must be between 1 and 65535")
	}
	client := &Client{client: modbusclient.NewTCPClient(), params: *params}
	if err := client.Connect(); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

// Connect establishes the Modbus TCP connection.
func (c *Client) Connect() error {
	address := fmt.Sprintf("%s:%d", c.params.Host, c.params.Port)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.client.Connect(ctx, address); err != nil {
		return fmt.Errorf("connect to Modbus inverter %s: %w", address, err)
	}
	return nil
}

// Close closes the Modbus TCP connection.
func (c *Client) Close() {
	if c.client != nil {
		if err := c.client.Close(); err != nil {
			log.Printf("close Modbus connection: %v", err)
		}
	}
}

// Pv reads the SG15RT read-only FC04 register range and maps values to the
// application keys used by the WebSocket transport.
func (c *Client) Pv() (map[inverter.DataKey]any, error) {
	startAddress := firstRegisterReference - 1
	data := make([]byte, int(registerCount)*2)
	for offset := uint16(0); offset < registerCount; {
		count := min(maxReadRegisterCount, registerCount-offset)
		chunkStart := startAddress + offset
		req, err := packet.NewReadInputRegistersRequestTCP(c.params.UnitID, chunkStart, count)
		if err != nil {
			return nil, fmt.Errorf("build SG15RT register request at address %d: %w", chunkStart, err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := c.client.Do(ctx, req)
		cancel()
		if err != nil {
			return nil, fmt.Errorf("read SG15RT input registers at address %d: %w", chunkStart, err)
		}
		readResp, ok := resp.(*packet.ReadInputRegistersResponseTCP)
		if !ok {
			return nil, fmt.Errorf("unexpected Modbus response type %T", resp)
		}
		chunkRegisters, err := readResp.AsRegisters(chunkStart)
		if err != nil {
			return nil, fmt.Errorf("decode SG15RT input registers at address %d: %w", chunkStart, err)
		}
		copy(data[int(offset)*2:], chunkRegisters.Data())
		offset += count
	}
	registers, err := packet.NewRegisters(data, startAddress)
	if err != nil {
		return nil, fmt.Errorf("assemble SG15RT input registers: %w", err)
	}

	// Register pairs use low word first per the SG15RT README mapping.
	registers.WithByteOrder(packet.BigEndianLowWordFirst)
	return decodeValues(registers)
}

func decodeValues(registers *packet.Registers) (map[inverter.DataKey]any, error) {
	u16 := func(reference uint16) (uint16, error) { return registers.Uint16(reference - 1) }
	u32 := func(reference uint16) (uint32, error) { return registers.Uint32(reference - 1) }
	values := make(map[inverter.DataKey]any, 16)
	put16 := func(key inverter.DataKey, reference uint16, scale float64) error {
		raw, err := u16(reference)
		if err != nil {
			return fmt.Errorf("decode %s at reference %d: %w", key, reference, err)
		}
		values[key] = float64(raw) / scale
		return nil
	}
	put32 := func(key inverter.DataKey, reference uint16, scale float64) error {
		raw, err := u32(reference)
		if err != nil {
			return fmt.Errorf("decode %s at reference %d: %w", key, reference, err)
		}
		values[key] = float64(raw) / scale
		return nil
	}
	for _, read := range []func() error{
		func() error { return put16(inverter.TodayEnergy, 5003, 10) },
		func() error { return put32(inverter.TotalEnergy, 5004, 1) },
		func() error { return put32(inverter.TotalRunningTime, 5006, 1) },
		func() error {
			raw, err := u16(5008)
			if err != nil {
				return fmt.Errorf("decode %s at reference 5008: %w", inverter.InverterTemp, err)
			}
			values[inverter.InverterTemp] = float64(int16(raw)) / 10
			return nil
		},
		func() error { return put32(inverter.SunPower, 5017, 1000) },
		func() error { return put16(inverter.VoltagePhaseA, 5019, 10) },
		func() error { return put16(inverter.VoltagePhaseB, 5020, 10) },
		func() error { return put16(inverter.VoltagePhaseC, 5021, 10) },
		func() error { return put16(inverter.CurrentPhaseA, 5022, 10) },
		func() error { return put16(inverter.CurrentPhaseB, 5023, 10) },
		func() error { return put16(inverter.CurrentPhaseC, 5024, 10) },
		func() error { return put32(inverter.ActivePower, 5031, 1000) },
		func() error { return put16(inverter.GridFrequency, 5148, 100) },
		func() error {
			raw, err := u16(5038)
			if err != nil {
				return fmt.Errorf("decode %s at reference 5038: %w", inverter.Status, err)
			}
			values[inverter.Status] = statusString(raw)
			return nil
		},
		func() error { return put16(inverter.ArrayInsulationResistance, 5071, 1) },
		func() error { return put16(inverter.BusVoltage, 5147, 10) },
	} {
		if err := read(); err != nil {
			return nil, err
		}
	}
	return values, nil
}

func statusString(code uint16) string {
	if code == 0 {
		return "on grid operation"
	}
	return fmt.Sprintf("unknown (%d)", code)
}
