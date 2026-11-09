package cmd

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "assinador-cli",
		Short:         "CLI em Go para o assinador-verificador.jar",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(
		newAssinarCmd(),
		newVerificarCmd(),
		newServidorCmd(),
	)

	return rootCmd
}
