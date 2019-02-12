package app

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/chzyer/readline"
	"github.com/marclop/elasticsearch-cli/cli"
	"github.com/marclop/elasticsearch-cli/client"
	"github.com/marclop/elasticsearch-cli/poller"
)

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
				client.NewHTTP(defaultConfig, client.NewMock()),
				nil,
				channel,
				&poller.IndexPoller{},
			},
			&Application{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client:       client.NewHTTP(defaultConfig, client.NewMock()),
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
				client: client.NewHTTP(defaultConfig, client.NewMock(
					client.MockResponse{Response: http.Response{
						StatusCode: 200,
						Request:    &http.Request{Method: "GET"},
						Body:       client.NewStringBody(`{"health": "yellow", "status": "open", "index": "elastic", "pri": "1"}`),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}},
				)),
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
				client: client.NewHTTP(defaultConfig, client.NewMock(
					client.MockResponse{Response: http.Response{
						StatusCode: 200,
						Request:    &http.Request{Method: "GET", URL: new(url.URL)},
						Body:       client.NewStringBody(`{"health": "yellow", "status": "open", "index": "elastic", "pri": "1"}`),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}},
					client.MockResponse{Response: http.Response{
						StatusCode: 200,
						Request:    &http.Request{Method: "GET", URL: new(url.URL)},
						Body:       client.NewStringBody(`{"health": "yellow", "status": "open", "index": "elastic", "pri": "1"}`),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}},
				)),
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
				client: client.NewHTTP(defaultConfig, client.NewMock(
					client.MockResponse{Error: errors.New("someerror")},
				)),
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
				client: client.NewHTTP(defaultConfig, client.NewMock()),
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
			client.NewHTTP(defaultConfig, client.NewMock()),
		},
		{
			"host is not modified due invalid schema",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, client.NewMock()),
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
			client.NewHTTP(defaultConfig, client.NewMock()),
		},
		{
			"port is modified",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, client.NewMock()),
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
			client.NewHTTP(defaultConfig, client.NewMock()),
		},
		{
			"user is modified",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, client.NewMock()),
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
			client.NewHTTP(defaultConfig, client.NewMock()),
		},
		{
			"password is modified",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, client.NewMock()),
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
			client.NewHTTP(defaultConfig, client.NewMock()),
		},
		{
			"Verbose is modified",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, client.NewMock()),
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
			client.NewHTTP(defaultConfig, client.NewMock()),
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

func TestApplication_getClusterPrompt(t *testing.T) {
	type fields struct {
		config       *Config
		client       *client.HTTP
		formatFunc   Formatter
		output       io.Writer
		indexChannel chan []string
		parser       *cli.InputParser
		poller       Poller
		repl         *readline.Instance
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"When the cluster is green, returns the greenPrompt",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, client.NewMock(
					client.MockResponse{Response: http.Response{
						StatusCode: 200,
						Request:    &http.Request{Method: "GET"},
						Body: client.NewStringBody(`{
							"cluster_name": "myCluster",
							"status": "green",
							"timed_out": false,
							"number_of_nodes": 3,
							"number_of_data_nodes": 2,
							"active_primary_shards": 3,
							"active_shards": 6,
							"relocating_shards": 0,
							"initializing_shards": 0,
							"unassigned_shards": 0,
							"delayed_unassigned_shards": 0,
							"number_of_pending_tasks": 0,
							"number_of_in_flight_fetch": 0,
							"task_max_waiting_in_queue_millis": 0,
							"active_shards_percent_as_number": 100.0
						  }`),
						Header: http.Header{"Content-Type": []string{"application/json"}},
					}},
				)),
				repl:       &readline.Instance{},
				formatFunc: cli.Format,
				output:     &bytes.Buffer{},
			},
			GreenPrompt,
		},
		{
			"When the cluster is yellow, returns the yellowPrompt",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, client.NewMock(
					client.MockResponse{Response: http.Response{
						StatusCode: 200,
						Request:    &http.Request{Method: "GET"},
						Body: client.NewStringBody(`{
							"cluster_name": "myCluster",
							"status": "yellow",
							"timed_out": false,
							"number_of_nodes": 3,
							"number_of_data_nodes": 2,
							"active_primary_shards": 3,
							"active_shards": 6,
							"relocating_shards": 0,
							"initializing_shards": 0,
							"unassigned_shards": 0,
							"delayed_unassigned_shards": 0,
							"number_of_pending_tasks": 0,
							"number_of_in_flight_fetch": 0,
							"task_max_waiting_in_queue_millis": 0,
							"active_shards_percent_as_number": 100.0
						  }`),
						Header: http.Header{"Content-Type": []string{"application/json"}},
					}},
				)),
				repl:       &readline.Instance{},
				formatFunc: cli.Format,
				output:     &bytes.Buffer{},
			},
			YellowPrompt,
		},
		{
			"When the cluster is red, returns the redPrompt",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, client.NewMock(
					client.MockResponse{Response: http.Response{
						StatusCode: 200,
						Request:    &http.Request{Method: "GET"},
						Body: client.NewStringBody(`{
							"cluster_name": "myCluster",
							"status": "red",
							"timed_out": false,
							"number_of_nodes": 3,
							"number_of_data_nodes": 2,
							"active_primary_shards": 3,
							"active_shards": 6,
							"relocating_shards": 0,
							"initializing_shards": 0,
							"unassigned_shards": 0,
							"delayed_unassigned_shards": 0,
							"number_of_pending_tasks": 0,
							"number_of_in_flight_fetch": 0,
							"task_max_waiting_in_queue_millis": 0,
							"active_shards_percent_as_number": 100.0
						  }`),
						Header: http.Header{"Content-Type": []string{"application/json"}},
					}},
				)),
				repl:       &readline.Instance{},
				formatFunc: cli.Format,
				output:     &bytes.Buffer{},
			},
			RedPrompt,
		},
		{
			"When the request returns an unparsable body, returns the defaultPrompt",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, client.NewMock(
					client.MockResponse{Response: http.Response{
						StatusCode: 200,
						Request:    &http.Request{Method: "GET"},
						Body:       client.NewStringBody(`{"cluster_name": ,,,"myCluster",}`),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}},
				)),
				repl:       &readline.Instance{},
				formatFunc: cli.Format,
				output:     &bytes.Buffer{},
			},
			DefaultPrompt,
		},
		{
			"When the request returns an error, returns the defaultPrompt",
			fields{
				config: &Config{
					Verbose:      false,
					PollInterval: 10,
				},
				client: client.NewHTTP(defaultConfig, client.NewMock(
					client.MockResponse{Error: errors.New("someerror")},
				)),
				repl:       &readline.Instance{},
				formatFunc: cli.Format,
				output:     &bytes.Buffer{},
			},
			DefaultPrompt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				config:       tt.fields.config,
				client:       tt.fields.client,
				formatFunc:   tt.fields.formatFunc,
				output:       tt.fields.output,
				indexChannel: tt.fields.indexChannel,
				parser:       tt.fields.parser,
				poller:       tt.fields.poller,
				repl:         tt.fields.repl,
			}
			if got := app.getClusterPrompt(); got != tt.want {
				t.Errorf("Application.getClusterPrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}
