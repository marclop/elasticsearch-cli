package poller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

type mockClient struct {
	content string
	fail    bool
	header  http.Header
}

func (c *mockClient) HandleCall(_, _, _ string) (*http.Response, error) {
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

func TestNewIndexPoller(t *testing.T) {
	channel := make(chan []string, 1)
	type args struct {
		client client
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
				make(chan bool, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewIndexPoller(tt.args.client, tt.args.c, tt.args.poll)
			got.controlChannel = tt.want.controlChannel
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIndexPoller() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexPoller_run(t *testing.T) {
	channel := make(chan []string, 1)
	type fields struct {
		client   client
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

func TestIndexPollerStart(t *testing.T) {
	type fields struct {
		client   client
		endpoint string
		pollRate time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"",
			fields{
				&mockClient{
					`yellow open   elastic      dBoWJXLBSRuumXa-a-QN1w   5   1          0            0       650b           650b
yellow open   found        oBcPStMpTD2BZtQ9j2ff3w   5   1          0            0       650b           650b
yellow open   wat          s0uzswacS2-jPJJgKb8r7w   5   1          0            0       650b           650b`,
					false,
					nil,
				},
				defaultPollingEndpoint,
				time.Duration(1 * time.Second),
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
			indicesChannel := make(chan []string, 1)
			controlChannel := make(chan bool, 1)
			w := &IndexPoller{
				client:         tt.fields.client,
				endpoint:       tt.fields.endpoint,
				channel:        indicesChannel,
				pollRate:       tt.fields.pollRate,
				controlChannel: controlChannel,
			}
			go w.Start()
			got := <-w.channel
			w.Stop()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IndexPoller.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}
