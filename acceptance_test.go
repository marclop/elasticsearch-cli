// +build acceptance

package main

import (
	"encoding/json"
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/elastic/elasticsearch-cli/utils"
)

type Response struct {
	Acknowledged bool   `json:"acknowledged"`
	Status       int    `json:"status"`
	Tagline      string `json:"tagline"`
}

func TestElasticsearchCli_NonInteractive(t *testing.T) {
	binary := "bin/elasticsearch-cli"
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		want    *Response
	}{
		{
			"GetRootSucceeds",
			[]string{
				"GET",
				"/",
				"",
			},
			false,
			nil,
		},
		{
			"RunInvalidCommandFails",
			[]string{
				"asda",
				"/",
				"",
			},
			true,
			nil,
		},
		{
			"CreateIndexWorks",
			[]string{
				"PUT",
				"MyTestIndex",
				"",
			},
			false,
			&Response{},
		},
		{
			"DELETEIndexWorks",
			[]string{
				"DELETE",
				"MyTestIndex",
				"",
			},
			false,
			&Response{},
		},
		{
			"CreateIndexWithBodyWorks",
			[]string{
				"PUT",
				"MySettingsTestIndex",
				`{"settings": {"index": {"number_of_shards" : 1, "number_of_replicas" : 0} }}`,
			},
			false,
			&Response{},
		},
		{
			"GetIndexSettingsucceeds",
			[]string{
				"GET",
				"MySettingsTestIndex/_settings",
				"",
			},
			false,
			nil,
		},
		{
			"DELETEIndexWorks",
			[]string{
				"DELETE",
				"MySettingsTestIndex",
				"",
			},
			false,
			&Response{},
		},
	}

	for _, tt := range tests {
		out, err := exec.Command(binary, tt.args...).Output()
		if (err != nil) != tt.wantErr {
			t.Errorf("Command %s %v error = %v, wantErr %v", binary, tt.args, err, tt.wantErr)
		}
		t.Logf("%s result: \n%s", tt.name, string(out))

		if tt.want != nil {
			json.Unmarshal(out, tt.want)

			if tt.want.Status == 400 {
				t.Error("Create Index failed: Command", binary, tt.args)
			}
			if !tt.want.Acknowledged {
				t.Error("Operation failed: Command", binary, tt.args)
			}
		}
	}
}

func TestElasticsearchCli_Interactive(t *testing.T) {
	binary := "bin/elasticsearch-cli"
	tests := []struct {
		name    string
		lines   []string
		wantErr bool
		want    []string
	}{
		{
			"GetRootSucceeds",
			[]string{
				"GET /",
				"exit",
			},
			false,
			[]string{
				"You Know, for Search",
			},
		},
		{
			"GetRootSucceeds",
			[]string{
				"PUT a",
				"GET a",
				"DELETE a",
				"exit",
			},
			false,
			[]string{
				`"a": {`,
				`"acknowledged": true`,
			},
		},
	}

	for testN, tt := range tests {
		command := exec.Command(binary)

		stdin, err := command.StdinPipe()
		if (err != nil) != tt.wantErr {
			t.Errorf("Command %s %v error = %v, wantErr %v", binary, tt.lines, err, tt.wantErr)
		}
		stdout, err := command.StdoutPipe()
		if (err != nil) != tt.wantErr {
			t.Errorf("Command %s %v error = %v, wantErr %v", binary, tt.lines, err, tt.wantErr)
		}

		err = command.Start()
		if (err != nil) != tt.wantErr {
			t.Errorf("Command %s %v error = %v, wantErr %v", binary, tt.lines, err, tt.wantErr)
		}

		for _, line := range tt.lines {
			io.WriteString(stdin, utils.ConcatStrings(line, "\n"))
		}
		stdin.Close()
		out := utils.ReadAllString(stdout)

		t.Logf("%s result: \n%s", tt.name, out)

		err = command.Wait()
		if (err != nil) != tt.wantErr {
			t.Errorf("Command %s %v error = %v, wantErr %v", binary, tt.lines, err, tt.wantErr)
		}

		for _, want := range tt.want {
			if !strings.Contains(out, want) {
				t.Errorf("[Test %d]: Didn't find \"%s\" in: %s", testN, want, out)
			}
		}
	}
}
