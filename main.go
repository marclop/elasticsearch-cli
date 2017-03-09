package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/elastic/elasticsearch-cli/app"
	"github.com/elastic/elasticsearch-cli/cli"
	"github.com/elastic/elasticsearch-cli/client"
	"github.com/elastic/elasticsearch-cli/poller"
)

func main() {
	hostFlag := flag.String("host", "http://localhost", "Set the ElasticSearch host url")
	portFlag := flag.Int("port", 9200, "Set the Elasticsearch Port")
	userFlag := flag.String("user", "", "Username for HTTP basic auth")
	passFlag := flag.String("pass", "", "Password for HTTP basic auth")
	timeoutFlag := flag.Int("timeout", 10, "Set the HTTP client timeout")
	pollFlag := flag.Int("poll", 10, "Set the poll interval for index autodiscovery")
	verboseFlag := flag.Bool("verbose", false, "Verbose request/response information")
	// helpFlag := flag.Bool("-help", false, "Prints this message")

	flag.Parse()
	args := flag.Args()
	if len(args) == 1 && args[0] == "version" {
		fmt.Printf("Elasticsearch CLI v%s\n", app.Version())
		return
	}

	clientConfig, err := client.NewClientConfig(*hostFlag, *portFlag, *userFlag, *passFlag, *timeoutFlag)
	if err != nil {
		log.Fatalf("[ERROR]: %s", err)
	}
	client := client.NewClient(clientConfig, nil)

	parser, err := cli.NewParser(args)
	if err != nil {
		log.Fatalf("[ERROR]: %s", err)
	}

	indicesChannel := make(chan []string, 1)
	poller := poller.NewIndexPoller(client, indicesChannel, *pollFlag)
	appConfig := app.NewApplicationConfig(*verboseFlag, *pollFlag)
	app := app.Init(appConfig, client, parser, indicesChannel, poller)

	if len(args) > 0 {
		err := app.HandleCli(parser.Method(), parser.URL(), parser.Body())
		if err != nil {
			log.Fatalf("[ERROR]: %s", err)
		}
	} else {
		app.Interactive()
	}
}
