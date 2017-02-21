package cli

import "github.com/chzyer/readline"

// Completer has the list for the interactive completions
var Completer = readline.NewPrefixCompleter(
	readline.PcItem("GET",
		readline.PcItem("_cat/indices"),
		readline.PcItem("_cat/pending_tasks"),
		readline.PcItem("_cat/repositories"),
		readline.PcItem("_cat/segments"),
		readline.PcItem("_cat/health"),
		readline.PcItem("_cat/nodes"),
		readline.PcItem("_cat/allocation"),
		readline.PcItem("_cat/shards"),
		readline.PcItem("_cat/recovery"),
		readline.PcItem("_cat/fielddata"),
		readline.PcItem("_cat/nodeattrs"),
		readline.PcItem("_cat/count"),
		readline.PcItem("_cat/plugins"),
		readline.PcItem("_cat/templates"),
		readline.PcItem("_cat/aliases"),
	),
	readline.PcItem("PUT"),
	readline.PcItem("POST"),
	readline.PcItem("HEAD"),
	readline.PcItem("DELETE"),
	readline.PcItem("set",
		readline.PcItem("user"),
		readline.PcItem("pass"),
		readline.PcItem("host"),
		readline.PcItem("port"),
		readline.PcItem("verbose"),
	),
)
