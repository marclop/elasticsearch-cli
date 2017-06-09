package cli

import (
	"reflect"
	"testing"

	"github.com/chzyer/readline"
)

func TestAssembleIndexCompleter(t *testing.T) {
	type args struct {
		indices []string
	}
	tests := []struct {
		name string
		args args
		want readline.PrefixCompleterInterface
	}{
		{
			"AssembleIndexCompleterSucceedsWithNoIndexes",
			args{
				[]string{},
			},
			Completer,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AssembleIndexCompleter(tt.args.indices); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AssembleIndexCompleter() = %v, want %v", got, tt.want)
			}
		})
	}
}
