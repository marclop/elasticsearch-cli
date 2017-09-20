package cli

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func encodeData(p string) string {
	var v bytes.Buffer
	json.Indent(&v, []byte(p), "", "  ")
	v.WriteString("\n")
	return v.String()
}

func TestFormat(t *testing.T) {
	type args struct {
		input       *http.Response
		verbose     bool
		interactive bool
		writer      *bytes.Buffer
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"FormatSucceeds",
			args{
				&http.Response{
					Body:   ioutil.NopCloser(strings.NewReader("MyPlainTextInput")),
					Status: "200 OK",
					Request: &http.Request{
						Method: "GET",
					},
				},
				false,
				false,
				&bytes.Buffer{},
			},
			`MyPlainTextInput
`,
		},
		{
			"FormatSucceedsWithJSONResponse",
			args{
				&http.Response{
					Body:   ioutil.NopCloser(strings.NewReader(`{"a":"b"}`)),
					Status: "200 OK",
					Request: &http.Request{
						Method: "GET",
					},
				},
				false,
				false,
				&bytes.Buffer{},
			},
			encodeData(`{"a":"b"}`),
		},
		{
			"FormatSucceedsWithHEADMethod",
			args{
				&http.Response{
					Body:   ioutil.NopCloser(strings.NewReader("")),
					Status: "200 OK",
					Request: &http.Request{
						Method: "HEAD",
						Header: http.Header{
							"Content-Type": []string{"application/json"},
						},
						URL: &url.URL{
							Path: "/",
						},
					},
				},
				false,
				false,
				&bytes.Buffer{},
			},
			`Response:     200 OK
Content-Type: application/json

`,
		},
		{
			"FormatSucceedsWithVerboseOn",
			args{
				&http.Response{
					Body:   ioutil.NopCloser(strings.NewReader("")),
					Status: "200 OK",
					Request: &http.Request{
						Method: "HEAD",
						Header: http.Header{
							"Content-Type": []string{"application/json"},
						},
						URL: &url.URL{
							Path: "/",
						},
					},
				},
				true,
				false,
				&bytes.Buffer{},
			},
			`Method:       HEAD
URL:          /
Response:     200 OK
Content-Type: application/json

`,
		},
		{
			"FormatSucceedsWithInteractive",
			args{
				&http.Response{
					Body:   ioutil.NopCloser(strings.NewReader("")),
					Status: "200 OK",
					Request: &http.Request{
						Method: "HEAD",
						Header: http.Header{
							"Content-Type": []string{"application/json"},
						},
						URL: &url.URL{
							Path: "/",
						},
					},
				},
				false,
				true,
				&bytes.Buffer{},
			},
			`Method:       HEAD
URL:          /
Response:     200 OK
Content-Type: application/json

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Format(tt.args.input, tt.args.verbose, tt.args.interactive, tt.args.writer)
			if tt.args.writer.String() != tt.want {
				t.Errorf("Format() = %v, want = %v", tt.args.writer.String(), tt.want)
			}
		})
	}
}
