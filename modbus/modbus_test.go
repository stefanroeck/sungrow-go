package modbus

import (
	"testing"

	"github.com/aldas/go-modbus-client/packet"
)

func TestDecodeSG15RTRegisters(t *testing.T) {
	startAddress := firstRegisterReference - 1
	data := make([]byte, registerCount*2)
	put := func(reference, value uint16) {
		address := reference - 1
		offset := int(address-startAddress) * 2
		data[offset] = byte(value >> 8)
		data[offset+1] = byte(value)
	}
	put(5003, 20)     // 2.0 kWh today
	put(5004, 0xDCEB) // low word of 56555
	put(5005, 0x0000) // high word
	put(5006, 14122)  // running-time low word
	put(5007, 0)
	put(5008, 0xFFF1) // signed -15
	put(5017, 2053)   // DC power low word
	put(5018, 0)
	put(5031, 2051) // AC power low word
	put(5032, 0)
	put(5038, 0)
	put(5148, 5000)

	registers, err := packet.NewRegisters(data, startAddress)
	if err != nil {
		t.Fatal(err)
	}
	registers.WithByteOrder(packet.BigEndianLowWordFirst)
	values, err := decodeValues(registers)
	if err != nil {
		t.Fatal(err)
	}
	checks := map[string]any{
		"todayEnergy":      float64(2),
		"totalEnergy":      float64(56555),
		"totalRunningTime": float64(14122),
		"inverterTemp":     float64(-1.5),
		"sunPower":         float64(2.053),
		"activePower":      float64(2.051),
		"gridFrequency":    float64(50),
		"status":           "on grid operation",
	}
	for key, want := range checks {
		if got := values[key]; got != want {
			t.Errorf("%s = %#v; want %#v", key, got, want)
		}
	}
}
