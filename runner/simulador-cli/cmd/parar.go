package cmd

import (
	"github.com/felip/runner/simulador-cli/internal/runner"
	"github.com/spf13/cobra"
)

type pararOptions struct {
	Porta int
}

func newPararCmd() *cobra.Command {
	options := &pararOptions{}

	command := &cobra.Command{
		Use:   "parar",
		Short: "Envia requisicao de encerramento ao simulador em execucao",
		Long: `Encerra o simulador lendo o PID salvo em ~/.hubsaude/simulador-{porta}.pid.

Se o PID nao for encontrado ou o processo ja nao existir, retorna erro com codigo de saida 1.
Use 'simulador-cli status' para verificar se o simulador esta ativo antes de parar.`,
		Args: cobra.NoArgs,
		Example: "  simulador-cli parar\n" +
			"  simulador-cli parar --porta 8443\n" +
			"  simulador-cli parar --porta 9000",
		RunE: options.run,
	}

	flags := command.Flags()
	flags.IntVar(&options.Porta, "porta", defaultPort, "Porta em que o simulador esta em execucao.")

	return command
}

func (o *pararOptions) run(command *cobra.Command, _ []string) error {
	if o.Porta < 1 || o.Porta > 65535 {
		return validationError("A porta informada em --porta deve estar entre 1 e 65535.")
	}

	result, err := runner.Stop(o.Porta)
	if err != nil {
		return wrapRuntimeError(err)
	}

	return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}
