package client

import "fmt"
import "time"

var defaultClientHeaders = map[string]string{
	"Content-Type": "application/json",
}

type Config struct {
	host    string
	port    int
	user    string
	pass    string
	timeout time.Duration
	headers map[string]string
}

func NewClientConfig(host string, port int, user string, pass string, timeout time.Duration) *Config {
	return &Config{
		host:    host,
		port:    port,
		user:    user,
		pass:    pass,
		timeout: timeout * time.Second,
		headers: defaultClientHeaders,
	}
}

func (c *Config) SetHeader(key string, value string) {
	c.headers[key] = value
}

func (c *Config) HttpAddress() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

func (c *Config) GetTimeout() time.Duration {
	return c.timeout
}

func (c *Config) SetHost(value string) {
	c.host = value
}

func (c *Config) SetPort(value int) {
	c.port = value
}

func (c *Config) SetUser(value string) {
	c.user = value
}

func (c *Config) SetPass(value string) {
	c.pass = value
}
