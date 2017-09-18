package app

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/marclop/elasticsearch-cli/cli"
	"github.com/marclop/elasticsearch-cli/client"
)

// Application contains the full application and its dependencies
type Application struct {
	config       *Config
	client       client.Client
	formatFunc   Formatter
	output       io.Writer
	indexChannel chan []string
	parser       *cli.InputParser
	poller       Poller
	repl         *readline.Instance
}

// Poller is the responsible for polling ElasticSearch and retrieving endpoints to autocomplete the CLI
type Poller interface {
	Run()
}

// Formatter formats the HTTPResponse to Stdout
type Formatter func(input *http.Response, verbose bool, interactive bool, writer io.Writer)

// Init ties all the application pieces together and returns a conveninent *Application struct
// that allows easy interaction with all the pieces of the application
func Init(config *Config, client client.Client, p *cli.InputParser, f Formatter, c chan []string, w Poller, o io.Writer) *Application {
	return &Application{
		config:       config,
		client:       client,
		formatFunc:   f,
		parser:       p,
		indexChannel: c,
		poller:       w,
		output:       o,
	}
}

// HandleCli handles the the interaction between the validated input and
// remote HTTP calls to the specified host including the call to the JSON formatter
func (app *Application) HandleCli(method string, url string, body string) error {
	res, err := app.client.HandleCall(method, url, body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	app.formatFunc(res, app.config.Verbose, app.repl != nil, app.output)
	return err
}

func (app *Application) initInteractive() {
	go app.refreshCompleter()
	go app.poller.Run()
	app.repl, _ = readline.NewEx(
		&readline.Config{
			Prompt:          "\x1b[34melasticsearch> \x1b[0m",
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",
			AutoComplete:    cli.Completer,
			HistoryFile:     "/tmp/elasticsearch-cli.history",
		},
	)
}

func (app *Application) refreshCompleter() {
	for {
		select {
		case indices := <-app.indexChannel:
			app.repl.Config.AutoComplete = cli.AssembleIndexCompleter(indices)
		}
	}
}

// Interactive runs the application like a readline / REPL
func (app *Application) Interactive() {
	app.initInteractive()
	defer app.repl.Close()
	for {
		line, err := app.repl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			return
		}

		if len(line) == 0 {
			continue
		}

		cleanLine := strings.TrimSpace(line)
		if cleanLine == "exit" || cleanLine == "quit" {
			return
		}

		lineSliced := strings.Fields(cleanLine)
		if lineSliced[0] == "set" {
			app.doSetCommands(lineSliced)
			continue
		}

		app.parser, err = cli.NewInputParser(lineSliced)
		if err != nil {
			log.Print("[ERROR]: ", err)
			continue
		}

		err = app.HandleCli(app.parser.Method, app.parser.URL, app.parser.Body)
		if err != nil {
			log.Print("[ERROR]: ", err)
		}

	}
}

func (app *Application) doSetCommands(lineSliced []string) {
	if len(lineSliced) == 3 {
		switch lineSliced[1] {
		case "host":
			err := app.client.SetHost(lineSliced[2])
			if err != nil {
				log.Print("[ERROR]: ", err)
			}
		case "port":
			port, err := strconv.Atoi(lineSliced[2])
			if err != nil {
				log.Print(lineSliced[2], " is not a valid port")
			} else {
				app.client.SetPort(port)
			}
		case "user":
			app.client.SetUser(lineSliced[2])
		case "pass":
			app.client.SetPass(lineSliced[2])
		}
	} else if (len(lineSliced) == 2) && (lineSliced[1] == "verbose") {
		app.config.Verbose = true
	}
}
