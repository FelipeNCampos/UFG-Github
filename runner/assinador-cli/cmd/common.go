package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/felip/runner/assinador-cli/internal/javaruntime"
	"github.com/felip/runner/assinador-cli/internal/runner"
	"github.com/spf13/cobra"
)

const (
	defaultJarPath = "./assinador-verificador/build/dist/assinador-verificador.jar"
	defaultPort    = 8080
)

var (
	aliasPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{3,64}$`)
	digitsOnly   = regexp.MustCompile(`^\d+$`)
)

type runtimeFlags struct {
	JavaBin string
	JarPath string
}

// ExecutionStrategy define como o assinador.jar será invocado
type ExecutionStrategy string

const (
	StrategyAuto   ExecutionStrategy = "auto"   // Tenta HTTP, fallback para direto
	StrategyHTTP   ExecutionStrategy = "http"   // Força HTTP
	StrategyDirect ExecutionStrategy = "direct" // Força direto
)

func (s ExecutionStrategy) String() string {
	return string(s)
}

// ParseExecutionStrategy converte entrada do usuário para estratégia
// Entrada vazia padrão para StrategyAuto (novo comportamento)
func ParseExecutionStrategy(input string) (ExecutionStrategy, error) {
	normalized := strings.ToLower(strings.TrimSpace(input))

	// Padrão: string vazia mapeia para auto
	if normalized == "" {
		return StrategyAuto, nil
	}

	switch normalized {
	case "auto":
		return StrategyAuto, nil
	case "http":
		return StrategyHTTP, nil
	case "direto", "direct": // Suporta português e inglês
		return StrategyDirect, nil
	default:
		return "", validationError(
			"Estrategia de modo invalida: %s. Use 'auto', 'http' ou 'direto'.",
			input,
		)
	}
}

func bindRuntimeFlags(command *cobra.Command, target *runtimeFlags) {
	flags := command.Flags()
	flags.StringVar(&target.JavaBin, "java-bin", "", "Executavel Java opcional para sobrescrever o runtime gerenciado automaticamente.")
	flags.StringVar(&target.JarPath, "jar", defaultJarPath, "Caminho para o arquivo assinador-verificador.jar.")
}

func newRunnerConfig(runtime runtimeFlags) (runner.Config, error) {
	javaBin, err := javaruntime.NewDefaultResolver().Resolve(runtime.JavaBin)
	if err != nil {
		return runner.Config{}, err
	}

	return runner.Config{
		JavaBin: javaBin,
		JarPath: runtime.JarPath,
	}, nil
}

func validationError(message string, args ...any) error {
	return &ExitError{
		Code:    2,
		Message: fmt.Sprintf(message, args...),
	}
}

func emitResult(stdout io.Writer, stderr io.Writer, result runner.Result) error {
	if result.Stdout != "" {
		if _, err := io.WriteString(stdout, result.Stdout); err != nil {
			return err
		}
		if !strings.HasSuffix(result.Stdout, "\n") {
			if _, err := io.WriteString(stdout, "\n"); err != nil {
				return err
			}
		}
	}

	if result.Stderr != "" {
		if _, err := io.WriteString(stderr, result.Stderr); err != nil {
			return err
		}
		if !strings.HasSuffix(result.Stderr, "\n") {
			if _, err := io.WriteString(stderr, "\n"); err != nil {
				return err
			}
		}
	}

	if result.ExitCode != 0 {
		return &ExitError{Code: result.ExitCode}
	}

	return nil
}

func wrapRuntimeError(err error) error {
	if err == nil {
		return nil
	}
	return &ExitError{
		Code:    2,
		Message: err.Error(),
	}
}

func modeToJarValue(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "direto":
		return "direct", nil
	case "http":
		return "http", nil
	default:
		return "", validationError("Modo invalido: %s. Use --modo direto ou --modo http.", raw)
	}
}

func ensureValidJSONPaths(inputPath string, outputPath string) error {
	inputPath = strings.TrimSpace(inputPath)
	outputPath = strings.TrimSpace(outputPath)

	if inputPath == "" {
		return validationError("A flag obrigatoria --entrada nao foi informada.")
	}
	if outputPath == "" {
		return validationError("A flag obrigatoria --saida nao foi informada.")
	}

	inputInfo, err := os.Stat(inputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return validationError("Arquivo de entrada nao encontrado: %s", inputPath)
		}
		return validationError("Nao foi possivel acessar o arquivo de entrada: %s", inputPath)
	}
	if inputInfo.IsDir() {
		return validationError("O caminho informado em --entrada precisa apontar para um arquivo JSON.")
	}

	if !strings.HasSuffix(strings.ToLower(inputPath), ".json") {
		return validationError("O arquivo informado em --entrada deve ter extensao .json.")
	}
	if !strings.HasSuffix(strings.ToLower(outputPath), ".json") {
		return validationError("O arquivo informado em --saida deve ter extensao .json.")
	}

	normalizedInput := filepath.Clean(inputPath)
	normalizedOutput := filepath.Clean(outputPath)
	if normalizedInput == normalizedOutput {
		return validationError("Os caminhos de entrada e saida nao podem apontar para o mesmo arquivo.")
	}

	outputDir := filepath.Dir(outputPath)
	if outputDir != "." {
		outputInfo, err := os.Stat(outputDir)
		if err != nil {
			if os.IsNotExist(err) {
				return validationError("O diretorio de saida nao existe: %s", outputDir)
			}
			return validationError("Nao foi possivel acessar o diretorio de saida: %s", outputDir)
		}
		if !outputInfo.IsDir() {
			return validationError("O caminho pai de --saida nao e um diretorio valido: %s", outputDir)
		}
	}

	return nil
}

func ensureValidAlias(alias string) error {
	if strings.TrimSpace(alias) == "" {
		return validationError("A flag obrigatoria --alias nao foi informada.")
	}
	if !aliasPattern.MatchString(alias) {
		return validationError("O alias deve ter entre 3 e 64 caracteres e usar apenas letras, numeros, '-' ou '_'.")
	}
	return nil
}

func ensureValidPKCS11Library(path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}

	normalized := strings.ToLower(path)
	if strings.HasSuffix(normalized, ".dll") || strings.HasSuffix(normalized, ".so") || strings.HasSuffix(normalized, ".dylib") {
		return nil
	}

	return validationError("A biblioteca informada em --biblioteca-pkcs11 deve terminar com .dll, .so ou .dylib.")
}

func ensureValidPKCS11Slot(slot string) error {
	if strings.TrimSpace(slot) == "" {
		return nil
	}
	if !digitsOnly.MatchString(slot) {
		return validationError("O valor informado em --slot-pkcs11 deve ser um numero inteiro nao negativo.")
	}
	return nil
}

func ensurePort(command *cobra.Command, mode string, port int) error {
	if mode == "http" {
		if port < 1 || port > 65535 {
			return validationError("A porta informada em --porta deve estar entre 1 e 65535.")
		}
		return nil
	}

	if command.Flags().Changed("porta") {
		return validationError("A flag --porta so pode ser usada com --modo http.")
	}

	return nil
}

func parseOptionalTimeout(command *cobra.Command, timeout int) (int, error) {
	if !command.Flags().Changed("timeout") {
		return 0, nil
	}
	if timeout < 1 {
		return 0, validationError("O valor informado em --timeout deve ser um numero positivo de minutos.")
	}
	return timeout, nil
}

func appendPortIfNeeded(args []string, mode string, port int) []string {
	if mode == "http" {
		return append(args, "--port", strconv.Itoa(port))
	}
	return args
}
