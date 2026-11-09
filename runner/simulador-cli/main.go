package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/felip/runner/simulador-cli/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()

	if err := rootCmd.Execute(); err != nil {
		var exitErr *cmd.ExitError
		if errors.As(err, &exitErr) {
			if exitErr.Message != "" {
				fmt.Fprintln(os.Stderr, exitErr.Message)
			}
			os.Exit(exitErr.Code)
		}

		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
