package client

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/elasticsearch-cli/utils"
)

var defaultClientHeaders = map[string]string{
	"Content-Type": "application/json",
}

type Config struct {
	hostPort *hostPort
	user     string
	pass     string
	timeout  time.Duration
	headers  map[string]string
}

type hostPort struct {
	Host string
	Port int
}

// NewClientConfig handles the parameters that will be used in the HTTP Client
// If a socket is passed as a URL (http://<host>:<port>), the complex URL will prevail
// from the passed port
func NewClientConfig(host string, port int, user string, pass string, timeout int) (*Config, error) {
	err := validateSchema(host)
	if err != nil {
		return nil, err
	}
	hp := newHostPortString(host, port)

	return &Config{
		hostPort: hp,
		user:     user,
		pass:     pass,
		timeout:  time.Duration(timeout) * time.Second,
		headers:  defaultClientHeaders,
	}, nil
}

func validateSchema(host string) error {
	if !(strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://")) {
		return fmt.Errorf("Host doesn't contain a valid HTTP protocol (http|https) => %s", host)
	}
	return nil
}

func newHostPortString(host string, port int) *hostPort {
	urlString := strings.Split(host, "/")[2]
	if strings.Contains(urlString, ":") {
		urlStringPort := strings.Split(urlString, ":")[1]
		intedPort, _ := strconv.Atoi(urlStringPort)
		return &hostPort{strings.Join(strings.Split(host, ":")[0:2], ":"), intedPort}
	}
	return &hostPort{strings.Join(strings.Split(host, ":")[0:2], ":"), port}
}

// SetHeader that will be sent with the request
func (c *Config) SetHeader(key string, value string) {
	c.headers[key] = value
}

// HTTPAdress returns the host and port combination so it can
// be used by the Client http://host:port
func (c *Config) HTTPAdress() string {
	return utils.ConcatStrings(c.hostPort.Host, ":", strconv.Itoa(c.hostPort.Port))
}

// Timeout returns the configured HTTP timeout
func (c *Config) Timeout() time.Duration {
	return c.timeout
}

// SetHost modifies the target host
func (c *Config) SetHost(value string) error {
	err := validateSchema(value)
	if err != nil {
		return err
	}
	c.hostPort = newHostPortString(value, c.hostPort.Port)
	return nil
}

// SetPort modifies the target port
func (c *Config) SetPort(value int) {
	c.hostPort.Port = value
}

// SetUser modifies the user (HTTP Basic Auth)
func (c *Config) SetUser(value string) {
	c.user = value
}

// SetPass modifies the password (HTTP Basic Auth)
func (c *Config) SetPass(value string) {
	c.pass = value
}
