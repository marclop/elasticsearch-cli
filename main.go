package main

import "github.com/marclop/elasticsearch-cli/cmd"

// Version of elasticsearch-cli, populated at compile time
var Version string

func main() {
	cmd.Execute(Version)
}
