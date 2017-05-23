package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/chzyer/readline"
	"github.com/elastic/elasticsearch-cli/cli"
	"github.com/elastic/elasticsearch-cli/client"
	"github.com/elastic/elasticsearch-cli/poller"
)

type mockClient struct {
	content string
	fail    bool
	header  http.Header
	props   map[string]interface{}
}

func (c *mockClient) HandleCall(string, string, string) (*http.Response, error) {
	var err error
	if c.fail {
		err = fmt.Errorf("fail")
	}

	return &http.Response{
		Header: c.header,
		Body:   ioutil.NopCloser(strings.NewReader(c.content)),
		Request: &http.Request{
			Method: "GET",
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			URL: &url.URL{
				Path: "/",
			},
		},
	}, err
}

func (c *mockClient) SetHost(host string) error {
	c.props["host"] = host
	return nil
}

func (c *mockClient) SetPort(port int) {
	c.props["port"] = port
}
func (c *mockClient) SetUser(user string) {
	c.props["user"] = user
}
func (c *mockClient) SetPass(pass string) {
	c.props["pass"] = pass
}

func TestInit(t *testing.T) {
	channel := make(chan []string, 1)
	type args struct {
		config *Config
		client client.Client
		parser Parser
		c      chan []string
		w      Poller
	}
	tests := []struct {
		name string
		args args
		want *Application
	}{
		{
			"InitApplicationSucceeds",
			args{
				&Config{
					false,
					10,
				},
				&mockClient{},
				&cli.InputParser{},
				channel,
				&poller.IndexPoller{},
			},
			&Application{
				config: &Config{
					false,
					10,
				},
				client:       &mockClient{},
				parser:       &cli.InputParser{},
				indexChannel: channel,
				poller:       &poller.IndexPoller{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(tt.args.config, tt.args.client, tt.args.parser, tt.args.c, tt.args.w); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_HandleCli(t *testing.T) {
	type fields struct {
		config       *Config
		client       client.Client
		formatter    cli.Formatter
		indexChannel chan []string
		parser       Parser
		poller       Poller
		repl         *readline.Instance
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
		wantErr bool
	}{
		{
			"HandleCliSucceeds",
			fields{
				config: &Config{
					false,
					10,
				},
				client: &mockClient{
					content: `{"health": "yellow", "status": "open", "index": "elastic", "pri": "1"}`,
					fail:    false,
					header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				},
			},
			args{
				"GET",
				"",
				"",
			},
			false,
		},
		{
			"HandleCliInteractiveSucceeds",
			fields{
				config: &Config{
					false,
					10,
				},
				client: &mockClient{
					content: `{"health": "yellow", "status": "open", "index": "elastic", "pri": "1"}`,
					fail:    false,
					header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				},
				repl: &readline.Instance{},
			},
			args{
				"GET",
				"/",
				"",
			},
			false,
		},
		{
			"HandleCliFailsDueHandleCallFailing",
			fields{
				config: &Config{
					false,
					10,
				},
				client: &mockClient{
					content: `{"health": "yellow", "status": "open", "index": "elastic", "pri": "1"}`,
					fail:    true,
					header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				},
			},
			args{
				"GET",
				"/",
				"",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				config:       tt.fields.config,
				client:       tt.fields.client,
				formatter:    tt.fields.formatter,
				indexChannel: tt.fields.indexChannel,
				parser:       tt.fields.parser,
				poller:       tt.fields.poller,
				repl:         tt.fields.repl,
			}
			if err := app.HandleCli(tt.args.method, tt.args.url, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Application.HandleCli() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_doSetCommands(t *testing.T) {
	type fields struct {
		config       *Config
		client       client.Client
		formatter    cli.Formatter
		indexChannel chan []string
		parser       Parser
		poller       Poller
		repl         *readline.Instance
	}
	type args struct {
		lineSliced []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   client.Client
	}{
		{
			"SetHostIsCalled",
			fields{
				config: &Config{
					false,
					10,
				},
				client: &mockClient{
					props: make(map[string]interface{}, 1),
				},
				repl: &readline.Instance{},
			},
			args{
				[]string{
					"set",
					"host",
					"http://localhost",
				},
			},
			&mockClient{
				props: map[string]interface{}{
					"host": "http://localhost",
				},
			},
		},
		{
			"SetPortIsCalled",
			fields{
				config: &Config{
					false,
					10,
				},
				client: &mockClient{
					props: make(map[string]interface{}, 1),
				},
				repl: &readline.Instance{},
			},
			args{
				[]string{
					"set",
					"port",
					"9201",
				},
			},
			&mockClient{
				props: map[string]interface{}{
					"port": 9201,
				},
			},
		},
		{
			"SetUserIsCalled",
			fields{
				config: &Config{
					false,
					10,
				},
				client: &mockClient{
					props: make(map[string]interface{}, 1),
				},
				repl: &readline.Instance{},
			},
			args{
				[]string{
					"set",
					"user",
					"elastic",
				},
			},
			&mockClient{
				props: map[string]interface{}{
					"user": "elastic",
				},
			},
		},
		{
			"SetPassIsCalled",
			fields{
				config: &Config{
					false,
					10,
				},
				client: &mockClient{
					props: make(map[string]interface{}, 1),
				},
				repl: &readline.Instance{},
			},
			args{
				[]string{
					"set",
					"pass",
					"elastic",
				},
			},
			&mockClient{
				props: map[string]interface{}{
					"pass": "elastic",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				config:       tt.fields.config,
				client:       tt.fields.client,
				formatter:    tt.fields.formatter,
				indexChannel: tt.fields.indexChannel,
				parser:       tt.fields.parser,
				poller:       tt.fields.poller,
				repl:         tt.fields.repl,
			}
			app.doSetCommands(tt.args.lineSliced)
			if !reflect.DeepEqual(app.client, tt.want) {
				t.Errorf("app.client = %v, want %v", app.client, tt.want)
			}
		})
	}
}
