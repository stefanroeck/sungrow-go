package params

import (
	"testing"
)

func TestParse_AllFlags(t *testing.T) {

	args := []string{
		"-protocol=wss",
		"-host=inverter.local",
		"-port=9090",
		"-user=admin",
		"-password=secret",
		"-path=/custom/path",
		"-mqtt.server=ssl://mqtt.broker:8883",
		"-mqtt.user=mqttuser",
		"-mqtt.password=mqttpass",
		"-mqtt.clientId=client123",
		"-mqtt.topic=topic/test",
		"-mqtt.skipSSLVerify=false",
		"-sleep=42",
	}

	cfg := ParseFlags(args)

	// inverter params
	if got, want := cfg.InverterParams.Protocol, "wss"; got != want {
		t.Errorf("Protocol = %q; want %q", got, want)
	}
	if got, want := cfg.InverterParams.Host, "inverter.local"; got != want {
		t.Errorf("Host = %q; want %q", got, want)
	}
	if got, want := cfg.InverterParams.Port, 9090; got != want {
		t.Errorf("Port = %d; want %d", got, want)
	}
	if got, want := cfg.InverterParams.User, "admin"; got != want {
		t.Errorf("User = %q; want %q", got, want)
	}
	if got, want := cfg.InverterParams.Password, "secret"; got != want {
		t.Errorf("Password = %q; want %q", got, want)
	}
	if got, want := cfg.InverterParams.Path, "/custom/path"; got != want {
		t.Errorf("Path = %q; want %q", got, want)
	}

	// mqtt params
	if got, want := cfg.MqttParams.Server, "ssl://mqtt.broker:8883"; got != want {
		t.Errorf("MQTT Server = %q; want %q", got, want)
	}
	if got, want := cfg.MqttParams.User, "mqttuser"; got != want {
		t.Errorf("MQTT User = %q; want %q", got, want)
	}
	if got, want := cfg.MqttParams.Password, "mqttpass"; got != want {
		t.Errorf("MQTT Password = %q; want %q", got, want)
	}
	if got, want := cfg.MqttParams.ClientId, "client123"; got != want {
		t.Errorf("MQTT ClientId = %q; want %q", got, want)
	}
	if got, want := cfg.MqttParams.Topic, "topic/test"; got != want {
		t.Errorf("MQTT Topic = %q; want %q", got, want)
	}
	if got, want := cfg.MqttParams.SkipSSLVerify, false; got != want {
		t.Errorf("MQTT SkipSSLVerify = %v; want %v", got, want)
	}

	// sleep time
	if got, want := cfg.SleepBetweenCallsSecs, 42; got != want {
		t.Errorf("SleepBetweenCallsSecs = %d; want %d", got, want)
	}
}
