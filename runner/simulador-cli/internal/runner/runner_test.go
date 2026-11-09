package runner

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// --- IsPortInUse ---

func TestIsPortInUseWhenListening(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	if !IsPortInUse(port) {
		t.Fatalf("expected port %d to be reported as in use", port)
	}
}

func TestIsPortInUseWhenFree(t *testing.T) {
	if IsPortInUse(59980) {
		t.Fatalf("expected port 59980 to be reported as free")
	}
}

// --- StartServer ---

func TestStartServerReturnsStartupJSON(t *testing.T) {
	tempDir := t.TempDir()
	javaPath, _ := createFakeJava(t)
	jarPath := filepath.Join(tempDir, "simulador.jar")
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to create fake jar: %v", err)
	}

	t.Setenv("FAKE_JAVA_BEHAVIOR", "start_success")

	// Servidor TLS mock que representa o health endpoint do simulador
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	port := srv.Listener.Addr().(*net.TCPAddr).Port
	t.Cleanup(func() { removePID(port) })

	result, err := Config{
		JavaBin:        javaPath,
		JarPath:        jarPath,
		HealthCheckURL: srv.URL + "/actuator/health",
	}.StartServer(port, []string{"server", "start", "--port", strconv.Itoa(port)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("expected successful startup, got exit code %d", result.ExitCode)
	}
	if !strings.Contains(result.Stdout, "\"operation\": \"server-start\"") {
		t.Fatalf("unexpected startup output: %s", result.Stdout)
	}
	if result.PID == 0 {
		t.Fatalf("expected a child process pid for the background server")
	}
}

func TestStartServerReturnsErrorWhenJarFails(t *testing.T) {
	tempDir := t.TempDir()
	javaPath, _ := createFakeJava(t)
	jarPath := filepath.Join(tempDir, "simulador.jar")
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to create fake jar: %v", err)
	}

	t.Setenv("FAKE_JAVA_BEHAVIOR", "start_error")

	// Processo encerra com erro antes do health check ter chance de responder
	result, err := Config{
		JavaBin:        javaPath,
		JarPath:        jarPath,
		HealthCheckURL: "https://localhost:59999/actuator/health",
	}.StartServer(8443, []string{"server", "start", "--port", "8443"})
	if err != nil {
		t.Fatalf("unexpected hard error: %v", err)
	}
	if result.ExitCode == 0 {
		t.Fatalf("expected non-zero exit code on jar failure")
	}
}

func TestStartServerKillsProcessWhenHealthCheckFails(t *testing.T) {
	tempDir := t.TempDir()
	javaPath, _ := createFakeJava(t)
	jarPath := filepath.Join(tempDir, "simulador.jar")
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to create fake jar: %v", err)
	}

	t.Setenv("FAKE_JAVA_BEHAVIOR", "start_success")

	start := time.Now()
	result, err := Config{
		JavaBin:        javaPath,
		JarPath:        jarPath,
		HealthCheckURL: "://invalid-health-url",
	}.StartServer(8443, []string{"server", "start", "--port", "8443"})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatalf("expected health check error")
	}
	if result.ExitCode == 0 {
		t.Fatalf("expected non-zero exit code after killing unhealthy process")
	}
	if elapsed > 2*time.Second {
		t.Fatalf("expected unhealthy process to be killed quickly, took %s", elapsed)
	}
}

// --- validate ---

func TestValidateReturnsHelpfulErrorWhenJavaIsMissing(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "simulador.jar")
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to create fake jar: %v", err)
	}

	_, err := Config{
		JavaBin: filepath.Join(tempDir, "missing-java"),
		JarPath: jarPath,
	}.StartServer(8443, []string{"server", "start", "--port", "8443"})
	if err == nil {
		t.Fatalf("expected error when java is missing")
	}
	if !strings.Contains(err.Error(), "Java nao encontrado") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateReturnsHelpfulErrorWhenJarIsMissing(t *testing.T) {
	tempDir := t.TempDir()
	javaPath, _ := createFakeJava(t)

	_, err := Config{
		JavaBin: javaPath,
		JarPath: filepath.Join(tempDir, "missing.jar"),
	}.StartServer(8443, []string{"server", "start", "--port", "8443"})
	if err == nil {
		t.Fatalf("expected error when jar is missing")
	}
	if !strings.Contains(err.Error(), "simulador.jar nao encontrado") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Stop ---

func TestStopKillsProcessByPID(t *testing.T) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "ping", "-n", "100", "127.0.0.1")
	} else {
		cmd = exec.Command("sleep", "100")
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start background process: %v", err)
	}

	port := 59985
	if err := savePID(port, cmd.Process.Pid); err != nil {
		t.Fatalf("failed to save PID: %v", err)
	}
	t.Cleanup(func() { removePID(port) })

	result, err := Stop(port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result.Stdout, "SUCCESS") {
		t.Fatalf("expected SUCCESS in result, got: %s", result.Stdout)
	}

	if err := cmd.Wait(); err == nil {
		t.Fatalf("expected process to be killed, but it exited cleanly")
	}

	if _, statErr := os.Stat(pidFilePath(port)); !os.IsNotExist(statErr) {
		t.Fatalf("expected PID file to be removed after stop")
	}
}

func TestStopReturnsErrorWhenSimuladorNotRunning(t *testing.T) {
	port := 59990
	removePID(port) // garante estado limpo

	_, err := Stop(port)
	if err == nil {
		t.Fatalf("expected error when simulador is not running")
	}
	if !strings.Contains(err.Error(), strconv.Itoa(port)) {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Status ---

func TestStatusReturnsOnlineWhenRunning(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/info" {
			fmt.Fprint(w, `{"status":"ONLINE","version":"1.0.0"}`)
		}
	}))
	defer srv.Close()

	port := srv.Listener.Addr().(*net.TCPAddr).Port

	result, err := Status(port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result.Stdout, "ONLINE") {
		t.Fatalf("unexpected result: %s", result.Stdout)
	}
}

func TestStatusReturnsRunningWhenEndpointNotFound(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>404 Not Found</h1>")
	}))
	defer srv.Close()

	port := srv.Listener.Addr().(*net.TCPAddr).Port

	result, err := Status(port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result.Stdout, "RUNNING") {
		t.Fatalf("expected RUNNING when server is up but endpoint missing, got: %s", result.Stdout)
	}
	if !strings.Contains(result.Stdout, "404") {
		t.Fatalf("expected status code 404 in message, got: %s", result.Stdout)
	}
}

func TestStatusReturnsOfflineWhenNotRunning(t *testing.T) {
	result, err := Status(59991)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result.Stdout, "OFFLINE") {
		t.Fatalf("expected OFFLINE in result, got: %s", result.Stdout)
	}
	if !strings.Contains(result.Stdout, "59991") {
		t.Fatalf("expected port in result, got: %s", result.Stdout)
	}
}

// --- validate (caminho é diretório) ---

func TestValidateReturnsErrorWhenJarIsDirectory(t *testing.T) {
	tempDir := t.TempDir()
	javaPath, _ := createFakeJava(t)

	_, err := Config{
		JavaBin: javaPath,
		JarPath: tempDir,
	}.StartServer(8443, []string{"server", "start", "--port", "8443"})
	if err == nil {
		t.Fatalf("expected error when jar path is a directory")
	}
	if !strings.Contains(err.Error(), "diretorio") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

// --- jarArgs ---

func TestJarArgsContainsJarFlagAndPath(t *testing.T) {
	c := Config{JavaBin: "java", JarPath: "/caminho/simulador.jar"}
	args := c.jarArgs([]string{"server", "start", "--port", "8443"})

	if len(args) < 4 {
		t.Fatalf("expected at least 4 args, got %d: %v", len(args), args)
	}
	if args[0] != "-jar" {
		t.Fatalf("expected args[0]==-jar, got %q", args[0])
	}
	if args[1] != "/caminho/simulador.jar" {
		t.Fatalf("expected args[1]==jar path, got %q", args[1])
	}
	if args[2] != "server" || args[3] != "start" {
		t.Fatalf("expected command args after jar path, got %v", args[2:])
	}
}

// --- helpers ---

func createFakeJava(t *testing.T) (string, string) {
	t.Helper()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "fake-java.log")

	if runtime.GOOS == "windows" {
		javaPath := filepath.Join(tempDir, "java.cmd")
		script := "@echo off\r\n" +
			"if not \"%FAKE_JAVA_LOG%\"==\"\" echo %*>>\"%FAKE_JAVA_LOG%\"\r\n" +
			"if \"%FAKE_JAVA_BEHAVIOR%\"==\"start_success\" (\r\n" +
			"  powershell -NoProfile -Command \"Start-Sleep -Seconds 3\" >NUL\r\n" +
			"  exit /b 0\r\n" +
			")\r\n" +
			"if \"%FAKE_JAVA_BEHAVIOR%\"==\"start_error\" (\r\n" +
			"  >&2 echo {\"status\":\"ERROR\",\"type\":\"RUNTIME_ERROR\"}\r\n" +
			"  exit /b 1\r\n" +
			")\r\n" +
			"powershell -NoProfile -Command \"Start-Sleep -Seconds 3\" >NUL\r\n"
		if err := os.WriteFile(javaPath, []byte(script), 0o755); err != nil {
			t.Fatalf("failed to write fake java script: %v", err)
		}
		return javaPath, logPath
	}

	javaPath := filepath.Join(tempDir, "java")
	script := "#!/bin/sh\n" +
		"if [ -n \"$FAKE_JAVA_LOG\" ]; then printf '%s\\n' \"$*\" >> \"$FAKE_JAVA_LOG\"; fi\n" +
		"case \"$FAKE_JAVA_BEHAVIOR\" in\n" +
		"  start_success)\n" +
		"    sleep 3\n" +
		"    exit 0\n" +
		"    ;;\n" +
		"  start_error)\n" +
		"    printf '{\"status\":\"ERROR\",\"type\":\"RUNTIME_ERROR\"}\\n' >&2\n" +
		"    exit 1\n" +
		"    ;;\n" +
		"esac\n" +
		"sleep 3\n"
	if err := os.WriteFile(javaPath, []byte(script), 0o755); err != nil {
		t.Fatalf("failed to write fake java script: %v", err)
	}
	return javaPath, logPath
}
