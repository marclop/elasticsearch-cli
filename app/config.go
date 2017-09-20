package app

// Config for elasticsearch-cli Application
type Config struct {
	User         string `mapstructure:"user"`
	Pass         string `mapstructure:"pass"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Verbose      bool   `mapstructure:"verbose"`
	PollInterval int    `mapstructure:"poll-interval"`
	Timeout      int    `mapstructure:"timeout"`
}

// NewApplicationConfig acts as the Factory for the elasticsearch-cli config
func NewApplicationConfig(verbose bool) *Config {
	return &Config{
		Verbose: verbose,
	}
}
