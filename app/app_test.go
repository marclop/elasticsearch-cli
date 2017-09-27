package app

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/chzyer/readline"
	"github.com/marclop/elasticsearch-cli/cli"
	"github.com/marclop/elasticsearch-cli/client"
	"github.com/marclop/elasticsearch-cli/poller"
)

type mockCaller struct {
	content string
	fail    bool
	header  http.Header
	props   map[string]interface{}
}

func (c *mockCaller) Do(*http.Request) (*http.Response, error) {
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

var defaultConfig = func(c *client.Config, _ error) *client.Config { return c }(
	client.NewClientConfig("http://localhost", 9200, "user", "pass", 10, false),
)

func TestInitialize(t *testing.T) {
	channel := make(chan []string, 1)
	type args struct {
		config *Config
		client *client.HTTP
		f      Formatter
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
					Verbose:      false,
					PollInterval: 10,
				},
				client.NewHTTP(defaultConfig, &mockCaller{}),
				nil,
				channel,
				&poller.IndexPoller{},
			},
			&Application{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client:       client.NewHTTP(defaultConfig, &mockCaller{}),
				indexChannel: channel,
				poller:       &poller.IndexPoller{},
				formatFunc:   nil,
				output:       nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initialize(tt.args.config, tt.args.client, nil, tt.args.c, tt.args.w, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_HandleCli(t *testing.T) {
	type fields struct {
		config       *Config
		client       *client.HTTP
		format       Formatter
		indexChannel chan []string
		poller       Poller
		repl         *readline.Instance
		output       io.Writer
	}
	type args struct {
		input []string
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
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, &mockCaller{
					content: `{"health": "yellow", "status": "open", "index": "elastic", "pri": "1"}`,
					fail:    false,
					header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}),
				format: cli.Format,
				output: &bytes.Buffer{},
			},
			args{
				[]string{
					"GET",
					"",
					"",
				},
			},
			false,
		},
		{
			"HandleCliInteractiveSucceeds",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, &mockCaller{
					content: `{"health": "yellow", "status": "open", "index": "elastic", "pri": "1"}`,
					fail:    false,
					header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}),
				repl:   &readline.Instance{},
				format: cli.Format,
				output: &bytes.Buffer{},
			},
			args{
				[]string{
					"GET",
					"/",
					"",
				},
			},
			false,
		},
		{
			"HandleCliFailsDueHandleCallFailing",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, &mockCaller{
					content: `{"health": "yellow", "status": "open", "index": "elastic", "pri": "1"}`,
					fail:    true,
					header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}),
				format: cli.Format,
				output: &bytes.Buffer{},
			},
			args{
				[]string{
					"GET",
					"/",
					"",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				config:       tt.fields.config,
				client:       tt.fields.client,
				formatFunc:   tt.fields.format,
				indexChannel: tt.fields.indexChannel,
				poller:       tt.fields.poller,
				repl:         tt.fields.repl,
				output:       tt.fields.output,
			}
			if err := app.HandleCli(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Application.HandleCli() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_doSetCommands(t *testing.T) {
	type fields struct {
		config       *Config
		client       *client.HTTP
		format       Formatter
		indexChannel chan []string
		parser       *cli.InputParser
		poller       Poller
		repl         *readline.Instance
		output       io.Writer
	}
	type args struct {
		lineSliced []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *client.HTTP
	}{
		{
			"host is modified",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, &mockCaller{
					props: make(map[string]interface{}, 1),
				}),
				repl:   &readline.Instance{},
				format: cli.Format,
				output: &bytes.Buffer{},
			},
			args{
				[]string{
					"set",
					"host",
					"http://localhost",
				},
			},
			client.NewHTTP(defaultConfig, &mockCaller{
				props: map[string]interface{}{
					"host": "http://localhost",
				},
			}),
		},
		{
			"host is not modified due invalid schema",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, &mockCaller{
					props: make(map[string]interface{}, 1),
				}),
				repl:   &readline.Instance{},
				format: cli.Format,
				output: &bytes.Buffer{},
			},
			args{
				[]string{
					"set",
					"host",
					"INVALID://localhost",
				},
			},
			client.NewHTTP(defaultConfig, &mockCaller{
				props: map[string]interface{}{
					"host": "INVALID://localhost",
				},
			}),
		},
		{
			"port is modified",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, &mockCaller{
					props: make(map[string]interface{}, 1),
				}),
				repl:   &readline.Instance{},
				format: cli.Format,
				output: &bytes.Buffer{},
			},
			args{
				[]string{
					"set",
					"port",
					"9201",
				},
			},
			client.NewHTTP(defaultConfig, &mockCaller{
				props: map[string]interface{}{
					"port": 9201,
				},
			}),
		},
		{
			"user is modified",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, &mockCaller{
					props: make(map[string]interface{}, 1),
				}),
				repl:   &readline.Instance{},
				format: cli.Format,
				output: &bytes.Buffer{},
			},
			args{
				[]string{
					"set",
					"user",
					"elastic",
				},
			},
			client.NewHTTP(defaultConfig, &mockCaller{
				props: map[string]interface{}{
					"user": "elastic",
				},
			}),
		},
		{
			"password is modified",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, &mockCaller{
					props: make(map[string]interface{}, 1),
				}),
				repl:   &readline.Instance{},
				format: cli.Format,
				output: &bytes.Buffer{},
			},
			args{
				[]string{
					"set",
					"pass",
					"elastic",
				},
			},
			client.NewHTTP(defaultConfig, &mockCaller{
				props: map[string]interface{}{
					"pass": "elastic",
				},
			}),
		},
		{
			"Verbose is modified",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, &mockCaller{
					props: make(map[string]interface{}, 1),
				}),
				repl:   &readline.Instance{},
				format: cli.Format,
				output: &bytes.Buffer{},
			},
			args{
				[]string{
					"set",
					"verbose",
				},
			},
			client.NewHTTP(defaultConfig, &mockCaller{
				props: map[string]interface{}{
					"verbose": "elastic",
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				config:       tt.fields.config,
				client:       tt.fields.client,
				formatFunc:   tt.fields.format,
				indexChannel: tt.fields.indexChannel,
				parser:       tt.fields.parser,
				poller:       tt.fields.poller,
				repl:         tt.fields.repl,
				output:       tt.fields.output,
			}
			app.doSetCommands(tt.args.lineSliced)
			if !reflect.DeepEqual(app.client.Config, tt.want.Config) {
				t.Errorf("app.client = %v, want %v", app.client.Config, tt.want.Config)
			}
		})
	}
}
