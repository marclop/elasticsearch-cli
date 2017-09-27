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
	Insecure     bool   `mapstructure:"insecure"`
	Headers      map[string]string
}
