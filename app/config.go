package app

// Config for elasticsearch-cli Application
type Config struct {
	verbose      bool
	pollInterval int
}

// NewApplicationConfig acts as the Factory for the elasticsearch-cli config
func NewApplicationConfig(verbose bool, pollInterval int) *Config {
	return &Config{
		verbose:      verbose,
		pollInterval: pollInterval,
	}
}
