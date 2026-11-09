package cmd

import (
	"strconv"

	"github.com/felip/runner/simulador-cli/internal/provisioner"
	"github.com/felip/runner/simulador-cli/internal/runner"
	"github.com/spf13/cobra"
)

type iniciarOptions struct {
	runtimeFlags
	Porta int
}

func newIniciarCmd() *cobra.Command {
	options := &iniciarOptions{}

	command := &cobra.Command{
		Use:   "iniciar",
		Short: "Inicia o simulador.jar na porta indicada",
		Long: `Inicia o simulador.jar como servidor HTTP em segundo plano.

Antes de iniciar, o CLI verifica automaticamente:
  - Se o simulador.jar esta disponivel (baixa via release.json se necessario)
  - Se Java esta instalado (provisiona Temurin 21 em ~/.hubsaude/jre/ se ausente)
  - Se a porta esta livre (erro imediato se ja estiver em uso)

Aguarda o JSON de confirmacao de inicializacao antes de retornar.`,
		Args: cobra.NoArgs,
		Example: "  simulador-cli iniciar\n" +
			"  simulador-cli iniciar --porta 8443\n" +
			"  simulador-cli iniciar --porta 9000 --jar /caminho/para/simulador.jar\n" +
			"  simulador-cli iniciar --java-bin /usr/lib/jvm/java-21/bin/java",
		RunE: options.run,
	}

	flags := command.Flags()
	flags.IntVar(&options.Porta, "porta", defaultPort, "Porta em que o simulador sera iniciado.")
	bindRuntimeFlags(command, &options.runtimeFlags)

	return command
}

func (o *iniciarOptions) run(command *cobra.Command, _ []string) error {
	if o.Porta < 1 || o.Porta > 65535 {
		return validationError("A porta informada em --porta deve estar entre 1 e 65535.")
	}

	if runner.IsPortInUse(o.Porta) {
		return validationError(
			"A porta %d ja esta em uso. Verifique se o simulador ja esta em execucao com:\n  simulador-cli status --porta %d",
			o.Porta, o.Porta,
		)
	}

	if err := provisioner.EnsureJar(o.JarPath, command.OutOrStdout()); err != nil {
		return wrapRuntimeError(err)
	}

	resolvedJava, err := provisioner.EnsureJava(o.JavaBin, command.OutOrStdout())
	if err != nil {
		return wrapRuntimeError(err)
	}
	o.JavaBin = resolvedJava

	args := []string{"server", "start", "--port", strconv.Itoa(o.Porta)}

	result, err := newRunnerConfig(o.runtimeFlags).StartServer(o.Porta, args)
	if err != nil {
		return wrapRuntimeError(err)
	}

	return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}
