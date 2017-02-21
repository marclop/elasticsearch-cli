package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ClientInterface interface {
	HandleCall(string, string, string) (*http.Response, error)
	SetHost(string)
	SetPort(int)
	SetUser(string)
	SetPass(string)
}

type Client struct {
	config  *Config
	address string
	client  *http.Client
}

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		client: &http.Client{
			Timeout: config.GetTimeout(),
		},
	}
}

// TODO: Bulk operations

// HandleCall is responsible to perform HTTP requests against the secified url
// it heavivly relies on the underlying net/http.Client.
//
// Because we have to inject the `Content-Type: application/json`, client.Do is used.
func (c *Client) HandleCall(method string, url string, body string) (*http.Response, error) {
	var bodyIoReader io.Reader
	if body != "" {
		bodyIoReader = strings.NewReader(body)
	}

	req, err := c.createRequest(method, c.getFullUrl(url), bodyIoReader)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
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

func (c *Client) getFullUrl(url string) string {
	return fmt.Sprintf("%s%s", c.config.HttpAddress(), url)
}

func (c *Client) SetHost(value string) {
	c.config.host = value
}

func (c *Client) SetPort(value int) {
	c.config.port = value
}

func (c *Client) SetUser(value string) {
	c.config.user = value
}

func (c *Client) SetPass(value string) {
	c.config.pass = value
}
