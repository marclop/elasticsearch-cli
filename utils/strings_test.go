package utils

import (
	"io"
	"strings"
	"testing"
)

func TestReaderToString(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"NewReaderToStringSucceeds",
			args{
				strings.NewReader("MyTestString"),
			},
			"MyTestString",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReaderToString(tt.args.reader); got != tt.want {
				t.Errorf("ReaderToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringInSlice(t *testing.T) {
	type args struct {
		a    string
		list []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"StringInSliceIsFound",
			args{
				"MyFoundString",
				[]string{
					"My",
					"String",
					"MyFoundString",
				},
			},
			true,
		},
		{
			"StringInSliceIsNotFound",
			args{
				"MyMissingString",
				[]string{
					"My",
					"String",
					"MyString",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInSlice(tt.args.a, tt.args.list); got != tt.want {
				t.Errorf("StringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcatStrings(t *testing.T) {
	type args struct {
		strs []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"ConcatStringsSucceeds",
			args{
				[]string{
					"a",
					"b",
					"/",
					"c",
				},
			},
			"ab/c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConcatStrings(tt.args.strs...); got != tt.want {
				t.Errorf("ConcatStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}
