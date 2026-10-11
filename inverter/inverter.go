// Package inverter defines the transport-neutral inverter API and connection
// parameters shared by the supported communication implementations.
package inverter

// Client is the common interface implemented by inverter communication
// transports.
type Client interface {
	Connect() error
	Pv() (map[string]any, error)
	Close()
}

// Params contains connection settings shared by inverter transports. Transport
// specific settings belong in the corresponding implementation's own config
// type.
type Params struct {
	Host string
}
