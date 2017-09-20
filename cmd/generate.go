package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	generatedLocation string

	generateCmd = &cobra.Command{
		Use:       "generate",
		Short:     "Generates elasticsearch-cli docs",
		ValidArgs: []string{"docs"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	docCmd = &cobra.Command{
		Use:   "docs",
		Short: "Generates the command tree documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(generatedLocation); os.IsNotExist(err) {
				os.Mkdir(generatedLocation, os.ModePerm)
			}

			return doc.GenMarkdownTree(RootCmd, generatedLocation)
		},
	}
)

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(docCmd)
	docCmd.Flags().StringVarP(&generatedLocation, "location", "l", "./docs", "Set the location of the generated output")
}
