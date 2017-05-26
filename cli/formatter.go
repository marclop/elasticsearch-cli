package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Format formats the HTTPResponse to Stdout
func Format(input *http.Response, verbose bool, interactive bool, writer io.Writer) {
	content, _ := ioutil.ReadAll(input.Body)
	var out bytes.Buffer
	err := json.Indent(&out, content, "", "  ")

	if interactive || verbose {
		fmt.Fprintf(writer, "Method:       %s\n", strings.ToUpper(input.Request.Method))
		fmt.Fprintf(writer, "URL:          %s\n", strings.ToLower(input.Request.URL.Path))
		if !verbose {
			fmt.Fprintln(writer)
		}
	}
	if verbose || input.Request.Method == "HEAD" {
		fmt.Fprintln(writer, "Response:    ", input.Status)
		fmt.Fprintf(writer, "Content-Type: %s\n\n", input.Request.Header["Content-Type"][0])
	}

	if input.Request.Method == "HEAD" {
		return
	}

	if err != nil {
		fmt.Fprintln(writer, strings.TrimSpace(string(content)))
		return
	}
	fmt.Fprintln(writer, strings.TrimSpace(out.String()))
}
