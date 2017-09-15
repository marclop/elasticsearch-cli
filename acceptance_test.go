// +build acceptance

package main

import (
	"encoding/json"
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/marclop/elasticsearch-cli/utils"
)

type Response struct {
	Acknowledged bool   `json:"acknowledged"`
	Status       int    `json:"status"`
	Tagline      string `json:"tagline"`
}

func TestElasticsearchCliNonInteractive(t *testing.T) {
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
		{
			"SettingAnInvalidHostFails",
			[]string{
				"-host",
				"INVALID:wut",
				"GET",
				"/",
			},
			true,
			nil,
		},
		{
			"SettingAnInvalidHostPortFails",
			[]string{
				"-host",
				"http://localhost:asda",
				"GET",
				"/",
			},
			true,
			nil,
		},
		{
			"SettingAnInvalidPortFails",
			[]string{
				"-port",
				"asda",
				"GET",
				"/",
			},
			true,
			nil,
		},
	}

	for _, tt := range tests {
		out, err := exec.Command(binary, tt.args...).CombinedOutput()
		if (err != nil) != tt.wantErr {
			t.Errorf("Command %s %v error = %v, output = %v, wantErr = %v", binary, tt.args, err, string(out), tt.wantErr)
		}

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

func TestElasticsearchCliInteractive(t *testing.T) {
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
		{
			"SetInvalidPortFails",
			[]string{
				"set port asda",
			},
			false,
			[]string{
				`asda is not a valid port`,
			},
		},
		{
			"SetInvalidHostFails",
			[]string{
				"set host asda",
			},
			false,
			[]string{
				`host "asda" is invalid`,
			},
		},
		{
			"SetInvalidHostPortFails",
			[]string{
				"set host http://localhost:INVALID",
			},
			false,
			[]string{
				`host "http://localhost:INVALID" is invalid`,
			},
		},
	}

	for testN, tt := range tests {
		command := exec.Command(binary)

		stdin, err := command.StdinPipe()
		if (err != nil) != tt.wantErr {
			t.Errorf("[Test %d FAIL]: Command %s %v error = %v, wantErr %v", testN, binary, tt.lines, err, tt.wantErr)
		}
		stdout, err := command.StdoutPipe()
		if (err != nil) != tt.wantErr {
			t.Errorf("[Test %d FAIL]: Command %s %v error = %v, wantErr %v", testN, binary, tt.lines, err, tt.wantErr)
		}

		stderr, err := command.StderrPipe()
		if (err != nil) != tt.wantErr {
			t.Errorf("[Test %d FAIL]: Command %s %v error = %v, wantErr %v", testN, binary, tt.lines, err, tt.wantErr)
		}

		err = command.Start()
		if (err != nil) != tt.wantErr {
			t.Errorf("[Test %d FAIL]: Command %s %v error = %v, wantErr %v", testN, binary, tt.lines, err, tt.wantErr)
		}

		for _, line := range tt.lines {
			io.WriteString(stdin, utils.ConcatStrings(line, "\n"))
		}

		stdin.Close()
		out := utils.ReadAllString(stdout)
		stdErrOut := utils.ReadAllString(stderr)

		if out != "" {
			t.Logf("[Test %d INFO]: %s stdout Result: \n%s", testN, tt.name, out)
		}
		if stdErrOut != "" {
			t.Logf("[Test %d INFO]: %s stderr Result: \n%s", testN, tt.name, stdErrOut)
		}

		err = command.Wait()
		if (err != nil) != tt.wantErr {
			t.Errorf("[Test %d]: Command %s %v error = %v, wantErr %v", testN, binary, tt.lines, err, tt.wantErr)
		}

		for _, want := range tt.want {
			if !strings.Contains(out, want) {
				if !strings.Contains(stdErrOut, want) {
					t.Errorf("[Test %d stderr FAIL]: Didn't find \"%s\" in: %s", testN, want, stdErrOut)
				}
				if out != "" {
					t.Errorf("[Test %d stdout FAIL]: Didn't find \"%s\" in: %s", testN, want, out)
				}
			}
		}
	}
}
