package runner

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	JavaBin        string
	JarPath        string
	HealthCheckURL string // substitui https://localhost:{port}/actuator/health; usado em testes
}

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	PID      int
}

// IsPortInUse verifica via TCP se algo já está escutando na porta.
func IsPortInUse(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Stop encerra o simulador lendo o PID salvo em ~/.hubsaude/simulador-{port}.pid e matando o processo.
func Stop(port int) (Result, error) {
	pid, err := loadPID(port)
	if err != nil {
		return Result{}, fmt.Errorf(
			"Nao foi possivel localizar o PID do simulador na porta %d. Verifique se ele esta em execucao com: simulador-cli status --porta %d",
			port, port,
		)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		removePID(port)
		return Result{}, fmt.Errorf("Processo com PID %d nao encontrado. O simulador pode ja ter sido encerrado.", pid)
	}

	if err := process.Kill(); err != nil {
		removePID(port)
		return Result{}, fmt.Errorf("Nao foi possivel encerrar o simulador (PID %d): %v", pid, err)
	}

	removePID(port)

	response := fmt.Sprintf(
		`{"status": "SUCCESS", "operation": "server-stop", "port": %d, "message": "Simulador encerrado.", "pid": %d}`,
		port, pid,
	)
	return Result{Stdout: response}, nil
}

// Status consulta GET /api/info no simulador via HTTPS.
// Três estados possíveis:
//   - OFFLINE : nenhum processo responde na porta
//   - RUNNING : processo responde, mas /api/info retornou status não-2xx
//   - (corpo) : /api/info retornou 2xx — repassa o JSON do simulador
func Status(port int) (Result, error) {
	url := fmt.Sprintf("https://localhost:%d/api/info", port)

	client := httpsClient(5 * time.Second)
	resp, err := client.Get(url)
	if err != nil {
		offline := fmt.Sprintf(
			`{"status": "OFFLINE", "port": %d, "message": "Simulador nao responde na porta %d."}`,
			port, port,
		)
		return Result{Stdout: offline}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		running := fmt.Sprintf(
			`{"status": "RUNNING", "port": %d, "message": "Servidor ativo na porta %d, mas GET /api/info retornou %d."}`,
			port, port, resp.StatusCode,
		)
		return Result{Stdout: running}, nil
	}

	body, _ := io.ReadAll(resp.Body)
	return Result{Stdout: string(body)}, nil
}

// StartServer inicia o simulador.jar em segundo plano e aguarda confirmacao via health check HTTPS
// em /actuator/health. Salva o PID em ~/.hubsaude/simulador-{port}.pid ao confirmar inicializacao.
func (c Config) StartServer(port int, args []string) (Result, error) {
	if err := c.validate(); err != nil {
		return Result{}, err
	}

	command := exec.Command(c.JavaBin, c.jarArgs(args)...)
	stderrPipe, err := command.StderrPipe()
	if err != nil {
		return Result{}, fmt.Errorf("Nao foi possivel preparar a captura do stderr do simulador: %w", err)
	}

	if err := command.Start(); err != nil {
		return Result{}, fmt.Errorf("Falha ao iniciar o Java. Verifique o valor de --java-bin e tente novamente: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	stderrCh := make(chan string, 1)
	waitCh := make(chan error, 1)
	healthCh := make(chan error, 1)

	go func() {
		data, _ := io.ReadAll(stderrPipe)
		stderrCh <- string(data)
	}()
	go func() {
		waitCh <- command.Wait()
	}()
	go func() {
		healthCh <- pollHealth(ctx, c.healthURL(port))
	}()

	select {
	case err := <-healthCh:
		if err != nil {
			if command.Process != nil {
				_ = command.Process.Kill()
			}
			waitErr := <-waitCh
			stderr := <-stderrCh
			return Result{
				Stderr:   stderr,
				ExitCode: exitCodeFrom(waitErr),
			}, fmt.Errorf("Simulador nao ficou disponivel: %v", err)
		}
		pid := command.Process.Pid
		if saveErr := savePID(port, pid); saveErr != nil {
			fmt.Fprintf(os.Stderr, "Aviso: nao foi possivel salvar PID: %v\n", saveErr)
		}
		response := fmt.Sprintf(
			`{"status": "SUCCESS", "operation": "server-start", "port": %d, "message": "Simulador iniciado.", "pid": %d}`,
			port, pid,
		)
		return Result{Stdout: response, PID: pid}, nil

	case waitErr := <-waitCh:
		cancel()
		stderr := <-stderrCh
		result := Result{
			Stderr:   stderr,
			ExitCode: exitCodeFrom(waitErr),
		}
		if result.ExitCode == 0 {
			return Result{}, fmt.Errorf("O processo Java encerrou inesperadamente sem iniciar o simulador.")
		}
		return result, nil
	}
}

func (c Config) healthURL(port int) string {
	if c.HealthCheckURL != "" {
		return c.HealthCheckURL
	}
	return fmt.Sprintf("https://localhost:%d/actuator/health", port)
}

func pollHealth(ctx context.Context, url string) error {
	client := httpsClient(3 * time.Second)
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("erro ao preparar requisicao de health check: %v", err)
		}
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return nil
			}
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout aguardando simulador ficar disponivel")
		case <-time.After(500 * time.Millisecond):
		}
	}
}

func httpsClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		},
	}
}

func pidFilePath(port int) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Sprintf("simulador-%d.pid", port)
	}
	return filepath.Join(home, ".hubsaude", fmt.Sprintf("simulador-%d.pid", port))
}

func savePID(port, pid int) error {
	path := pidFilePath(port)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0o644)
}

func loadPID(port int) (int, error) {
	data, err := os.ReadFile(pidFilePath(port))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

func removePID(port int) {
	os.Remove(pidFilePath(port))
}

func (c Config) validate() error {
	if strings.TrimSpace(c.JavaBin) == "" {
		return fmt.Errorf("Nenhum executavel Java foi informado. Use --java-bin para apontar para o Java instalado.")
	}
	if strings.TrimSpace(c.JarPath) == "" {
		return fmt.Errorf("Nenhum caminho de JAR foi informado. Use --jar para apontar para o simulador.jar.")
	}

	if _, err := exec.LookPath(c.JavaBin); err != nil {
		return fmt.Errorf("Java nao encontrado: %s. Verifique o PATH ou informe um caminho valido com --java-bin.", c.JavaBin)
	}

	info, err := os.Stat(c.JarPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("simulador.jar nao encontrado: %s. Execute 'simulador-cli iniciar' apos provisionar o JAR ou informe um caminho valido com --jar.", c.JarPath)
		}
		return fmt.Errorf("Nao foi possivel acessar o JAR informado em --jar: %s", c.JarPath)
	}
	if info.IsDir() {
		return fmt.Errorf("O valor informado em --jar precisa apontar para um arquivo, nao para um diretorio: %s", c.JarPath)
	}

	return nil
}

func (c Config) jarArgs(args []string) []string {
	result := make([]string, 0, len(args)+2)
	result = append(result, "-jar", c.JarPath)
	result = append(result, args...)
	return result
}

func exitCodeFrom(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return 1
}
