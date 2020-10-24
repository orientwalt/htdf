package main

import (
	"os"

	"github.com/orientwalt/htdf/cmd/iris/cmd"
	_ "github.com/orientwalt/htdf/lite/statik"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		os.Exit(1)
	}
}
