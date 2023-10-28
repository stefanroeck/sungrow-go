package main

import "github.com/sroeck/sungrow-go/ws"

var (
	pvKeys = ws.Keys{
		"I18N_COMMON_TOTAL_DCPOWER":                   "sunPower",
		"I18N_COMMON_PV_DAYILY_ENERGY_GENERATION":     "sunPowerDay",
		"I18N_COMMON_PV_TOTAL_ENERGY_GENERATION":      "sunPowerTotal",
		"I18N_COMMON_FEED_NETWORK_TOTAL_ACTIVE_POWER": "netFeedIn",
		"I18N_COMMON_DAILY_FEED_NETWORK_VOLUME":       "netFeedInDay",
		"I18N_COMMON_TOTAL_FEED_NETWORK_VOLUME":       "netFeedInTotal",
		"I18N_CONFIG_KEY_4060":                        "netPower",
		"I18N_COMMON_ENERGY_GET_FROM_GRID_DAILY":      "netPowerDay",
		"I18N_COMMON_TOTAL_ELECTRIC_GRID_GET_POWER":   "netPowerTotal",
		"I18N_COMMON_LOAD_TOTAL_ACTIVE_POWER":         "consumption",
		"I18N_CONFIG_KEY_1001188":                     "consumptionRate",
		"I18N_COMMON_AIR_TEM_INSIDE_MACHINE":          "inverterTemp",
		"I18N_COMMON_TOTAL_GRID_RUNNING_TIME":		   "totalRunningTime",
		"I18N_COMMON_DAILY_POWER_YIELD":               "todayEnergy",
		"I18N_COMMON_TOTAL_YIELD":                     "totalEnergy",
		"I18N_COMMON_TOTAL_ACTIVE_POWER":              "activePower",
		"I18N_COMMON_BUS_VOLTAGE":                     "busVoltage",
	}

	batteryKeys = ws.Keys{
		"I18N_COMMON_BATTERY_SOC":         "batteryLevel",
		"I18N_CONFIG_KEY_3907":            "batteryCharge",
		"I18N_CONFIG_KEY_3921":            "batteryDischarge",
		"I18N_COMMON_BATTARY_HEALTH":      "batteryHealth",
		"I18N_COMMON_BATTERY_TEMPERATURE": "batteryTemp",
	}
)
