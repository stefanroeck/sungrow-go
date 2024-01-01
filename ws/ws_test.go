package ws

import (
	"reflect"
	"testing"
)

var (
	testHost   = "192.168.2.100"
	testPort = 8082
	testPath = "/ws/home/overview"
)

func TestNewWS(t *testing.T) {
	type fields struct {
		host string
		port int
		path string
	}

	tests := []struct {
		name    string
		fields  fields
		want    *WS
		wantErr error
	}{
		{
			name:    "1",
			fields:  fields{testHost, testPort, testPath},
			want:    &WS{host: testHost, port: testPort, path: testPath},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWS(tt.fields.host, tt.fields.port, tt.fields.path)

			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewWS() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWS() = %v, want %v", got, tt.want)
			}
		})
	}
}
