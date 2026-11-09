package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// --- integração: ciclo completo status → parar → status ---

func TestLifecycleStatusPararStatus(t *testing.T) {
	shutdownReceived := make(chan struct{})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/info":
			fmt.Fprint(w, `{"status":"ONLINE","version":"test"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/shutdown":
			fmt.Fprint(w, `{"status":"SUCCESS","message":"Simulador encerrado."}`)
			go func() { close(shutdownReceived) }()
		}
	}))
	defer srv.Close()

	port := srv.Listener.Addr().(*net.TCPAddr).Port
	portStr := strconv.Itoa(port)

	// 1. status com servidor ativo → ONLINE
	stdout, _, err := executeRootCommand("status", "--porta", portStr)
	if err != nil {
		t.Fatalf("status: unexpected error: %v", err)
	}
	if !strings.Contains(stdout, "ONLINE") {
		t.Fatalf("status: expected ONLINE, got: %s", stdout)
	}

	// 2. parar → envia POST /shutdown
	_, _, err = executeRootCommand("parar", "--porta", portStr)
	if err != nil {
		t.Fatalf("parar: unexpected error: %v", err)
	}
	select {
	case <-shutdownReceived:
	case <-time.After(2 * time.Second):
		t.Fatalf("parar: /shutdown nao foi chamado dentro do timeout")
	}

	// 3. simula desligamento: fecha o servidor fake
	srv.Close()

	// 4. status sem servidor → OFFLINE (sem erro de processo)
	stdout, _, err = executeRootCommand("status", "--porta", portStr)
	if err != nil {
		t.Fatalf("status offline: unexpected error: %v", err)
	}
	if !strings.Contains(stdout, "OFFLINE") {
		t.Fatalf("status offline: expected OFFLINE, got: %s", stdout)
	}
}

// --- iniciar ---

func TestIniciarRejectsPortAlreadyInUse(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	_, _, err = executeRootCommand("iniciar", "--porta", strconv.Itoa(port))
	assertValidationError(t, err, "ja esta em uso")
}

func TestIniciarRejectsPortZero(t *testing.T) {
	_, _, err := executeRootCommand("iniciar", "--porta", "0")
	assertValidationError(t, err, "--porta")
}

func TestIniciarRejectsPortAboveRange(t *testing.T) {
	_, _, err := executeRootCommand("iniciar", "--porta", "99999")
	assertValidationError(t, err, "--porta")
}

func TestIniciarStartsJarWithCorrectPort(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "simulador.jar")
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to write jar: %v", err)
	}

	// Obtém uma porta livre e a libera imediatamente para o comando usar
	freePort := findFreePort(t)

	javaPath, logPath := createFakeJavaForCommandTests(t)
	t.Setenv("FAKE_JAVA_BEHAVIOR", "start_success")
	t.Setenv("FAKE_JAVA_LOG", logPath)

	portStr := strconv.Itoa(freePort)
	stdout, stderr, err := executeRootCommand(
		"iniciar",
		"--porta", portStr,
		"--java-bin", javaPath,
		"--jar", jarPath,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr, got: %s", stderr)
	}
	if !strings.Contains(stdout, "\"operation\": \"server-start\"") {
		t.Fatalf("unexpected stdout: %s", stdout)
	}

	logBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read fake java log: %v", err)
	}
	logContent := string(logBytes)
	if !strings.Contains(logContent, "server start") || !strings.Contains(logContent, "--port "+portStr) {
		t.Fatalf("unexpected log content: %s", logContent)
	}
}

func TestIniciarPropagatesJarError(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "simulador.jar")
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to write jar: %v", err)
	}

	javaPath, _ := createFakeJavaForCommandTests(t)
	t.Setenv("FAKE_JAVA_BEHAVIOR", "start_error")

	_, stderr, err := executeRootCommand(
		"iniciar",
		"--java-bin", javaPath,
		"--jar", jarPath,
	)

	var exitErr *ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected ExitError, got %T: %v", err, err)
	}
	if exitErr.Code == 0 {
		t.Fatalf("expected non-zero exit code")
	}
	_ = stderr
}

// --- parar ---

func TestPararRejectsPortZero(t *testing.T) {
	_, _, err := executeRootCommand("parar", "--porta", "0")
	assertValidationError(t, err, "--porta")
}

func TestPararSendsPostToShutdown(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/shutdown" {
			called = true
			fmt.Fprint(w, `{"status":"SUCCESS","message":"Simulador encerrado."}`)
		}
	}))
	defer srv.Close()

	port := srv.Listener.Addr().(*net.TCPAddr).Port

	stdout, _, err := executeRootCommand("parar", "--porta", strconv.Itoa(port))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected POST /shutdown to be called")
	}
	if !strings.Contains(stdout, "SUCCESS") {
		t.Fatalf("unexpected stdout: %s", stdout)
	}
}

func TestPararReturnsErrorWhenSimuladorNotRunning(t *testing.T) {
	_, _, err := executeRootCommand("parar", "--porta", "59992")
	if err == nil {
		t.Fatalf("expected error when simulador is not running")
	}
	var exitErr *ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected ExitError, got %T", err)
	}
	if exitErr.Code != 1 {
		t.Fatalf("expected exit code 1, got %d", exitErr.Code)
	}
}

// --- status ---

func TestStatusRejectsPortZero(t *testing.T) {
	_, _, err := executeRootCommand("status", "--porta", "0")
	assertValidationError(t, err, "--porta")
}

func TestStatusReturnsOnlineWhenRunning(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/info" {
			fmt.Fprint(w, `{"status":"ONLINE","version":"1.0.0"}`)
		}
	}))
	defer srv.Close()

	port := srv.Listener.Addr().(*net.TCPAddr).Port

	stdout, _, err := executeRootCommand("status", "--porta", strconv.Itoa(port))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout, "ONLINE") {
		t.Fatalf("unexpected stdout: %s", stdout)
	}
}

func TestStatusReturnsOfflineWhenNotRunning(t *testing.T) {
	stdout, _, err := executeRootCommand("status", "--porta", "59993")
	if err != nil {
		t.Fatalf("unexpected error (status should not error on offline): %v", err)
	}
	if !strings.Contains(stdout, "OFFLINE") {
		t.Fatalf("expected OFFLINE in stdout, got: %s", stdout)
	}
}

// --- helpers ---

// findFreePort abre e fecha imediatamente um listener para obter uma porta livre.
func findFreePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func executeRootCommand(args ...string) (string, string, error) {
	rootCmd := NewRootCommand()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stderr)
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	return stdout.String(), stderr.String(), err
}

func assertValidationError(t *testing.T, err error, fragment string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected validation error")
	}
	var exitErr *ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected ExitError, got %T", err)
	}
	if exitErr.Code != 2 {
		t.Fatalf("expected exit code 2, got %d", exitErr.Code)
	}
	if fragment != "" && !strings.Contains(exitErr.Message, fragment) {
		t.Fatalf("expected error message to contain %q, got: %s", fragment, exitErr.Message)
	}
}

func createFakeJavaForCommandTests(t *testing.T) (string, string) {
	t.Helper()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "fake-java.log")

	if runtime.GOOS == "windows" {
		javaPath := filepath.Join(tempDir, "java.cmd")
		script := "@echo off\r\n" +
			"if not \"%FAKE_JAVA_LOG%\"==\"\" echo %*>>\"%FAKE_JAVA_LOG%\"\r\n" +
			"if \"%FAKE_JAVA_BEHAVIOR%\"==\"start_success\" (\r\n" +
			"  echo {\r\n" +
			"  echo   \"status\": \"SUCCESS\",\r\n" +
			"  echo   \"operation\": \"server-start\",\r\n" +
			"  echo   \"port\": 8443\r\n" +
			"  echo }\r\n" +
			"  powershell -NoProfile -Command \"Start-Sleep -Seconds 3\" >NUL\r\n" +
			"  exit /b 0\r\n" +
			")\r\n" +
			"if \"%FAKE_JAVA_BEHAVIOR%\"==\"start_error\" (\r\n" +
			"  >&2 echo {\"status\":\"ERROR\",\"type\":\"RUNTIME_ERROR\"}\r\n" +
			"  exit /b 1\r\n" +
			")\r\n" +
			"echo {\"status\":\"SUCCESS\"}\r\n"
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
		"    cat <<'EOF'\n" +
		"{\n" +
		"  \"status\": \"SUCCESS\",\n" +
		"  \"operation\": \"server-start\",\n" +
		"  \"port\": 8443\n" +
		"}\n" +
		"EOF\n" +
		"    sleep 3\n" +
		"    exit 0\n" +
		"    ;;\n" +
		"  start_error)\n" +
		"    printf '{\"status\":\"ERROR\",\"type\":\"RUNTIME_ERROR\"}\\n' >&2\n" +
		"    exit 1\n" +
		"    ;;\n" +
		"esac\n" +
		"printf '{\"status\":\"SUCCESS\"}\\n'\n"
	if err := os.WriteFile(javaPath, []byte(script), 0o755); err != nil {
		t.Fatalf("failed to write fake java script: %v", err)
	}
	return javaPath, logPath
}
