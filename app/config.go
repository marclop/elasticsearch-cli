package app

// Config for elasticsearch-cli Application
type Config struct {
	verbose      bool
	pollInterval int
}

// NewApplicationConfig acts as the Factory for the elasticsearch-cli config
func NewApplicationConfig(verbose bool) *Config {
	return &Config{
		verbose: verbose,
	}
}
