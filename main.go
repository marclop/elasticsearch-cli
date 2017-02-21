package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/elastic/elasticsearch-cli/app"
	"github.com/elastic/elasticsearch-cli/cli"
	"github.com/elastic/elasticsearch-cli/client"
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
		fmt.Printf("Elasticsearch CLI v%s\n", app.GetVersion())
		return
	}

	clientConfig := client.NewClientConfig(*hostFlag, *portFlag, *userFlag, *passFlag, time.Duration(*timeoutFlag))
	client := client.NewClient(clientConfig)
	parser, err := cli.NewParser(args)
	if err != nil {
		fmt.Println(err)
	}

	// TODO: Reenable poll for index auto discovery
	appConfig := app.NewApplicationConfig(*verboseFlag, 5)
	app := app.Init(appConfig, client, parser)

	if len(args) > 0 {
		err := app.HandleCli(parser.Method(), parser.URL(), parser.Body())
		if err != nil {
			fmt.Println(err)
		}
	} else {
		app.Interactive()
	}
}
