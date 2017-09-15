package poller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/marclop/elasticsearch-cli/client"
)

type mockClient struct {
	content string
	fail    bool
	header  http.Header
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
		},
	}, err
}

func (c *mockClient) SetHost(string) error {
	return nil
}

func (c *mockClient) SetPort(int)    {}
func (c *mockClient) SetUser(string) {}
func (c *mockClient) SetPass(string) {}

func TestNewIndexPoller(t *testing.T) {
	channel := make(chan []string, 1)
	type args struct {
		client client.Client
		c      chan []string
		poll   int
	}
	tests := []struct {
		name string
		args args
		want *IndexPoller
	}{
		{
			"NewIndexPollerSucceeds",
			args{
				&mockClient{},
				channel,
				10,
			},
			&IndexPoller{
				&mockClient{},
				defaultPollingEndpoint,
				channel,
				time.Duration(10) * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIndexPoller(tt.args.client, tt.args.c, tt.args.poll); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIndexPoller() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexPoller_run(t *testing.T) {
	channel := make(chan []string, 1)
	type fields struct {
		client   client.Client
		endpoint string
		channel  chan []string
		pollRate time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"RunPollerSucceeds",
			fields{
				&mockClient{
					`yellow open   elastic      dBoWJXLBSRuumXa-a-QN1w   5   1          0            0       650b           650b
yellow open   found        oBcPStMpTD2BZtQ9j2ff3w   5   1          0            0       650b           650b
yellow open   wat          s0uzswacS2-jPJJgKb8r7w   5   1          0            0       650b           650b`,
					false,
					nil,
				},
				defaultPollingEndpoint,
				channel,
				time.Duration(10 * time.Second),
			},
			[]string{
				"elastic",
				"found",
				"wat",
			},
		},
		{
			"RunPollerSucceedsFails",
			fields{
				&mockClient{
					"",
					true,
					nil,
				},
				defaultPollingEndpoint,
				channel,
				time.Duration(10 * time.Second),
			},
			nil,
		},
		{
			"RunPollerSucceeds",
			fields{
				&mockClient{
					`[
  {
    "health": "yellow",
    "status": "open",
    "index": "elastic",
    "pri": "1",
    "rep": "1",
    "docs.count": "150000",
    "docs.deleted": "0",
    "store.size": "31.7mb",
    "pri.store.size": "31.7mb"
  },
  {
    "health": "yellow",
    "status": "open",
    "index": "found",
    "pri": "1",
    "rep": "1",
    "docs.count": "150000",
    "docs.deleted": "0",
    "store.size": "31.7mb",
    "pri.store.size": "31.7mb"
  },
  {
    "health": "yellow",
    "status": "open",
    "index": "wat",
    "pri": "1",
    "rep": "1",
    "docs.count": "150000",
    "docs.deleted": "0",
    "store.size": "31.7mb",
    "pri.store.size": "31.7mb"
  }
]`,
					false,
					http.Header{
						"Content-Type": []string{"application/json"},
					},
				},
				defaultPollingEndpoint,
				channel,
				time.Duration(10 * time.Second),
			},
			[]string{
				"elastic",
				"found",
				"wat",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &IndexPoller{
				client:   tt.fields.client,
				endpoint: tt.fields.endpoint,
				channel:  tt.fields.channel,
				pollRate: tt.fields.pollRate,
			}
			if got := w.run(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IndexPoller.run() = %v, want %v", got, tt.want)
			}
		})
	}
}
