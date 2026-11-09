package cmd

import "github.com/spf13/cobra"

type verificarOptions struct {
	runtimeFlags
	Entrada          string
	Saida            string
	Modo             string
	Alias            string
	BibliotecaPKCS11 string
	SlotPKCS11       string
	Porta            int
}

func newVerificarCmd() *cobra.Command {
	options := &verificarOptions{}

	command := &cobra.Command{
		Use:   "verificar",
		Short: "Executa o comando validate do assinador-verificador.jar",
		Args:  cobra.NoArgs,
		Example: "  assinador-cli verificar --entrada ./assinado.json --saida ./resultado.json --modo direto\n" +
			"  assinador-cli verificar --entrada ./assinado.json --saida ./resultado.json --modo http --porta 8080",
		RunE: options.run,
	}

	flags := command.Flags()
	flags.StringVar(&options.Entrada, "entrada", "", "Arquivo JSON de entrada.")
	flags.StringVar(&options.Saida, "saida", "", "Arquivo JSON de saida.")
	flags.StringVar(&options.Modo, "modo", "auto",
		"Estrategia de execucao: auto, http ou direto.\n"+
			"  auto   (padrao): Tenta servidor HTTP, fallback para direto se indisponivel.\n"+
			"  http   : Usa apenas servidor HTTP, falha se indisponivel.\n"+
			"  direto : Usa apenas execucao direta (sem servidor).")
	flags.StringVar(&options.Alias, "alias", "", "Nao suportado em verificar; mantido apenas para detectar uso invalido.")
	flags.StringVar(&options.BibliotecaPKCS11, "biblioteca-pkcs11", "", "Nao suportado em verificar; mantido apenas para detectar uso invalido.")
	flags.StringVar(&options.SlotPKCS11, "slot-pkcs11", "", "Nao suportado em verificar; mantido apenas para detectar uso invalido.")
	flags.IntVar(&options.Porta, "porta", defaultPort, "Porta do servidor HTTP quando --modo http for usado.")
	bindRuntimeFlags(command, &options.runtimeFlags)

	return command
}

func (o *verificarOptions) run(command *cobra.Command, _ []string) error {
	if err := ensureValidJSONPaths(o.Entrada, o.Saida); err != nil {
		return err
	}
	if command.Flags().Changed("alias") || command.Flags().Changed("biblioteca-pkcs11") || command.Flags().Changed("slot-pkcs11") {
		return validationError("O comando verificar nao aceita --alias, --biblioteca-pkcs11 nem --slot-pkcs11.")
	}

	// Converte --modo para estratégia
	strategy, err := ParseExecutionStrategy(o.Modo)
	if err != nil {
		return err
	}

	// Se HTTP foi forçado, valida porta
	if strategy == StrategyHTTP {
		if err := ensurePort(command, "http", o.Porta); err != nil {
			return err
		}
	}

	// Constrói argumentos (sem --mode)
	args := []string{
		"validate",
		"--pathin", o.Entrada,
		"--pathout", o.Saida,
	}

	// Usa RunWithStrategy
	config, err := newRunnerConfig(o.runtimeFlags)
	if err != nil {
		return wrapRuntimeError(err)
	}
	result, err := config.RunWithStrategy(args, strategy.String(), o.Porta)
	if err != nil {
		return wrapRuntimeError(err)
	}

	return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}
