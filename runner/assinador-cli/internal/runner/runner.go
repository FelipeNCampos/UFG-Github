package runner

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Config struct {
	JavaBin string
	JarPath string
}

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	PID      int
}

// IsServerAvailable verifica se o servidor responde na porta dada
// Usa teste de conexão TCP (funciona independente de implementação HTTP)
// Retorna true se conexão sucede, false se timeout ou recusada
func IsServerAvailable(port int, timeoutSecs int) bool {
	address := fmt.Sprintf("localhost:%d", port)
	timeout := time.Duration(timeoutSecs) * time.Second

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

// DetermineExecutionMode decide o modo final baseado na estratégia
// Para "auto": retorna "http" se servidor disponível, senão "direct"
// Para "http": retorna "http"
// Para "direct": retorna "direct"
func DetermineExecutionMode(strategy string, port int) (string, error) {
	switch strings.ToLower(strategy) {
	case "direct":
		return "direct", nil

	case "http":
		return "http", nil

	case "auto":
		// Tenta HTTP com 1 segundo de timeout
		if IsServerAvailable(port, 1) {
			return "http", nil
		}
		// Fallback para direto
		return "direct", nil

	default:
		return "", fmt.Errorf("estrategia desconhecida: %s", strategy)
	}
}

// ApplyExecutionMode modifica os argumentos do comando para setar --mode e --port
// Remove qualquer flag --mode existente, adiciona nova --mode, e --port se necessário
func ApplyExecutionMode(args []string, mode string, port int) []string {
	// Remove flag --mode existente
	filtered := RemoveModeFromArgs(args)

	// Adiciona nova --mode
	filtered = append(filtered, "--mode", mode)

	// Adiciona --port se modo é http e não está presente
	if mode == "http" {
		hasPort := false
		for j := 0; j < len(filtered); j++ {
			if filtered[j] == "--port" {
				hasPort = true
				break
			}
		}
		if !hasPort {
			filtered = append(filtered, "--port", fmt.Sprintf("%d", port))
		}
	}

	return filtered
}

// RemoveModeFromArgs remove qualquer flag --mode e seu valor dos argumentos
func RemoveModeFromArgs(args []string) []string {
	result := []string{}
	i := 0
	for i < len(args) {
		if args[i] == "--mode" {
			i += 2 // Pula flag e valor
			continue
		}
		result = append(result, args[i])
		i++
	}
	return result
}

// RunWithStrategy executa com seleção automática de modo baseada na estratégia
// Se strategy é "auto": tenta HTTP primeiro, fallback para direto se necessário
// Se strategy é "http": força modo HTTP
// Se strategy é "direct": força modo direto
func (c Config) RunWithStrategy(args []string, strategy string, port int) (Result, error) {
	if err := c.Validate(); err != nil {
		return Result{}, err
	}

	// Determina modo final baseado em estratégia e disponibilidade
	mode, err := DetermineExecutionMode(strategy, port)
	if err != nil {
		return Result{}, err
	}

	// Aplica modo aos argumentos
	finalArgs := ApplyExecutionMode(args, mode, port)

	// Executa com modo determinado
	return c.Run(finalArgs)
}

func (c Config) Run(args []string) (Result, error) {
	if err := c.Validate(); err != nil {
		return Result{}, err
	}

	command := exec.Command(c.JavaBin, c.jarArgs(args)...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	result := Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if err == nil {
		return result, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.ExitCode = exitErr.ExitCode()
		return result, nil
	}

	return Result{}, fmt.Errorf("Falha ao executar o Java. Verifique o valor de --java-bin e tente novamente: %w", err)
}

func (c Config) StartServer(args []string) (Result, error) {
	if err := c.Validate(); err != nil {
		return Result{}, err
	}

	command := exec.Command(c.JavaBin, c.jarArgs(args)...)
	stdoutPipe, err := command.StdoutPipe()
	if err != nil {
		return Result{}, fmt.Errorf("Nao foi possivel preparar a captura do stdout do servidor: %w", err)
	}
	stderrPipe, err := command.StderrPipe()
	if err != nil {
		return Result{}, fmt.Errorf("Nao foi possivel preparar a captura do stderr do servidor: %w", err)
	}

	if err := command.Start(); err != nil {
		return Result{}, fmt.Errorf("Falha ao iniciar o Java. Verifique o valor de --java-bin e tente novamente: %w", err)
	}

	type stdoutReadResult struct {
		Output string
		Err    error
	}

	stdoutCh := make(chan stdoutReadResult, 1)
	stderrCh := make(chan string, 1)
	waitCh := make(chan error, 1)

	go func() {
		output, readErr := readFirstJSONObject(stdoutPipe)
		stdoutCh <- stdoutReadResult{Output: output, Err: readErr}
	}()

	go func() {
		data, _ := io.ReadAll(stderrPipe)
		stderrCh <- string(data)
	}()

	go func() {
		waitCh <- command.Wait()
	}()

	select {
	case stdoutResult := <-stdoutCh:
		if stdoutResult.Err == nil && strings.TrimSpace(stdoutResult.Output) != "" {
			return Result{
				Stdout: stdoutResult.Output,
				PID:    command.Process.Pid,
			}, nil
		}

		waitErr := <-waitCh
		stderr := <-stderrCh
		result := Result{
			Stdout:   stdoutResult.Output,
			Stderr:   stderr,
			ExitCode: exitCodeFromWait(waitErr),
		}
		if result.ExitCode == 0 {
			return Result{}, fmt.Errorf("O processo Java encerrou sem retornar o JSON de inicializacao do servidor.")
		}
		return result, nil

	case waitErr := <-waitCh:
		stderr := <-stderrCh
		stdoutResult := <-stdoutCh
		result := Result{
			Stdout:   stdoutResult.Output,
			Stderr:   stderr,
			ExitCode: exitCodeFromWait(waitErr),
		}
		if result.ExitCode == 0 && strings.TrimSpace(result.Stdout) == "" && strings.TrimSpace(result.Stderr) == "" {
			return Result{}, fmt.Errorf("O processo Java encerrou sem retornar nenhuma saida.")
		}
		return result, nil
	}
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.JavaBin) == "" {
		return fmt.Errorf("Nenhum executavel Java foi informado. Use --java-bin para apontar para o Java instalado.")
	}
	if strings.TrimSpace(c.JarPath) == "" {
		return fmt.Errorf("Nenhum caminho de JAR foi informado. Use --jar para apontar para o assinador-verificador.jar.")
	}

	if _, err := exec.LookPath(c.JavaBin); err != nil {
		return fmt.Errorf("Java nao encontrado: %s. Verifique o PATH ou informe um caminho valido com --java-bin.", c.JavaBin)
	}

	info, err := os.Stat(c.JarPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("JAR nao encontrado: %s. Gere o artefato do modulo Java ou informe um caminho valido com --jar.", c.JarPath)
		}
		return fmt.Errorf("Nao foi possivel acessar o JAR informado em --jar: %s", c.JarPath)
	}
	if info.IsDir() {
		return fmt.Errorf("O valor informado em --jar precisa apontar para um arquivo, nao para um diretorio: %s", c.JarPath)
	}

	return nil
}

func (c Config) jarArgs(args []string) []string {
	jarArgs := make([]string, 0, len(args)+2)
	jarArgs = append(jarArgs, "-jar", c.JarPath)
	jarArgs = append(jarArgs, args...)
	return jarArgs
}

func exitCodeFromWait(err error) int {
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	return 1
}

func readFirstJSONObject(reader io.Reader) (string, error) {
	buffered := bufio.NewReader(reader)
	var output bytes.Buffer
	started := false
	depth := 0
	inString := false
	escaped := false

	for {
		b, err := buffered.ReadByte()
		if err != nil {
			if !started {
				return "", err
			}
			return output.String(), err
		}

		if !started {
			if isSpaceByte(b) {
				continue
			}
			if b != '{' {
				continue
			}

			started = true
			depth = 1
			output.WriteByte(b)
			continue
		}

		output.WriteByte(b)

		if inString {
			if escaped {
				escaped = false
				continue
			}
			if b == '\\' {
				escaped = true
				continue
			}
			if b == '"' {
				inString = false
			}
			continue
		}

		switch b {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return output.String(), nil
			}
		}
	}
}

func isSpaceByte(b byte) bool {
	switch b {
	case ' ', '\n', '\r', '\t':
		return true
	default:
		return false
	}
}
