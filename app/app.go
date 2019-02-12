package app

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/marclop/elasticsearch-cli/cli"
	"github.com/marclop/elasticsearch-cli/client"
	"github.com/marclop/elasticsearch-cli/elasticsearch"
	"github.com/marclop/elasticsearch-cli/poller"
)

const (
	// DefaultPrompt is used when the cluster status cannot be retrieved
	DefaultPrompt = "\x1b[34melasticsearch> \x1b[0m"
	// GreenPrompt is used when the cluster status is green
	GreenPrompt = "\x1b[32melasticsearch> \x1b[0m"
	// YellowPrompt is used when the cluster status is yellow
	YellowPrompt = "\x1b[33melasticsearch> \x1b[0m"
	// RedPrompt is used when the cluster status is red
	RedPrompt = "\x1b[31melasticsearch> \x1b[0m"
)

// clusterPrompts are used in interactive mode for different prompt colors
var clusterPrompts = map[string]string{
	"green":  GreenPrompt,
	"yellow": YellowPrompt,
	"red":    RedPrompt,
}

// Application contains the full application and its dependencies
type Application struct {
	config       *Config
	client       *client.HTTP
	formatFunc   Formatter
	output       io.Writer
	indexChannel chan []string
	parser       *cli.InputParser
	poller       Poller
	repl         *readline.Instance
}

// Poller is the responsible for polling ElasticSearch and retrieving endpoints to autocomplete the CLI
type Poller interface {
	// Start the IndexPoller indefinitely, which will get the cluster indexList
	// And will send the results back to the channel
	Start()
	// Stop makes the indexPoller stop querying the Elasticsearch endpoint
	// additionally closing all of the channels
	Stop()
}

// Formatter formats the HTTPResponse to Stdout
type Formatter func(input *http.Response, verbose bool, interactive bool, writer io.Writer) error

// New creates a new instance of elasticsearch-cli from the passed Config
func New(config *Config) (*Application, error) {
	clientConfig, err := client.NewClientConfig(config.Host, config.Port, config.User, config.Pass, config.Timeout, config.Insecure)
	if err != nil {
		return nil, err
	}

	if config.Headers != nil {
		for header, value := range config.Headers {
			clientConfig.SetHeader(header, value)
		}
	}

	httpClient := client.NewHTTP(clientConfig, config.Client)
	if err != nil {
		return nil, err
	}

	indicesChannel := make(chan []string, 1)
	indexPoller := poller.NewIndexPoller(httpClient, indicesChannel, config.PollInterval)
	return initialize(config, httpClient, cli.Format, indicesChannel, indexPoller, os.Stdout), nil
}

func initialize(config *Config, client *client.HTTP, f Formatter, c chan []string, w Poller, o io.Writer) *Application {
	log.SetOutput(os.Stderr)
	return &Application{
		config:       config,
		client:       client,
		formatFunc:   f,
		indexChannel: c,
		poller:       w,
		output:       o,
	}
}

// HandleCli handles the the interaction between the validated input and
// remote HTTP calls to the specified host including the call to the JSON formatter
func (app *Application) HandleCli(args []string) error {
	input, err := cli.NewInputParser(args)
	if err != nil {
		return err
	}

	res, err := app.client.HandleCall(input.Method, input.URL, input.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return app.formatFunc(res, app.config.Verbose, app.repl != nil, app.output)
}

func (app *Application) initInteractive() {
	app.repl, _ = readline.NewEx(
		&readline.Config{
			Prompt:          app.getClusterPrompt(),
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",
			AutoComplete:    cli.Completer,
			HistoryFile:     "/tmp/elasticsearch-cli.history",
		},
	)
	go app.refreshCompleter()
	go app.poller.Start()
}

func (app *Application) getClusterPrompt() string {
	res, err := app.client.HandleCall("GET", "/_cluster/health", "")
	if err != nil {
		return DefaultPrompt
	}

	var clusterHealth elasticsearch.Health
	err = json.NewDecoder(res.Body).Decode(&clusterHealth)
	if err != nil {
		return DefaultPrompt
	}

	return clusterPrompts[clusterHealth.Status]
}

func (app *Application) refreshCompleter() {
	for {
		select {
		case indices, ok := <-app.indexChannel:
			if !ok {
				return
			}
			app.repl.Config.AutoComplete = cli.AssembleIndexCompleter(indices)
		}
	}
}

// Interactive runs the application like a readline / REPL
func (app *Application) Interactive() error {
	app.initInteractive()
	defer app.poller.Stop()
	for {
		app.repl.Config.Prompt = app.getClusterPrompt()
		line, err := app.repl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		if len(line) == 0 {
			continue
		}

		cleanLine := strings.TrimSpace(line)
		if cleanLine == "exit" || cleanLine == "quit" {
			break
		}

		input := strings.Fields(cleanLine)
		if input[0] == "set" {
			app.doSetCommands(input)
			continue
		}

		if err := app.HandleCli(input); err != nil {
			log.Print("[ERROR]: ", err)
		}
	}

	return app.repl.Close()
}

func (app *Application) doSetCommands(input []string) {
	if len(input) == 3 {
		switch input[1] {
		case "host":
			err := app.client.SetHost(input[2])
			if err != nil {
				log.Print("[ERROR]: ", err)
			}
		case "port":
			port, err := strconv.Atoi(input[2])
			if err != nil {
				log.Print(input[2], " is not a valid port")
			} else {
				app.client.Config.HostPort.Port = port
			}
		case "user":
			app.client.Config.User = input[2]
		case "pass":
			app.client.Config.Pass = input[2]
		}
	}

	if (len(input) == 2) && (input[1] == "verbose") {
		app.config.Verbose = true
	}
}
