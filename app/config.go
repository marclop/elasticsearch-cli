package app

// Config for elasticsearch-cli Application
type Config struct {
	Verbose      bool
	PollInterval int
}

// NewApplicationConfig acts as the Factory for the elasticsearch-cli config
func NewApplicationConfig(verbose bool) *Config {
	return &Config{
		Verbose: verbose,
	}
}
