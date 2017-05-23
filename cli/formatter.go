package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type JSONFormatter struct {
	input       *http.Response
	interactive bool
}

type Formatter interface {
	Format(bool)
}

// NewJSONFormatter initializes a JSON Formatter non interactive
func NewJSONFormatter(input *http.Response) *JSONFormatter {
	return &JSONFormatter{
		input:       input,
		interactive: false,
	}
}

// NewIteractiveJSONFormatter initializes a JSON Formatter interactive
func NewIteractiveJSONFormatter(input *http.Response) *JSONFormatter {
	return &JSONFormatter{
		input:       input,
		interactive: true,
	}
}

// Format formats the HTTPResponse to Stdout
func (f *JSONFormatter) Format(verbose bool) {
	content, _ := ioutil.ReadAll(f.input.Body)
	var out bytes.Buffer
	err := json.Indent(&out, content, "", "  ")

	if f.interactive || verbose {
		fmt.Printf("Method:       %s\n", strings.ToUpper(f.input.Request.Method))
		fmt.Printf("URL:          %s\n", strings.ToLower(f.input.Request.URL.Path))
		if !verbose {
			fmt.Println()
		}
	}
	if verbose || f.input.Request.Method == "HEAD" {
		fmt.Println("Response:    ", f.input.Status)
		fmt.Printf("Content-Type: %s\n\n", f.input.Request.Header["Content-Type"][0])
	}

	if f.input.Request.Method == "HEAD" {
		return
	}

	if err != nil {
		fmt.Println(strings.TrimSpace(string(content)))
		return
	}
	fmt.Println(strings.TrimSpace(out.String()))
}
