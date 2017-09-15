package cli

import (
	"fmt"
	"strings"

	"github.com/marclop/elasticsearch-cli/utils"
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
	method string
	url    string
	body   string
}

// NewInputParser initializes the parser and validates the input
func NewInputParser(input []string) (*InputParser, error) {
	url := "/"
	body := ""
	method := "GET"

	if len(input) == 0 {
		return nil, nil
	}

	method = strings.ToUpper(input[0])

	if len(input) > 1 {
		url = strings.ToLower(input[1])
	}

	if len(input) == 3 {
		body = input[2]
	}

	p := &InputParser{
		method: strings.ToUpper(method),
		url:    strings.ToLower(url),
		body:   body,
	}

	err := p.Validate()
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Validate makes sure that the parsed method and URL are valid
func (p *InputParser) Validate() error {
	if !utils.StringInSlice(p.method, supportedMethods) {
		return fmt.Errorf("method \"%s\" is not supported", p.method)
	}
	p.ensureURLIsPrefixed()
	return nil
}

func (p *InputParser) ensureURLIsPrefixed() {
	if !strings.HasPrefix(p.url, "/") {
		p.url = utils.ConcatStrings("/", p.url)
	}
}

// Method returns the parsed Method in uppercase
func (p *InputParser) Method() string {
	return p.method
}

// URL returns the parsed URL in in lowercase
func (p *InputParser) URL() string {
	return p.url
}

// Body returns the parsed request body
func (p *InputParser) Body() string {
	return p.body
}
