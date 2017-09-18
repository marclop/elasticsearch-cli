package cli

import (
	"fmt"
	"strings"

	"github.com/marclop/elasticsearch-cli/utils"
)

const (
	defaultURL = "/"
)

var supportedMethods = []string{
	"GET",
	"HEAD",
	"DELETE",
	"PUT",
	"POST",
}

// InputParser is the struct that parses the input into something usable by the
// application
type InputParser struct {
	Method string
	URL    string
	Body   string
}

// NewInputParser initializes the parser and validates the input
func NewInputParser(input []string) (*InputParser, error) {
	if len(input) == 0 {
		return nil, nil
	}

	var inputParser = &InputParser{
		Method: strings.ToUpper(input[0]),
		URL:    defaultURL,
		Body:   "",
	}

	if len(input) > 1 {
		inputParser.URL = strings.ToLower(input[1])
	}

	if len(input) > 2 {
		inputParser.Body = strings.Join(input[2:], "")
	}

	err := inputParser.Validate()
	if err != nil {
		return nil, err
	}

	return inputParser, nil
}

// Validate makes sure that the parsed Method and URL are valid
func (p *InputParser) Validate() error {
	if !utils.StringInSlice(p.Method, supportedMethods) {
		return fmt.Errorf("Method \"%s\" is not supported", p.Method)
	}
	p.ensureURLIsPrefixed()
	return nil
}

func (p *InputParser) ensureURLIsPrefixed() {
	if !strings.HasPrefix(p.URL, "/") {
		p.URL = utils.ConcatStrings("/", p.URL)
	}
}
