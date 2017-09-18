package cli

import (
	"reflect"
	"testing"
)

func TestNewInputParser(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name    string
		args    args
		want    *InputParser
		wantErr bool
	}{
		{
			"NewParserSucceedsWhenProvidedCorrectMethod",
			args{
				[]string{
					"GET",
					"/",
					"",
				},
			},
			&InputParser{
				"GET",
				"/",
				"",
			},
			false,
		},
		{
			"NewParserSucceedsWhenURLIsBlank",
			args{
				[]string{
					"GET",
					"",
					"",
				},
			},
			&InputParser{
				"GET",
				"/",
				"",
			},
			false,
		},
		{
			"NewParserFailsWhenMethodIsBlank",
			args{
				[]string{
					"",
					"/",
					"",
				},
			},
			nil,
			true,
		},
		{
			"NewParserFailsWhenMethodIsInvalid",
			args{
				[]string{
					"WAT",
					"/",
					"",
				},
			},
			nil,
			true,
		},
		{
			"NewParserSucceedsWhenIsInteractive",
			args{
				[]string{},
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewInputParser(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Validate(t *testing.T) {
	type fields struct {
		Method string
		URL    string
		Body   string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"ValidateSucceeds",
			fields{
				"GET",
				"/",
				"",
			},
			false,
		},
		{
			"ValidateSucceedsWhenMethodIsEmpty",
			fields{
				"",
				"/",
				"",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InputParser{
				Method: tt.fields.Method,
				URL:    tt.fields.URL,
				Body:   tt.fields.Body,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Parser.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
