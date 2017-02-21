package elasticshell

// Config for elasticshell Application
type Config struct {
	verbose      bool
	pollInterval int
	interactive  bool
}

// NewApplicationConfig acts as the Factory for the elasticshell config
func NewApplicationConfig(verbose bool, pollInterval int) *Config {
	return &Config{
		verbose:      verbose,
		pollInterval: pollInterval,
	}
}
