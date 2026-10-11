// Package inverter defines the transport-neutral inverter API and connection
// parameters shared by the supported communication implementations.
package inverter

// DataKey identifies a field in the map returned by Client.Pv.
type DataKey string

// PvDataKeys are the transport-neutral keys returned by Client.Pv.
const (
	SunPower                  DataKey = "sunPower"
	SunPowerDay               DataKey = "sunPowerDay"
	SunPowerTotal             DataKey = "sunPowerTotal"
	NetFeedIn                 DataKey = "netFeedIn"
	NetFeedInDay              DataKey = "netFeedInDay"
	NetFeedInTotal            DataKey = "netFeedInTotal"
	NetPower                  DataKey = "netPower"
	NetPowerDay               DataKey = "netPowerDay"
	NetPowerTotal             DataKey = "netPowerTotal"
	Consumption               DataKey = "consumption"
	ConsumptionRate           DataKey = "consumptionRate"
	InverterTemp              DataKey = "inverterTemp"
	TotalRunningTime          DataKey = "totalRunningTime"
	TodayEnergy               DataKey = "todayEnergy"
	TotalEnergy               DataKey = "totalEnergy"
	ActivePower               DataKey = "activePower"
	BusVoltage                DataKey = "busVoltage"
	Status                    DataKey = "status"
	VoltagePhaseA             DataKey = "voltagePhaseA"
	VoltagePhaseB             DataKey = "voltagePhaseB"
	VoltagePhaseC             DataKey = "voltagePhaseC"
	CurrentPhaseA             DataKey = "currentPhaseA"
	CurrentPhaseB             DataKey = "currentPhaseB"
	CurrentPhaseC             DataKey = "currentPhaseC"
	GridFrequency             DataKey = "gridFrequency"
	ArrayInsulationResistance DataKey = "arrayInsulationResistance"
	BatteryLevel              DataKey = "batteryLevel"
	BatteryCharge             DataKey = "batteryCharge"
	BatteryDischarge          DataKey = "batteryDischarge"
	BatteryHealth             DataKey = "batteryHealth"
	BatteryTemp               DataKey = "batteryTemp"
)

// Client is the common interface implemented by inverter communication
// transports.
type Client interface {
	Connect() error
	Pv() (map[DataKey]any, error)
	Close()
}

// Params contains connection settings shared by inverter transports. Transport
// specific settings belong in the corresponding implementation's own config
// type.
type Params struct {
	Host string
}
