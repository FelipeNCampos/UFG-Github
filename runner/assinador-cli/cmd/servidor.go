package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
)

type servidorActionOptions struct {
	runtimeFlags
	Porta   int
	Timeout int
}

func newServidorCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "servidor",
		Short: "Gerencia o modo server do assinador-verificador.jar",
		Args:  cobra.NoArgs,
	}

	command.AddCommand(
		newServidorIniciarCmd(),
		newServidorStatusCmd(),
		newServidorPararCmd(),
	)

	return command
}

func newServidorIniciarCmd() *cobra.Command {
	options := &servidorActionOptions{}

	command := &cobra.Command{
		Use:   "iniciar",
		Short: "Inicia o servidor HTTP do assinador-verificador.jar",
		Args:  cobra.NoArgs,
		Example: "  assinador-cli servidor iniciar --porta 8080\n" +
			"  assinador-cli servidor iniciar --porta 8081 --timeout 30",
		RunE: options.runIniciar,
	}

	flags := command.Flags()
	flags.IntVar(&options.Porta, "porta", defaultPort, "Porta do servidor HTTP.")
	flags.IntVar(&options.Timeout, "timeout", 0, "Tempo limite de inatividade, em minutos.")
	bindRuntimeFlags(command, &options.runtimeFlags)

	return command
}

func newServidorStatusCmd() *cobra.Command {
	options := &servidorActionOptions{}

	command := &cobra.Command{
		Use:     "status",
		Short:   "Consulta o status do servidor HTTP do assinador-verificador.jar",
		Args:    cobra.NoArgs,
		Example: "  assinador-cli servidor status --porta 8080",
		RunE:    options.runStatus,
	}

	flags := command.Flags()
	flags.IntVar(&options.Porta, "porta", defaultPort, "Porta do servidor HTTP.")
	bindRuntimeFlags(command, &options.runtimeFlags)

	return command
}

func newServidorPararCmd() *cobra.Command {
	options := &servidorActionOptions{}

	command := &cobra.Command{
		Use:     "parar",
		Short:   "Encerra o servidor HTTP do assinador-verificador.jar",
		Args:    cobra.NoArgs,
		Example: "  assinador-cli servidor parar --porta 8080",
		RunE:    options.runParar,
	}

	flags := command.Flags()
	flags.IntVar(&options.Porta, "porta", defaultPort, "Porta do servidor HTTP.")
	bindRuntimeFlags(command, &options.runtimeFlags)

	return command
}

func (o *servidorActionOptions) runIniciar(command *cobra.Command, _ []string) error {
	if o.Porta < 1 || o.Porta > 65535 {
		return validationError("A porta informada em --porta deve estar entre 1 e 65535.")
	}

	timeout, err := parseOptionalTimeout(command, o.Timeout)
	if err != nil {
		return err
	}

	args := []string{"server", "start", "--port", strconv.Itoa(o.Porta)}
	if timeout > 0 {
		args = append(args, "--timeout", strconv.Itoa(timeout))
	}

	config, err := newRunnerConfig(o.runtimeFlags)
	if err != nil {
		return wrapRuntimeError(err)
	}

	result, err := config.StartServer(args)
	if err != nil {
		return wrapRuntimeError(err)
	}

	return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}

func (o *servidorActionOptions) runStatus(command *cobra.Command, _ []string) error {
	if o.Porta < 1 || o.Porta > 65535 {
		return validationError("A porta informada em --porta deve estar entre 1 e 65535.")
	}

	config, err := newRunnerConfig(o.runtimeFlags)
	if err != nil {
		return wrapRuntimeError(err)
	}

	result, err := config.Run([]string{"server", "status", "--port", strconv.Itoa(o.Porta)})
	if err != nil {
		return wrapRuntimeError(err)
	}

	return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}

func (o *servidorActionOptions) runParar(command *cobra.Command, _ []string) error {
	if o.Porta < 1 || o.Porta > 65535 {
		return validationError("A porta informada em --porta deve estar entre 1 e 65535.")
	}

	config, err := newRunnerConfig(o.runtimeFlags)
	if err != nil {
		return wrapRuntimeError(err)
	}

	result, err := config.Run([]string{"server", "stop", "--port", strconv.Itoa(o.Porta)})
	if err != nil {
		return wrapRuntimeError(err)
	}

	return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}
