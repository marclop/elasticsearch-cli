package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/elastic/elasticsearch-cli/cli"
	"github.com/elastic/elasticsearch-cli/client"
	"github.com/elastic/elasticsearch-cli/elasticshell"
)

func main() {
	hostFlag := flag.String("host", "http://localhost", "Set the ElasticSearch host url")
	portFlag := flag.Int("port", 9200, "Set the Elasticsearch Port")
	userFlag := flag.String("user", "", "Username for HTTP basic auth")
	passFlag := flag.String("pass", "", "Password for HTTP basic auth")
	timeoutFlag := flag.Int("timeout", 10, "Set the HTTP client timeout")
	// pollFlag := flag.Int("poll", 5, "Set the poll interval")
	verboseFlag := flag.Bool("verbose", false, "Verbose request/response information")

	flag.Parse()
	args := flag.Args()
	if len(args) == 1 && args[0] == "version" {
		fmt.Printf("Elasticsearch-cli v%s\n", elasticshell.GetVersion())
		return
	}

	clientConfig := client.NewClientConfig(*hostFlag, *portFlag, *userFlag, *passFlag, time.Duration(*timeoutFlag))
	client := client.NewClient(clientConfig)
	parser, err := cli.NewParser(args)
	if err != nil {
		fmt.Println(err)
	}

	// TODO: Reenable poll for index auto discovery
	elasticShellConfig := elasticshell.NewApplicationConfig(*verboseFlag, 5)
	elasticShell := elasticshell.Init(elasticShellConfig, client, parser)

	if len(args) > 0 {
		err := elasticShell.HandleCli(parser.Method(), parser.URL(), parser.Body())
		if err != nil {
			fmt.Println(err)
		}
	} else {
		elasticShell.Interactive()
	}
}
