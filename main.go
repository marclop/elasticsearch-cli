package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/elastic/elasticsearch-cli/app"
	"github.com/elastic/elasticsearch-cli/cli"
	"github.com/elastic/elasticsearch-cli/client"
	"github.com/elastic/elasticsearch-cli/poller"
)

// Version of elasticsearch-cli, populated at compile time
var Version string

func main() {
	hostFlag := flag.String("host", "http://localhost", "Set the ElasticSearch host url")
	portFlag := flag.Int("port", 9200, "Set the Elasticsearch Port")
	userFlag := flag.String("user", "", "Username for HTTP basic auth")
	passFlag := flag.String("pass", "", "Password for HTTP basic auth")
	timeoutFlag := flag.Int("timeout", 10, "Set the HTTP timeout")
	pollFlag := flag.Int("poll", 10, "Set the poll interval for index / endpoint autodiscovery")
	verboseFlag := flag.Bool("verbose", false, "Verbose request/response information")

	flag.Parse()
	args := flag.Args()
	if len(args) == 1 && args[0] == "version" {
		fmt.Printf("Elasticsearch CLI v%s\n", Version)
		return
	}

	clientConfig, err := client.NewClientConfig(*hostFlag, *portFlag, *userFlag, *passFlag, *timeoutFlag)
	if err != nil {
		log.Fatalf("[ERROR]: %s", err)
	}

	httpClient := client.NewHTTPClient(clientConfig, nil)
	parser, err := cli.NewInputParser(args)
	if err != nil {
		log.Fatalf("[ERROR]: %s", err)
	}

	indicesChannel := make(chan []string, 1)
	indexPoller := poller.NewIndexPoller(httpClient, indicesChannel, *pollFlag)
	appConfig := app.NewApplicationConfig(*verboseFlag)
	application := app.Init(appConfig, httpClient, parser, cli.Format, indicesChannel, indexPoller, os.Stdout)

	if len(args) > 0 {
		err := application.HandleCli(parser.Method(), parser.URL(), parser.Body())
		if err != nil {
			log.Fatalf("[ERROR]: %s", err)
		}
	} else {
		application.Interactive()
	}
}
