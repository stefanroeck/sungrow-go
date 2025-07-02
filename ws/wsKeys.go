package ws

var (
	pvKeys = Keys{
		"I18N_COMMON_TOTAL_DCPOWER":                     {Name: "sunPower", KeyType: KeyTypes.Number},
		"I18N_COMMON_PV_DAYILY_ENERGY_GENERATION":       {Name: "sunPowerDay", KeyType: KeyTypes.Number},
		"I18N_COMMON_PV_TOTAL_ENERGY_GENERATION":        {Name: "sunPowerTotal", KeyType: KeyTypes.Number},
		"I18N_COMMON_FEED_NETWORK_TOTAL_ACTIVE_POWER":   {Name: "netFeedIn", KeyType: KeyTypes.Number},
		"I18N_COMMON_DAILY_FEED_NETWORK_VOLUME":         {Name: "netFeedInDay", KeyType: KeyTypes.Number},
		"I18N_COMMON_TOTAL_FEED_NETWORK_VOLUME":         {Name: "netFeedInTotal", KeyType: KeyTypes.Number},
		"I18N_CONFIG_KEY_4060":                          {Name: "netPower", KeyType: KeyTypes.Number},
		"I18N_COMMON_ENERGY_GET_FROM_GRID_DAILY":        {Name: "netPowerDay", KeyType: KeyTypes.Number},
		"I18N_COMMON_TOTAL_ELECTRIC_GRID_GET_POWER":     {Name: "netPowerTotal", KeyType: KeyTypes.Number},
		"I18N_COMMON_LOAD_TOTAL_ACTIVE_POWER":           {Name: "consumption", KeyType: KeyTypes.Number},
		"I18N_CONFIG_KEY_1001188":                       {Name: "consumptionRate", KeyType: KeyTypes.Number},
		"I18N_COMMON_AIR_TEM_INSIDE_MACHINE":            {Name: "inverterTemp", KeyType: KeyTypes.Number},
		"I18N_COMMON_TOTAL_GRID_RUNNING_TIME":           {Name: "totalRunningTime", KeyType: KeyTypes.Number},
		"I18N_COMMON_DAILY_POWER_YIELD":                 {Name: "todayEnergy", KeyType: KeyTypes.Number},
		"I18N_COMMON_TOTAL_YIELD":                       {Name: "totalEnergy", KeyType: KeyTypes.Number},
		"I18N_COMMON_TOTAL_ACTIVE_POWER":                {Name: "activePower", KeyType: KeyTypes.Number},
		"I18N_COMMON_BUS_VOLTAGE":                       {Name: "busVoltage", KeyType: KeyTypes.Number},
		"I18N_COMMON_RUNNING_STATE":                     {Name: "status", KeyType: KeyTypes.String},
		"I18N_COMMONUA":                                 {Name: "voltagePhaseA", KeyType: KeyTypes.Number},
		"I18N_COMMON_UB":                                {Name: "voltagePhaseB", KeyType: KeyTypes.Number},
		"I18N_COMMON_UC":                                {Name: "voltagePhaseC", KeyType: KeyTypes.Number},
		"I18N_COMMON_FRAGMENT_RUN_TYPE1":                {Name: "currentPhaseA", KeyType: KeyTypes.Number},
		"I18N_COMMON_IB":                                {Name: "currentPhaseB", KeyType: KeyTypes.Number},
		"I18N_COMMON_IC":                                {Name: "currentPhaseC", KeyType: KeyTypes.Number},
		"I18N_COMMON_GRID_FREQUENCY":                    {Name: "gridFrequency", KeyType: KeyTypes.Number},
		"I18N_COMMON_SQUARE_ARRAY_INSULATION_IMPEDANCE": {Name: "arrayInsulationResistance", KeyType: KeyTypes.Number},
		"I18N_COMMON_BATTERY_SOC":                       {Name: "batteryLevel", KeyType: KeyTypes.Number},
		"I18N_CONFIG_KEY_3907":                          {Name: "batteryCharge", KeyType: KeyTypes.Number},
		"I18N_CONFIG_KEY_3921":                          {Name: "batteryDischarge", KeyType: KeyTypes.Number},
		"I18N_COMMON_BATTARY_HEALTH":                    {Name: "batteryHealth", KeyType: KeyTypes.Number},
		"I18N_COMMON_BATTERY_TEMPERATURE":               {Name: "batteryTemp", KeyType: KeyTypes.Number},
	}

	valueMapping = map[string]string{
		"I18N_COMMON_STANDBY":    "standby",
		"I18N_COMMON_STATUS_RUN": "running",
	}
)
