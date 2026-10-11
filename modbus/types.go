package modbus

// Params contains connection settings for the Modbus TCP transport.
type Params struct {
	Host   string
	Port   int
	UnitID uint8
}
