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

const (
	methodFormat      = "Method:      "
	urlFormat         = "URL:         "
	responseFormat    = "Response:    "
	contentTypeFormat = "Content-Type:"
)

// Format formats the HTTPResponse to the output io.Writer
func Format(input *http.Response, verbose bool, interactive bool, output io.Writer) error {
	content, err := ioutil.ReadAll(input.Body)
	if err != nil {
		return err
	}
	// Removes any extra spaces the body might be carrying
	content = []byte(strings.TrimSpace(string(content)))

	var headers = new(bytes.Buffer)
	if interactive || verbose {
		headers.WriteString(
			fmt.Sprintln(methodFormat, strings.ToUpper(input.Request.Method)),
		)
		headers.WriteString(
			fmt.Sprintln(urlFormat, strings.ToLower(input.Request.URL.Path)),
		)
	}

	if verbose || input.Request.Method == "HEAD" {
		headers.WriteString(
			fmt.Sprintln(responseFormat, input.Status),
		)
		headers.WriteString(
			fmt.Sprintln(contentTypeFormat, input.Request.Header["Content-Type"][0]),
		)
	}

	// Print the headers
	if headers.String() != "" {
		headers.WriteString(fmt.Sprintln())
		headers.WriteTo(output)
	}

	// Print the response content
	var payload bytes.Buffer
	err = json.Indent(&payload, content, "", "  ")
	if err != nil {
		payload.Write(content)
	}

	if payload.String() != "" {
		payload.WriteString("\n")
		payload.WriteTo(output)
	}

	return nil
}
