package cli

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/elastic/elasticsearch-cli/utils"
)

var supportedMethods = []string{
	"GET",
	"HEAD",
	"DELETE",
	"PUT",
	"POST",
}

type Parser struct {
	method string
	url    string
	body   string
}

type ParserInterface interface {
	Validate() error
	Method() string
	URL() string
	Body() string
}

// NewParser initializes the parser and validates the input
func NewParser(input []string) (ParserInterface, error) {
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

	p := &Parser{
		method: strings.ToUpper(method),
		url:    strings.ToLower(url),
		body:   body,
	}
	return p, p.Validate()
}

//TODO: Use Hashicorp multierror

// Validate makes sure that the parsed method and URL are valid
func (p *Parser) Validate() error {
	if !utils.StringInSlice(p.method, supportedMethods) {
		return fmt.Errorf("method \"%s\" is not supported", p.method)
	}
	p.ensureURLIsPrefixed()
	return nil
}

func (p *Parser) ensureURLIsPrefixed() {
	if !strings.HasPrefix(p.url, "/") {
		var buffer bytes.Buffer
		buffer.WriteString("/")
		buffer.WriteString(p.url)
		p.url = buffer.String()
	}
}

// Method returns the parsed Method in uppercase
func (p *Parser) Method() string {
	return p.method
}

// URL returns the parsed URL in in lowercase
func (p *Parser) URL() string {
	return p.url
}

// Body returns the parsed request body
func (p *Parser) Body() string {
	return p.body
}
