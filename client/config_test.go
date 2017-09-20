package client

import (
	"reflect"
	"testing"
	"time"
)

func TestNewClientConfig(t *testing.T) {
	type args struct {
		host    string
		port    int
		user    string
		pass    string
		timeout int
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			"NewClientConfigSucceeds",
			args{
				"http://localhost",
				9200,
				"",
				"",
				10,
			},
			&Config{
				&hostPort{"http://localhost", 9200},
				"",
				"",
				time.Duration(10 * time.Second),
				defaultClientHeaders,
			},
			false,
		},
		{
			"NewClientConfigFailsDueBadHTTPSchema",
			args{
				"yolo://localhost",
				9200,
				"",
				"",
				10,
			},
			nil,
			true,
		},
		{
			"NewConfigHostSucceedsWhenURLIncludesPort",
			args{
				"http://localhost:9201",
				9200,
				"",
				"",
				10,
			},
			&Config{
				&hostPort{"http://localhost", 9201},
				"",
				"",
				time.Duration(10 * time.Second),
				defaultClientHeaders,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClientConfig(tt.args.host, tt.args.port, tt.args.user, tt.args.pass, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateSchema(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"ValidateSchemaSucceedsWithHTTP",
			args{
				"http://localhost",
			},
			false,
		},
		{
			"ValidateSchemaSucceedsWithHTTPS",
			args{
				"https://localhost",
			},
			false,
		},
		{
			"ValidateSchemaFailsWhenSchemaIsInvalid",
			args{
				"wat://localhost",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateSchema(tt.args.host); (err != nil) != tt.wantErr {
				t.Errorf("validateSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_HTTPAdress(t *testing.T) {
	type fields struct {
		hostPort *hostPort
		user     string
		pass     string
		timeout  time.Duration
		headers  map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"HTTPAdressIsValid",
			fields{
				&hostPort{
					"http://localhost",
					9200,
				},
				"",
				"",
				time.Duration(10 * time.Second),
				nil,
			},
			"http://localhost:9200",
		},
		{
			"HTTPAdressIsValidWithNonDefaultPort",
			fields{
				&hostPort{
					"http://localhost",
					9230,
				},
				"",
				"",
				time.Duration(10 * time.Second),
				nil,
			},
			"http://localhost:9230",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				HostPort: tt.fields.hostPort,
				User:     tt.fields.user,
				Pass:     tt.fields.pass,
				Timeout:  tt.fields.timeout,
				headers:  tt.fields.headers,
			}
			if got := c.HTTPAdress(); got != tt.want {
				t.Errorf("Config.HTTPAdress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetHost(t *testing.T) {
	type fields struct {
		hostPort *hostPort
		user     string
		pass     string
		timeout  time.Duration
		headers  map[string]string
	}
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"SetHostSucceeds",
			fields{
				&hostPort{
					"http://localhost",
					9200,
				},
				"",
				"",
				time.Duration(10 * time.Second),
				nil,
			},
			args{
				"https://localhost",
			},
			false,
		},
		{
			"SetHosFailsWhenHTTPSchemaIsInvalid",
			fields{
				&hostPort{
					"http://localhost",
					9200,
				},
				"",
				"",
				time.Duration(10 * time.Second),
				nil,
			},
			args{
				"invalid://localhost",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				HostPort: tt.fields.hostPort,
				User:     tt.fields.user,
				Pass:     tt.fields.pass,
				Timeout:  tt.fields.timeout,
				headers:  tt.fields.headers,
			}
			if err := c.SetHost(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Config.SetHost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && c.HostPort.Host != tt.args.value {
				t.Errorf("Config.SetHost() = %s, want %s", c.HostPort.Host, tt.args.value)
			}
		})
	}
}
