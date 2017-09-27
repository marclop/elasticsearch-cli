package client

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"

	"github.com/marclop/elasticsearch-cli/utils"
)

// HTTPCallerInterface is the HTTP implementation for caller
type HTTPCallerInterface interface {
	Do(*http.Request) (*http.Response, error)
}

// HTTP Wraps an http.Client with its config
type HTTP struct {
	Config *Config
	caller HTTPCallerInterface
}

// NewHTTP is the factory function for HTTP
func NewHTTP(config *Config, client HTTPCallerInterface) *HTTP {
	if client == nil {
		transport := http.DefaultTransport.(*http.Transport)
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: !config.insecure,
			}
		} else {
			transport.TLSClientConfig.InsecureSkipVerify = !config.insecure
		}
		client = &http.Client{
			Timeout:   config.Timeout,
			Transport: transport,
		}
	}

	return &HTTP{
		Config: config,
		caller: client,
	}
}

// TODO: Bulk operations

// HandleCall is responsible to perform HTTP requests against the secified url
// it relies on the underlying net/http.Client or Injected CallerInterface.
//
// Because we have to inject the `Content-Type: application/json`, client.Do is used.
func (c *HTTP) HandleCall(method, url, body string) (*http.Response, error) {
	var bodyIoReader io.Reader
	if body != "" {
		bodyIoReader = strings.NewReader(body)
	}

	req, err := c.createRequest(method, c.fullURL(url), bodyIoReader)
	if err != nil {
		return nil, err
	}

	return c.caller.Do(req)
}

func (c *HTTP) createRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range c.Config.headers {
		req.Header.Add(key, value)
	}

	if (c.Config.User != "") && (c.Config.Pass != "") {
		req.SetBasicAuth(c.Config.User, c.Config.Pass)
	}

	return req, nil
}

func (c *HTTP) fullURL(url string) string {
	return utils.ConcatStrings(c.Config.HTTPAdress(), url)
}

// SetHost modifies the target host
func (c *HTTP) SetHost(value string) error {
	return c.Config.SetHost(value)
}
