package runner

import (
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRunReturnsHelpfulErrorWhenJavaIsMissing(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "assinador-verificador.jar")
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to create fake jar: %v", err)
	}

	_, err := Config{
		JavaBin: filepath.Join(tempDir, "missing-java"),
		JarPath: jarPath,
	}.Run([]string{"server", "status"})
	if err == nil {
		t.Fatalf("expected a helpful error when java is missing")
	}
	if !strings.Contains(err.Error(), "Java nao encontrado") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunReturnsHelpfulErrorWhenJarIsMissing(t *testing.T) {
	tempDir := t.TempDir()
	javaPath, _ := createFakeJava(t)

	_, err := Config{
		JavaBin: javaPath,
		JarPath: filepath.Join(tempDir, "missing.jar"),
	}.Run([]string{"server", "status"})
	if err == nil {
		t.Fatalf("expected a helpful error when jar is missing")
	}
	if !strings.Contains(err.Error(), "JAR nao encontrado") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStartServerReturnsStartupJSON(t *testing.T) {
	tempDir := t.TempDir()
	javaPath, _ := createFakeJava(t)
	jarPath := filepath.Join(tempDir, "assinador-verificador.jar")
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to create fake jar: %v", err)
	}

	t.Setenv("FAKE_JAVA_BEHAVIOR", "start_success")

	result, err := Config{
		JavaBin: javaPath,
		JarPath: jarPath,
	}.StartServer([]string{"server", "start", "--port", "8080"})
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

func createFakeJava(t *testing.T) (string, string) {
	t.Helper()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "fake-java.log")

	if runtime.GOOS == "windows" {
		javaPath := filepath.Join(tempDir, "java.cmd")
		script := "@echo off\r\n" +
			"if not \"%FAKE_JAVA_LOG%\"==\"\" echo %*>>\"%FAKE_JAVA_LOG%\"\r\n" +
			"if \"%FAKE_JAVA_BEHAVIOR%\"==\"run_success\" (\r\n" +
			"  echo {\"status\":\"SUCCESS\",\"operation\":\"sign\"}\r\n" +
			"  exit /b 0\r\n" +
			")\r\n" +
			"if \"%FAKE_JAVA_BEHAVIOR%\"==\"run_error\" (\r\n" +
			"  >&2 echo {\"status\":\"ERROR\",\"type\":\"VALIDATION_ERROR\"}\r\n" +
			"  exit /b 2\r\n" +
			")\r\n" +
			"if \"%FAKE_JAVA_BEHAVIOR%\"==\"start_success\" (\r\n" +
			"  echo {\r\n" +
			"  echo   \"status\": \"SUCCESS\",\r\n" +
			"  echo   \"operation\": \"server-start\",\r\n" +
			"  echo   \"port\": 8080\r\n" +
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
		"  run_success)\n" +
		"    printf '{\"status\":\"SUCCESS\",\"operation\":\"sign\"}\\n'\n" +
		"    exit 0\n" +
		"    ;;\n" +
		"  run_error)\n" +
		"    printf '{\"status\":\"ERROR\",\"type\":\"VALIDATION_ERROR\"}\\n' >&2\n" +
		"    exit 2\n" +
		"    ;;\n" +
		"  start_success)\n" +
		"    cat <<'EOF'\n" +
		"{\n" +
		"  \"status\": \"SUCCESS\",\n" +
		"  \"operation\": \"server-start\",\n" +
		"  \"port\": 8080\n" +
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

// Testes para ExecutionStrategy e modo automático
func TestIsServerAvailableWhenListening(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("falha ao criar listener: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port

	if !IsServerAvailable(port, 1) {
		t.Fatalf("esperava que servidor estivesse disponivel na porta %d", port)
	}
}

func TestIsServerAvailableWhenNotListening(t *testing.T) {
	// Usa uma porta que dificilmente estará em uso
	if IsServerAvailable(59999, 1) {
		t.Fatalf("esperava que servidor nao estivesse disponivel na porta 59999")
	}
}

func TestDetermineExecutionModeDirect(t *testing.T) {
	mode, err := DetermineExecutionMode("direct", 8080)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if mode != "direct" {
		t.Fatalf("esperava 'direct', obteve %q", mode)
	}
}

func TestDetermineExecutionModeHTTP(t *testing.T) {
	mode, err := DetermineExecutionMode("http", 8080)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if mode != "http" {
		t.Fatalf("esperava 'http', obteve %q", mode)
	}
}

func TestDetermineExecutionModeAutoWithServer(t *testing.T) {
	listener, _ := net.Listen("tcp", "localhost:0")
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port

	mode, err := DetermineExecutionMode("auto", port)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if mode != "http" {
		t.Fatalf("esperava 'http' com servidor disponivel, obteve %q", mode)
	}
}

func TestDetermineExecutionModeAutoWithoutServer(t *testing.T) {
	mode, err := DetermineExecutionMode("auto", 59999)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if mode != "direct" {
		t.Fatalf("esperava 'direct' sem servidor, obteve %q", mode)
	}
}

func TestApplyExecutionModeAddsMode(t *testing.T) {
	args := []string{"sign", "--pathin", "e.json", "--pathout", "s.json"}
	result := ApplyExecutionMode(args, "http", 8080)

	hasMode := false
	hasPort := false
	for i := 0; i < len(result)-1; i++ {
		if result[i] == "--mode" && result[i+1] == "http" {
			hasMode = true
		}
		if result[i] == "--port" && result[i+1] == "8080" {
			hasPort = true
		}
	}

	if !hasMode {
		t.Fatalf("flag --mode nao foi adicionada aos argumentos: %v", result)
	}
	if !hasPort {
		t.Fatalf("flag --port nao foi adicionada para modo http: %v", result)
	}
}

func TestApplyExecutionModeRemovesExistingMode(t *testing.T) {
	args := []string{"sign", "--mode", "direct", "--pathin", "e.json"}
	result := ApplyExecutionMode(args, "http", 8080)

	count := 0
	for i := 0; i < len(result); i++ {
		if result[i] == "--mode" {
			count++
		}
	}

	if count != 1 {
		t.Fatalf("esperava exatamente uma flag --mode, encontrou %d em: %v", count, result)
	}
}

func TestRemoveModeFromArgs(t *testing.T) {
	args := []string{"sign", "--mode", "direct", "--pathin", "e.json", "--mode", "http"}
	result := RemoveModeFromArgs(args)

	for i := 0; i < len(result); i++ {
		if result[i] == "--mode" {
			t.Fatalf("flag --mode nao foi removida: %v", result)
		}
	}
}
