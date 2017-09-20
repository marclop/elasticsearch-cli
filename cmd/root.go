package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/marclop/elasticsearch-cli/app"
	"github.com/marclop/elasticsearch-cli/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version string

	// RootCmd represents the base command when called without any subcommands
	RootCmd = &cobra.Command{
		Use:               "elasticsearch-cli",
		Short:             "elasticsearch-cli provides a REPL console-like interface to interact with Elasticsearch",
		DisableAutoGenTag: true,
		ValidArgs:         append([]string{"version"}, cli.SupportedMethods...),
		ArgAliases:        append([]string{"version"}, cli.SupportedMethods...),
		RunE:              runESCLI,
	}
)

func runESCLI(cmd *cobra.Command, args []string) error {
	initConfig()

	var c app.Config
	err := viper.Unmarshal(&c)
	if err != nil {
		return err
	}

	esCli, err := app.New(&c)
	if err != nil {
		return err
	}

	// This fixes Cobra routing on children commands (the first arg is missing)
	if cmd.HasParent() {
		args = append([]string{cmd.Name()}, args...)
	}

	if len(args) > 0 {
		return esCli.HandleCli(args)
	}

	return esCli.Interactive()
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string) {
	version = v
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().String("cluster", "default", "config name, used to have multiple clusters configured in $HOME/.elasticsearch-cli/<env>")
	RootCmd.PersistentFlags().String("host", "http://localhost", "default elasticsearch URL")
	RootCmd.PersistentFlags().Int("port", 9200, "default elasticsearch port to use")
	RootCmd.PersistentFlags().StringP("user", "u", "", "username to use to authenticate (If not specified look for ES_USER environment variable)")
	RootCmd.PersistentFlags().StringP("pass", "p", "", "password to use to authenticate (If not specified, will look for ES_PASS environment variable)")
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose mode")
	RootCmd.PersistentFlags().Int("poll-interval", 10, "interval on which to poll Elasticsearch to provide index autocompletion")
	RootCmd.PersistentFlags().IntP("timeout", "t", 10, "http client timeout to the remote endpoint")
	viper.BindPFlags(RootCmd.PersistentFlags())

	for _, m := range cli.SupportedMethods {
		RootCmd.AddCommand(&cobra.Command{
			Use:     fmt.Sprintf("%s <relative endpoint> [body]", strings.ToLower(m)),
			Aliases: []string{m},
			Short:   fmt.Sprintf("Performs a %s operation against the remote endpoint", m),
			RunE:    runESCLI,
		})
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("ES")
	viper.AutomaticEnv()
	viper.AddConfigPath("$HOME/.elasticsearch-cli")
	viper.SetConfigName(viper.GetString("cluster"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using cluster config file:", viper.ConfigFileUsed())
	}
}
