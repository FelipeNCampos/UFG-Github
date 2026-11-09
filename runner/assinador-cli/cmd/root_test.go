package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestAssinarRejectsMissingAlias(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "entrada.json")
	outputPath := filepath.Join(tempDir, "saida.json")
	if err := os.WriteFile(inputPath, []byte(`{"resourceType":"Bundle"}`), 0o644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}

	_, _, err := executeRootCommand(
		"assinar",
		"--entrada", inputPath,
		"--saida", outputPath,
		"--modo", "direto",
	)
	if err == nil {
		t.Fatalf("expected validation error for missing alias")
	}

	var exitErr *ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected ExitError, got %T", err)
	}
	if exitErr.Code != 2 {
		t.Fatalf("expected exit code 2, got %d", exitErr.Code)
	}
	if !strings.Contains(exitErr.Message, "--alias") {
		t.Fatalf("unexpected error message: %s", exitErr.Message)
	}
}

func TestAssinarMapsPortugueseFlagsToJarArguments(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "entrada.json")
	outputPath := filepath.Join(tempDir, "saida.json")
	jarPath := filepath.Join(tempDir, "assinador-verificador.jar")
	if err := os.WriteFile(inputPath, []byte(`{"resourceType":"Bundle"}`), 0o644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to write jar: %v", err)
	}

	javaPath, logPath := createFakeJavaForCommandTests(t)
	t.Setenv("FAKE_JAVA_BEHAVIOR", "run_success")
	t.Setenv("FAKE_JAVA_LOG", logPath)

	stdout, stderr, err := executeRootCommand(
		"assinar",
		"--entrada", inputPath,
		"--saida", outputPath,
		"--modo", "direto",
		"--alias", "demo-signer",
		"--biblioteca-pkcs11", "token.dll",
		"--slot-pkcs11", "7",
		"--java-bin", javaPath,
		"--jar", jarPath,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr, got: %s", stderr)
	}
	if !strings.Contains(stdout, "\"status\":\"SUCCESS\"") {
		t.Fatalf("unexpected stdout: %s", stdout)
	}

	logBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read fake java log: %v", err)
	}
	logContent := string(logBytes)
	expectedFragments := []string{
		"-jar " + jarPath,
		"sign",
		"--pathin " + inputPath,
		"--pathout " + outputPath,
		"--mode direct",
		"--alias demo-signer",
		"--pkcs11-lib token.dll",
		"--pkcs11-slot 7",
	}
	for _, fragment := range expectedFragments {
		if !strings.Contains(logContent, fragment) {
			t.Fatalf("expected log to contain %q, got: %s", fragment, logContent)
		}
	}
}

func TestAssinarUsesManagedJavaWhenJavaBinIsNotProvided(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "entrada.json")
	outputPath := filepath.Join(tempDir, "saida.json")
	jarPath := filepath.Join(tempDir, "assinador-verificador.jar")
	if err := os.WriteFile(inputPath, []byte(`{"resourceType":"Bundle"}`), 0o644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to write jar: %v", err)
	}

	logPath := provisionManagedJavaForCommandTests(t)
	t.Setenv("FAKE_JAVA_BEHAVIOR", "run_success")
	t.Setenv("FAKE_JAVA_LOG", logPath)

	stdout, stderr, err := executeRootCommand(
		"assinar",
		"--entrada", inputPath,
		"--saida", outputPath,
		"--modo", "direto",
		"--alias", "demo-signer",
		"--jar", jarPath,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr, got: %s", stderr)
	}
	if !strings.Contains(stdout, "\"status\":\"SUCCESS\"") {
		t.Fatalf("unexpected stdout: %s", stdout)
	}

	logBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read managed java log: %v", err)
	}
	logContent := string(logBytes)
	if !strings.Contains(logContent, "--mode direct") || !strings.Contains(logContent, "--alias demo-signer") {
		t.Fatalf("unexpected managed java log content: %s", logContent)
	}
}

func TestVerificarRejectsUnsupportedAliasFlag(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "entrada.json")
	outputPath := filepath.Join(tempDir, "saida.json")
	if err := os.WriteFile(inputPath, []byte(`{"signature":"sig","inputDigestSha256":"digest"}`), 0o644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}

	_, _, err := executeRootCommand(
		"verificar",
		"--entrada", inputPath,
		"--saida", outputPath,
		"--modo", "direto",
		"--alias", "nao-permitido",
	)
	if err == nil {
		t.Fatalf("expected validation error for alias in verificar")
	}

	var exitErr *ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected ExitError, got %T", err)
	}
	if exitErr.Code != 2 {
		t.Fatalf("expected exit code 2, got %d", exitErr.Code)
	}
	if !strings.Contains(exitErr.Message, "nao aceita --alias") {
		t.Fatalf("unexpected error message: %s", exitErr.Message)
	}
}

func TestAssinarHttpMapsPortToJarArguments(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "entrada.json")
	outputPath := filepath.Join(tempDir, "saida.json")
	jarPath := filepath.Join(tempDir, "assinador-verificador.jar")
	if err := os.WriteFile(inputPath, []byte(`{"resourceType":"Bundle"}`), 0o644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to write jar: %v", err)
	}

	javaPath, logPath := createFakeJavaForCommandTests(t)
	t.Setenv("FAKE_JAVA_BEHAVIOR", "run_success")
	t.Setenv("FAKE_JAVA_LOG", logPath)

	_, _, err := executeRootCommand(
		"assinar",
		"--entrada", inputPath,
		"--saida", outputPath,
		"--modo", "http",
		"--porta", "8123",
		"--alias", "demo-signer",
		"--java-bin", javaPath,
		"--jar", jarPath,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read fake java log: %v", err)
	}
	logContent := string(logBytes)
	if !strings.Contains(logContent, "--mode http") || !strings.Contains(logContent, "--port 8123") {
		t.Fatalf("unexpected log content: %s", logContent)
	}
}

func TestVerificarHttpMapsPortToJarArguments(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "entrada.json")
	outputPath := filepath.Join(tempDir, "saida.json")
	jarPath := filepath.Join(tempDir, "assinador-verificador.jar")
	if err := os.WriteFile(inputPath, []byte(`{"signature":"sig","inputDigestSha256":"digest"}`), 0o644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to write jar: %v", err)
	}

	javaPath, logPath := createFakeJavaForCommandTests(t)
	t.Setenv("FAKE_JAVA_BEHAVIOR", "run_success")
	t.Setenv("FAKE_JAVA_LOG", logPath)

	_, _, err := executeRootCommand(
		"verificar",
		"--entrada", inputPath,
		"--saida", outputPath,
		"--modo", "http",
		"--porta", "9001",
		"--java-bin", javaPath,
		"--jar", jarPath,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read fake java log: %v", err)
	}
	logContent := string(logBytes)
	if !strings.Contains(logContent, "validate") || !strings.Contains(logContent, "--mode http") || !strings.Contains(logContent, "--port 9001") {
		t.Fatalf("unexpected log content: %s", logContent)
	}
}

func TestPropagatesJarExitCodeAndStderr(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "entrada.json")
	outputPath := filepath.Join(tempDir, "saida.json")
	jarPath := filepath.Join(tempDir, "assinador-verificador.jar")
	if err := os.WriteFile(inputPath, []byte(`{"signature":"sig","inputDigestSha256":"digest"}`), 0o644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to write jar: %v", err)
	}

	javaPath, _ := createFakeJavaForCommandTests(t)
	t.Setenv("FAKE_JAVA_BEHAVIOR", "run_error")

	stdout, stderr, err := executeRootCommand(
		"verificar",
		"--entrada", inputPath,
		"--saida", outputPath,
		"--modo", "direto",
		"--java-bin", javaPath,
		"--jar", jarPath,
	)
	if stdout != "" {
		t.Fatalf("expected empty stdout, got: %s", stdout)
	}
	if !strings.Contains(stderr, "\"VALIDATION_ERROR\"") {
		t.Fatalf("expected stderr to contain jar error output, got: %s", stderr)
	}

	var exitErr *ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected ExitError, got %T", err)
	}
	if exitErr.Code != 2 {
		t.Fatalf("expected propagated exit code 2, got %d", exitErr.Code)
	}
}

func TestServidorIniciarReturnsStartupJSON(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "assinador-verificador.jar")
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to write jar: %v", err)
	}

	javaPath, _ := createFakeJavaForCommandTests(t)
	t.Setenv("FAKE_JAVA_BEHAVIOR", "start_success")

	stdout, stderr, err := executeRootCommand(
		"servidor",
		"iniciar",
		"--porta", "8080",
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
}

func TestServidorStatusMapsToJarArguments(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "assinador-verificador.jar")
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to write jar: %v", err)
	}

	javaPath, logPath := createFakeJavaForCommandTests(t)
	t.Setenv("FAKE_JAVA_BEHAVIOR", "run_success")
	t.Setenv("FAKE_JAVA_LOG", logPath)

	_, _, err := executeRootCommand(
		"servidor",
		"status",
		"--porta", "9090",
		"--java-bin", javaPath,
		"--jar", jarPath,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read fake java log: %v", err)
	}
	logContent := string(logBytes)
	if !strings.Contains(logContent, "server status --port 9090") {
		t.Fatalf("unexpected log content: %s", logContent)
	}
}

func TestServidorPararMapsToJarArguments(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "assinador-verificador.jar")
	if err := os.WriteFile(jarPath, []byte("fake-jar"), 0o644); err != nil {
		t.Fatalf("failed to write jar: %v", err)
	}

	javaPath, logPath := createFakeJavaForCommandTests(t)
	t.Setenv("FAKE_JAVA_BEHAVIOR", "run_success")
	t.Setenv("FAKE_JAVA_LOG", logPath)

	_, _, err := executeRootCommand(
		"servidor",
		"parar",
		"--porta", "8088",
		"--java-bin", javaPath,
		"--jar", jarPath,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read fake java log: %v", err)
	}
	logContent := string(logBytes)
	if !strings.Contains(logContent, "server stop --port 8088") {
		t.Fatalf("unexpected log content: %s", logContent)
	}
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

func createFakeJavaForCommandTests(t *testing.T) (string, string) {
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
		"esac\n" +
		"printf '{\"status\":\"SUCCESS\"}\\n'\n"
	if err := os.WriteFile(javaPath, []byte(script), 0o755); err != nil {
		t.Fatalf("failed to write fake java script: %v", err)
	}
	return javaPath, logPath
}

func provisionManagedJavaForCommandTests(t *testing.T) string {
	t.Helper()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "managed-java.log")
	runnerHome := filepath.Join(tempDir, "runner")

	switch runtime.GOOS {
	case "windows":
		t.Setenv("APPDATA", tempDir)
		javaDir := filepath.Join(runnerHome, "java", "temurin-jre-17", "windows-x64", "bin")
		if err := os.MkdirAll(javaDir, 0o755); err != nil {
			t.Fatalf("failed to create managed runtime dir: %v", err)
		}
		script := "@echo off\r\n" +
			"if \"%1\"==\"-version\" exit /b 0\r\n" +
			"if not \"%FAKE_JAVA_LOG%\"==\"\" echo %*>>\"%FAKE_JAVA_LOG%\"\r\n" +
			"if \"%FAKE_JAVA_BEHAVIOR%\"==\"run_success\" (\r\n" +
			"  echo {\"status\":\"SUCCESS\",\"operation\":\"sign\"}\r\n" +
			"  exit /b 0\r\n" +
			")\r\n" +
			"echo {\"status\":\"SUCCESS\"}\r\n"
		if err := os.WriteFile(filepath.Join(javaDir, "java.cmd"), []byte(script), 0o755); err != nil {
			t.Fatalf("failed to write managed java script: %v", err)
		}
	default:
		t.Setenv("XDG_CONFIG_HOME", tempDir)
		javaDir := filepath.Join(runnerHome, "java", "temurin-jre-17", runtime.GOOS+"-x64", "bin")
		if err := os.MkdirAll(javaDir, 0o755); err != nil {
			t.Fatalf("failed to create managed runtime dir: %v", err)
		}
		script := "#!/bin/sh\n" +
			"if [ \"$1\" = \"-version\" ]; then exit 0; fi\n" +
			"if [ -n \"$FAKE_JAVA_LOG\" ]; then printf '%s\\n' \"$*\" >> \"$FAKE_JAVA_LOG\"; fi\n" +
			"printf '{\"status\":\"SUCCESS\",\"operation\":\"sign\"}\\n'\n"
		if err := os.WriteFile(filepath.Join(javaDir, "java"), []byte(script), 0o755); err != nil {
			t.Fatalf("failed to write managed java script: %v", err)
		}
	}

	return logPath
}
