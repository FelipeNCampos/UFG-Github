package cmd

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "simulador-cli",
		Short: "CLI para gerenciar o ciclo de vida do simulador.jar",
		Long: `simulador-cli gerencia o ciclo de vida do simulador.jar.

Provisiona automaticamente o JAR e o JRE quando ausentes:
  - simulador.jar e armazenado em ~/.hubsaude/simulador.jar
  - JRE Temurin 21 e provisionado em ~/.hubsaude/jre/ se Java nao for encontrado

Exemplos de uso:
  simulador-cli iniciar            # inicia o simulador na porta padrao (8443)
  simulador-cli status             # consulta se o simulador esta ativo
  simulador-cli parar              # encerra o simulador`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(
		newIniciarCmd(),
		newPararCmd(),
		newStatusCmd(),
	)

	return rootCmd
}
