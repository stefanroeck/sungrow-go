package ws

import "github.com/sroeck/sungrow-go/inverter"

var (
	pvKeys = Keys{
		"I18N_COMMON_TOTAL_DCPOWER":                     {Name: inverter.SunPower, KeyType: KeyTypes.Number},
		"I18N_COMMON_PV_DAYILY_ENERGY_GENERATION":       {Name: inverter.SunPowerDay, KeyType: KeyTypes.Number},
		"I18N_COMMON_PV_TOTAL_ENERGY_GENERATION":        {Name: inverter.SunPowerTotal, KeyType: KeyTypes.Number},
		"I18N_COMMON_FEED_NETWORK_TOTAL_ACTIVE_POWER":   {Name: inverter.NetFeedIn, KeyType: KeyTypes.Number},
		"I18N_COMMON_DAILY_FEED_NETWORK_VOLUME":         {Name: inverter.NetFeedInDay, KeyType: KeyTypes.Number},
		"I18N_COMMON_TOTAL_FEED_NETWORK_VOLUME":         {Name: inverter.NetFeedInTotal, KeyType: KeyTypes.Number},
		"I18N_CONFIG_KEY_4060":                          {Name: inverter.NetPower, KeyType: KeyTypes.Number},
		"I18N_COMMON_ENERGY_GET_FROM_GRID_DAILY":        {Name: inverter.NetPowerDay, KeyType: KeyTypes.Number},
		"I18N_COMMON_TOTAL_ELECTRIC_GRID_GET_POWER":     {Name: inverter.NetPowerTotal, KeyType: KeyTypes.Number},
		"I18N_COMMON_LOAD_TOTAL_ACTIVE_POWER":           {Name: inverter.Consumption, KeyType: KeyTypes.Number},
		"I18N_CONFIG_KEY_1001188":                       {Name: inverter.ConsumptionRate, KeyType: KeyTypes.Number},
		"I18N_COMMON_AIR_TEM_INSIDE_MACHINE":            {Name: inverter.InverterTemp, KeyType: KeyTypes.Number},
		"I18N_COMMON_TOTAL_GRID_RUNNING_TIME":           {Name: inverter.TotalRunningTime, KeyType: KeyTypes.Number},
		"I18N_COMMON_DAILY_POWER_YIELD":                 {Name: inverter.TodayEnergy, KeyType: KeyTypes.Number},
		"I18N_COMMON_TOTAL_YIELD":                       {Name: inverter.TotalEnergy, KeyType: KeyTypes.Number},
		"I18N_COMMON_TOTAL_ACTIVE_POWER":                {Name: inverter.ActivePower, KeyType: KeyTypes.Number},
		"I18N_COMMON_BUS_VOLTAGE":                       {Name: inverter.BusVoltage, KeyType: KeyTypes.Number},
		"I18N_COMMON_RUNNING_STATE":                     {Name: inverter.Status, KeyType: KeyTypes.String},
		"I18N_COMMONUA":                                 {Name: inverter.VoltagePhaseA, KeyType: KeyTypes.Number},
		"I18N_COMMON_UB":                                {Name: inverter.VoltagePhaseB, KeyType: KeyTypes.Number},
		"I18N_COMMON_UC":                                {Name: inverter.VoltagePhaseC, KeyType: KeyTypes.Number},
		"I18N_COMMON_FRAGMENT_RUN_TYPE1":                {Name: inverter.CurrentPhaseA, KeyType: KeyTypes.Number},
		"I18N_COMMON_IB":                                {Name: inverter.CurrentPhaseB, KeyType: KeyTypes.Number},
		"I18N_COMMON_IC":                                {Name: inverter.CurrentPhaseC, KeyType: KeyTypes.Number},
		"I18N_COMMON_GRID_FREQUENCY":                    {Name: inverter.GridFrequency, KeyType: KeyTypes.Number},
		"I18N_COMMON_SQUARE_ARRAY_INSULATION_IMPEDANCE": {Name: inverter.ArrayInsulationResistance, KeyType: KeyTypes.Number},
		"I18N_COMMON_BATTERY_SOC":                       {Name: inverter.BatteryLevel, KeyType: KeyTypes.Number},
		"I18N_CONFIG_KEY_3907":                          {Name: inverter.BatteryCharge, KeyType: KeyTypes.Number},
		"I18N_CONFIG_KEY_3921":                          {Name: inverter.BatteryDischarge, KeyType: KeyTypes.Number},
		"I18N_COMMON_BATTARY_HEALTH":                    {Name: inverter.BatteryHealth, KeyType: KeyTypes.Number},
		"I18N_COMMON_BATTERY_TEMPERATURE":               {Name: inverter.BatteryTemp, KeyType: KeyTypes.Number},
	}

	valueMapping = map[string]string{
		"I18N_COMMON_STANDBY":           "standby",
		"I18N_COMMON_STATUS_RUN":        "running",
		"I18N_COMMON_ON_GRID_OPERATION": "on grid operation",
	}
)
