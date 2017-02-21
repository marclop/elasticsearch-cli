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
	interactive bool
	method      string
	url         string
	body        string
}

type ParserInterface interface {
	Validate() error
	Method() string
	URL() string
	Body() string
}

// NewParser returns a new *Parser and validates the input
func NewParser(input []string) (*Parser, error) {
	url := "/"
	body := ""
	method := "GET"

	if len(input) == 1 {
		method = input[0]
	}

	if len(input) > 1 {
		url = input[1]
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

// NewIteractiveParser just returns a new *Parser, skips input validation
func NewIteractiveParser(line string) (*Parser, error) {
	p, err := NewParser(strings.Fields(line))
	if err != nil {
		return nil, err
	}

	p.interactive = true
	return p, nil
}

//TODO: Use Hashicorp multierror

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

func (p *Parser) Method() string {
	return p.method
}

func (p *Parser) URL() string {
	return p.url
}

func (p *Parser) Body() string {
	return p.body
}
