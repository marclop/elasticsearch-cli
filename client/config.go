package client

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/marclop/elasticsearch-cli/utils"
)

var defaultClientHeaders = map[string]string{
	"Content-Type": "application/json",
}

// Config contains the config http.Client that will be used for the http.Client
type Config struct {
	HostPort *hostPort
	User     string
	Pass     string
	Timeout  time.Duration
	headers  map[string]string
	insecure bool
}

type hostPort struct {
	Host string
	Port int
}

// NewClientConfig handles the parameters that will be used in the HTTP Client
// If a socket is passed as a URL (http://<host>:<port>), the complex URL will prevail
// from the passed port
func NewClientConfig(host string, port int, user string, pass string, timeout int, insecure bool) (*Config, error) {
	hp, err := newHostPortString(host, port)
	if err != nil {
		return nil, err
	}

	return &Config{
		HostPort: hp,
		User:     user,
		Pass:     pass,
		Timeout:  time.Duration(timeout) * time.Second,
		headers:  defaultClientHeaders,
		insecure: insecure,
	}, nil
}

func validateSchema(host string) error {
	if govalidator.IsURL(host) {
		return nil
	}
	return fmt.Errorf("host \"%s\" is invalid", host)
}

func newHostPortString(host string, port int) (*hostPort, error) {
	err := validateSchema(host)
	if err != nil {
		return nil, err
	}

	urlString := strings.Split(host, "/")[2]
	defaultHostPort := &hostPort{strings.Join(strings.Split(host, ":")[0:2], ":"), port}
	if strings.Contains(urlString, ":") {
		urlStringPort := strings.Split(urlString, ":")[1]
		intedPort, err := strconv.Atoi(urlStringPort)
		if err != nil {
			return nil, fmt.Errorf("invalid port \"%s\"", urlStringPort)
		}
		return &hostPort{defaultHostPort.Host, intedPort}, nil
	}
	return defaultHostPort, nil
}

// SetHeader that will be sent with the request
func (c *Config) SetHeader(key string, value string) {
	c.headers[key] = value
}

// HTTPAdress returns the host and port combination so it can
// be used by the Client http://host:port
func (c *Config) HTTPAdress() string {
	return utils.ConcatStrings(c.HostPort.Host, ":", strconv.Itoa(c.HostPort.Port))
}

// SetHost modifies the target host
func (c *Config) SetHost(value string) error {
	hostPort, err := newHostPortString(value, c.HostPort.Port)
	if err == nil {
		c.HostPort = hostPort
	}
	return err
}
