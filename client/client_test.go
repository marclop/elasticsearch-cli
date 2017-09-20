package client

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/marclop/elasticsearch-cli/utils"
)

type emptyMockCaller struct{}

func (c *emptyMockCaller) Do(*http.Request) (*http.Response, error) {
	var err error
	return &http.Response{}, err
}

func TestNewClient(t *testing.T) {
	type args struct {
		config *Config
		client HTTPCallerInterface
	}
	tests := []struct {
		name string
		args args
		want *HTTP
	}{
		{
			"NewClientHasMockClientInjected",
			args{
				config: &Config{&hostPort{"http://localhost", 9200}, "", "", time.Duration(10), nil},
				client: &emptyMockCaller{},
			},
			NewHTTP(
				&Config{&hostPort{"http://localhost", 9200}, "", "", time.Duration(10), nil},
				&emptyMockCaller{},
			),
		},
		{
			"NewClientHaNoInjections",
			args{
				config: &Config{&hostPort{"http://localhost", 9200}, "", "", time.Duration(10), nil},
				client: nil,
			},
			NewHTTP(
				&Config{&hostPort{"http://localhost", 9200}, "", "", time.Duration(10), nil},
				nil,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHTTP(tt.args.config, tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_HandleCall(t *testing.T) {
	type fields struct {
		config *Config
		caller HTTPCallerInterface
	}
	type args struct {
		method string
		url    string
		body   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			"HandleCallHTTPByEmptyMockCaller",
			fields{
				&Config{&hostPort{"http://localhost", 9200}, "", "", time.Duration(10), nil},
				&emptyMockCaller{},
			},
			args{
				"GET",
				"",
				"",
			},
			&http.Response{},
			false,
		},
		{
			"HandleCallHTTPSByEmptyMockCaller",
			fields{
				&Config{&hostPort{"https://localhost", 9200}, "", "", time.Duration(10), nil},
				&emptyMockCaller{},
			},
			args{
				"GET",
				"",
				"",
			},
			&http.Response{},
			false,
		},
		{
			"HandleCallWithBodyByEmptyMockCaller",
			fields{
				&Config{&hostPort{"http://localhost", 9200}, "", "", time.Duration(10), nil},
				&emptyMockCaller{},
			},
			args{
				"GET",
				"",
				"{\"hello\":true}",
			},
			&http.Response{},
			false,
		},
		{
			"HandleCallWithHeadersByEmptyMockCaller",
			fields{
				&Config{&hostPort{"http://localhost", 9200}, "", "", time.Duration(10), map[string]string{"Content-Type": "application/json"}},
				&emptyMockCaller{},
			},
			args{
				"GET",
				"",
				"",
			},
			&http.Response{},
			false,
		},
		{
			"HandleCallWithAuthAndHeadersByEmptyMockCaller",
			fields{
				&Config{&hostPort{"http://localhost", 9200}, "marc", "marc", time.Duration(10), map[string]string{"Content-Type": "application/json"}},
				&emptyMockCaller{},
			},
			args{
				"GET",
				"",
				"",
			},
			&http.Response{},
			false,
		},
		{
			"HandleCallWithInvalidMethodByEmptyMockCaller",
			fields{
				&Config{&hostPort{"http://localhost", 9200}, "", "", time.Duration(10), nil},
				&emptyMockCaller{},
			},
			args{
				"   ",
				"",
				"",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &HTTP{
				Config: tt.fields.config,
				caller: tt.fields.caller,
			}
			got, err := c.HandleCall(tt.args.method, tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.HandleCall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.HandleCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_createRequest(t *testing.T) {
	type fields struct {
		config *Config
		caller HTTPCallerInterface
	}
	type args struct {
		method string
		url    string
		body   io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			"createRequestWithCorrectMethod",
			fields{
				&Config{&hostPort{"http://localhost", 9200}, "", "", time.Duration(10), nil},
				&emptyMockCaller{},
			},
			args{
				"GET",
				"",
				nil,
			},
			&http.Request{
				Method:     "GET",
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     make(map[string][]string),
				Body:       nil,
				Host:       utils.ConcatStrings(),
			},
			false,
		},
		{
			"createRequestWithIncorrectMethod",
			fields{
				&Config{&hostPort{"http://localhost", 9200}, "", "", time.Duration(10), nil},
				&emptyMockCaller{},
			},
			args{
				"INVALID METHOD",
				"",
				nil,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &HTTP{
				Config: tt.fields.config,
				caller: tt.fields.caller,
			}
			got, err := c.createRequest(tt.args.method, tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.createRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				tt.want.URL, _ = url.Parse(tt.args.url)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.createRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
