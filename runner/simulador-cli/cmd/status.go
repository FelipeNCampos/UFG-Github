package cmd

import (
	"github.com/felip/runner/simulador-cli/internal/runner"
	"github.com/spf13/cobra"
)

type statusOptions struct {
	Porta int
}

func newStatusCmd() *cobra.Command {
	options := &statusOptions{}

	command := &cobra.Command{
		Use:   "status",
		Short: "Consulta o status do simulador em execucao",
		Long: `Consulta GET /api/info no simulador em execucao na porta indicada.

Retorna o JSON de status do simulador se estiver ativo.
Retorna {"status":"OFFLINE",...} se o simulador nao responder — sem erro de processo.

Codigos de saida:
  0  consulta realizada (independente de ONLINE ou OFFLINE)
  2  porta invalida`,
		Args: cobra.NoArgs,
		Example: "  simulador-cli status\n" +
			"  simulador-cli status --porta 8443\n" +
			"  simulador-cli status --porta 9000",
		RunE: options.run,
	}

	flags := command.Flags()
	flags.IntVar(&options.Porta, "porta", defaultPort, "Porta em que o simulador esta em execucao.")

	return command
}

func (o *statusOptions) run(command *cobra.Command, _ []string) error {
	if o.Porta < 1 || o.Porta > 65535 {
		return validationError("A porta informada em --porta deve estar entre 1 e 65535.")
	}

	result, err := runner.Status(o.Porta)
	if err != nil {
		return wrapRuntimeError(err)
	}

	return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}
