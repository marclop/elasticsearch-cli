package cli

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestNewJSONFormatter(t *testing.T) {
	type args struct {
		input *http.Response
	}
	tests := []struct {
		name string
		args args
		want Formatter
	}{
		{
			"NewJSONFormatterSucceds",
			args{
				&http.Response{},
			},
			&JSONFormatter{
				&http.Response{},
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJSONFormatter(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJSONFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIteractiveJSONFormatter(t *testing.T) {
	type args struct {
		input *http.Response
	}
	tests := []struct {
		name string
		args args
		want Formatter
	}{
		{
			"NewJInteractiveSONFormatterSucceds",
			args{
				&http.Response{},
			},
			&JSONFormatter{
				&http.Response{},
				true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIteractiveJSONFormatter(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIteractiveJSONFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONFormatter_Format(t *testing.T) {
	type fields struct {
		input       *http.Response
		interactive bool
	}
	type args struct {
		verbose bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"FormatSucceeds",
			fields{
				&http.Response{
					Body: ioutil.NopCloser(strings.NewReader("MyPlainTextInput")),
					Request: &http.Request{
						Method: "GET",
					},
				},
				false,
			},
			args{
				false,
			},
		},
		{
			"FormatSucceedsWithJSONResponse",
			fields{
				&http.Response{
					Body: ioutil.NopCloser(strings.NewReader(`{"a":"b"}`)),
					Request: &http.Request{
						Method: "GET",
					},
				},
				false,
			},
			args{
				false,
			},
		},
		{
			"FormatSucceedsWithHEADMethod",
			fields{
				&http.Response{
					Body: ioutil.NopCloser(strings.NewReader("")),
					Request: &http.Request{
						Method: "HEAD",
						Header: http.Header{
							"Content-Type": []string{"application/json"},
						},
					},
				},
				false,
			},
			args{
				false,
			},
		},
		{
			"FormatSucceedsWithVerboseOn",
			fields{
				&http.Response{
					Body: ioutil.NopCloser(strings.NewReader("")),
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
			},
			args{
				true,
			},
		},
		{
			"FormatSucceedsWithInteractive",
			fields{
				&http.Response{
					Body: ioutil.NopCloser(strings.NewReader("")),
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
			},
			args{
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &JSONFormatter{
				input:       tt.fields.input,
				interactive: tt.fields.interactive,
			}
			f.Format(tt.args.verbose)
		})
	}
}
