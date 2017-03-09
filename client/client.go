package client

import (
	"io"
	"net/http"
	"strings"

	"github.com/elastic/elasticsearch-cli/utils"
)

type ClientInterface interface {
	HandleCall(string, string, string) (*http.Response, error)
	SetHost(string) error
	SetPort(int)
	SetUser(string)
	SetPass(string)
}

type CallerInterface interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	config *Config
	caller CallerInterface
}

func NewClient(config *Config, client CallerInterface) ClientInterface {
	if client == nil {
		client = &http.Client{
			Timeout: config.Timeout(),
		}
	}
	return &Client{
		config: config,
		caller: client,
	}
}

// TODO: Bulk operations

// HandleCall is responsible to perform HTTP requests against the secified url
// it relies on the underlying net/http.Client or Injected CallerInterface.
//
// Because we have to inject the `Content-Type: application/json`, client.Do is used.
func (c *Client) HandleCall(method string, url string, body string) (*http.Response, error) {
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

func (c *Client) createRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range c.config.headers {
		req.Header.Add(key, value)
	}

	if (c.config.user != "") && (c.config.pass != "") {
		req.SetBasicAuth(c.config.user, c.config.pass)
	}

	return req, nil
}

func (c *Client) fullURL(url string) string {
	return utils.ConcatStrings(c.config.HTTPAdress(), url)
}

// SetHost modifies the target host
func (c *Client) SetHost(value string) error {
	return c.config.SetHost(value)
}

// SetPort modifies the target port
func (c *Client) SetPort(value int) {
	c.config.SetPort(value)
}

// SetUser modifies the user (HTTP Basic Auth)
func (c *Client) SetUser(value string) {
	c.config.user = value
}

// SetPass modifies the password (HTTP Basic Auth)
func (c *Client) SetPass(value string) {
	c.config.pass = value
}
