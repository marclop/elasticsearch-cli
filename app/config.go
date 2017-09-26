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
	Headers      map[string]string
}
