package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Formatter struct {
	input       *http.Response
	interactive bool
}

type FormatterInterface interface {
	FormatJSON(bool)
}

//TODO: Revisit implementation

func NewJSONFormatter(input *http.Response) *Formatter {
	return &Formatter{
		input:       input,
		interactive: false,
	}
}

func NewIteractiveJSONFormatter(input *http.Response) *Formatter {
	return &Formatter{
		input:       input,
		interactive: true,
	}
}

func (f *Formatter) FormatJSON(verbose bool) {
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
