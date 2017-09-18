package cli

import (
	"github.com/chzyer/readline"
	"github.com/marclop/elasticsearch-cli/utils"
)

var (
	getActions = []string{
		"/_aliases",
		"/_analyze",
		"/_count",
		"/_mapping",
		"/_recovery",
		"/_segments",
		"/_search",
		"/_search/_search_shards",
		"/_settings",
		"/_shard_stores",
		"/_stats",
		"/_validate/query",
	}
	putActions = []string{
		"/_mapping",
		"/_settings",
	}
	postActions = []string{
		"/_aliases",
		"/_cache/clear",
		"/_close",
		"/_flush/synced",
		"/_forcemerge",
		"/_open",
		"/_refresh",
		"/_rollover",
	}
)

var getAPI = []readline.PrefixCompleterInterface{
	readline.PcItem("_analyze"),
	readline.PcItem("_cat/indices"),
	readline.PcItem("_cat/pending_tasks"),
	readline.PcItem("_cat/repositories"),
	readline.PcItem("_cat/segments"),
	readline.PcItem("_cat/health"),
	readline.PcItem("_cat/nodes"),
	readline.PcItem("_cat/allocation"),
	readline.PcItem("_cat/shards"),
	readline.PcItem("_cat/recovery"),
	readline.PcItem("_cat/master"),
	readline.PcItem("_cat/fielddata"),
	readline.PcItem("_cat/nodeattrs"),
	readline.PcItem("_cat/count"),
	readline.PcItem("_cat/plugins"),
	readline.PcItem("_cat/templates"),
	readline.PcItem("_cat/aliases"),
	readline.PcItem("_cat/snapshots/"),
	readline.PcItem("_cluster/allocation/explain"),
	readline.PcItem("_cluster/settings"),
	readline.PcItem("_cluster/health"),
	readline.PcItem("_cluster/pending_tasks"),
	readline.PcItem("_cluster/state"),
	readline.PcItem("_cluster/stats"),
	readline.PcItem("_nodes/hot_threads"),
	// TODO: https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-nodes-stats.html
	readline.PcItem("_nodes/stats"),
	// TODO: Autodiscover Nodes? https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-nodes-info.html
	readline.PcItem("_nodes"),
	readline.PcItem("_search"),
	readline.PcItem("_search/template"),
	readline.PcItem("_stats"),
	readline.PcItem("_template"),
}

var postAPI = []readline.PrefixCompleterInterface{
	readline.PcItem("_cluster/reroute"),
	readline.PcItem("_flush"),
	readline.PcItem("_refresh"),
}

var putAPI = []readline.PrefixCompleterInterface{
	readline.PcItem("_cluster/settings"),
	readline.PcItem("_template/"),
}

var setCompleter = readline.PcItem("set",
	readline.PcItem("user"),
	readline.PcItem("pass"),
	readline.PcItem("host"),
	readline.PcItem("port"),
	readline.PcItem("verbose"),
)

// Completer has the initial list for the interactive completions
var Completer = readline.NewPrefixCompleter(
	readline.PcItem("GET", getAPI...),
	readline.PcItem("PUT", putAPI...),
	readline.PcItem("POST", postAPI...),
	readline.PcItem("HEAD"),
	readline.PcItem("DELETE"),
	setCompleter,
)

// AssembleIndexCompleter creates the autocompletion index for REPL
func AssembleIndexCompleter(indices []string) readline.PrefixCompleterInterface {
	var indexCompleterList []readline.PrefixCompleterInterface
	var indexGetOperations []readline.PrefixCompleterInterface
	var indexPutOperations []readline.PrefixCompleterInterface
	var indexPostOperations []readline.PrefixCompleterInterface

	for _, index := range indices {
		indexCompleterList = append(indexCompleterList, readline.PcItem(index))
	}

	for _, index := range indices {
		for _, endpoint := range getActions {
			indexGetOperations = append(indexGetOperations, readline.PcItem(utils.ConcatStrings(index, endpoint)))
		}

		for _, endpoint := range putActions {
			indexPutOperations = append(indexPutOperations, readline.PcItem(utils.ConcatStrings(index, endpoint)))
		}

		for _, endpoint := range postActions {
			indexPostOperations = append(indexPostOperations, readline.PcItem(utils.ConcatStrings(index, endpoint)))
		}

		for _, indexTarget := range indices {
			if index != indexTarget {
				indexPostOperations = append(indexPostOperations, readline.PcItem(utils.ConcatStrings(index, "/_shrink/", indexTarget)))
				indexPutOperations = append(indexPutOperations, readline.PcItem(utils.ConcatStrings(index, "/_alias/", indexTarget)))
			}
		}
	}

	getCompletions := append(getAPI, indexCompleterList...)
	getCompletions = append(getCompletions, indexGetOperations...)

	putCOmpletions := append(putAPI, indexCompleterList...)
	putCOmpletions = append(putCOmpletions, indexPutOperations...)

	postCompletions := append(postAPI, indexCompleterList...)
	postCompletions = append(postCompletions, indexPostOperations...)

	return readline.NewPrefixCompleter(
		readline.PcItem("GET", getCompletions...),
		readline.PcItem("PUT", putCOmpletions...),
		readline.PcItem("POST", postCompletions...),
		readline.PcItem("HEAD", indexCompleterList...),
		readline.PcItem("DELETE", indexCompleterList...),
		setCompleter,
	)
}
