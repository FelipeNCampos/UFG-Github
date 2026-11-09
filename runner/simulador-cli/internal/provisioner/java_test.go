package provisioner

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// --- EnsureJava: binário explícito ---

func TestEnsureJavaReturnsExplicitBinWhenFound(t *testing.T) {
	javaPath, _ := createFakeJavaBin(t)

	var out bytes.Buffer
	got, err := EnsureJava(javaPath, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Fatalf("expected resolved path, got empty string")
	}
}

func TestEnsureJavaErrorsOnMissingExplicitBin(t *testing.T) {
	tempDir := t.TempDir()
	missingPath := filepath.Join(tempDir, "no", "such", "java")

	var out bytes.Buffer
	_, err := EnsureJava(missingPath, &out)
	if err == nil {
		t.Fatalf("expected error for missing explicit java-bin")
	}
	if !strings.Contains(err.Error(), "--java-bin") {
		t.Fatalf("expected --java-bin in error, got: %v", err)
	}
}

// --- EnsureJava: JRE provisionado ---

func TestEnsureJavaUsesAlreadyProvisionedJRE(t *testing.T) {
	jreHome := t.TempDir()
	jreBin := filepath.Join(jreHome, "bin", fakeJavaBinName())
	if err := os.MkdirAll(filepath.Dir(jreBin), 0o755); err != nil {
		t.Fatalf("failed to create bin dir: %v", err)
	}
	if err := os.WriteFile(jreBin, []byte("fake"), 0o755); err != nil {
		t.Fatalf("failed to write fake java: %v", err)
	}

	t.Setenv("JAVA_JRE_HOME", jreHome)
	// Limpa PATH para que exec.LookPath("java") falhe
	t.Setenv("PATH", t.TempDir())

	var out bytes.Buffer
	got, err := EnsureJava("java", &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != jreBin {
		t.Fatalf("expected %q, got %q", jreBin, got)
	}
	if out.Len() != 0 {
		t.Fatalf("expected no output when JRE already provisioned, got: %s", out.String())
	}
}

// --- EnsureJava: download automático ---

func TestEnsureJavaDownloadsAndExtractsTemurin(t *testing.T) {
	jreHome := t.TempDir()
	archiveContent := createFakeJREArchive(t, "jdk-21-jre")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(archiveContent)
	}))
	defer srv.Close()

	t.Setenv("JAVA_ADOPTIUM_BASE_URL", srv.URL)
	t.Setenv("JAVA_JRE_HOME", jreHome)
	t.Setenv("PATH", t.TempDir()) // java não está no PATH

	var out bytes.Buffer
	got, err := EnsureJava("java", &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBin := filepath.Join(jreHome, "bin", fakeJavaBinName())
	if got != expectedBin {
		t.Fatalf("expected %q, got %q", expectedBin, got)
	}
	if _, err := os.Stat(expectedBin); err != nil {
		t.Fatalf("expected java binary at %s after extraction: %v", expectedBin, err)
	}
	if !strings.Contains(out.String(), "provisionado com sucesso") {
		t.Fatalf("expected success message in output, got: %s", out.String())
	}
}

func TestEnsureJavaErrorsWhenDownloadFails(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	t.Setenv("JAVA_ADOPTIUM_BASE_URL", srv.URL)
	t.Setenv("JAVA_JRE_HOME", t.TempDir())
	t.Setenv("PATH", t.TempDir())

	var out bytes.Buffer
	_, err := EnsureJava("java", &out)
	if err == nil {
		t.Fatalf("expected error when download fails")
	}
	if !strings.Contains(err.Error(), "Falha ao provisionar") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- extractZip / extractTarGz ---

func TestExtractZipCreatesExpectedFiles(t *testing.T) {
	content := createFakeJREZip(t, "jdk-21-jre")
	archive := writeTempFile(t, content, "test-*.zip")
	dest := t.TempDir()

	if err := extractZip(archive, dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFile := filepath.Join(dest, "jdk-21-jre", "bin", fakeJavaBinName())
	if _, err := os.Stat(expectedFile); err != nil {
		t.Fatalf("expected file %s after extraction: %v", expectedFile, err)
	}
}

func TestExtractTarGzCreatesExpectedFiles(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("tar.gz extraction not used on Windows")
	}

	content := createFakeJRETarGz(t, "jdk-21-jre")
	archive := writeTempFile(t, content, "test-*.tar.gz")
	dest := t.TempDir()

	if err := extractTarGz(archive, dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFile := filepath.Join(dest, "jdk-21-jre", "bin", "java")
	if _, err := os.Stat(expectedFile); err != nil {
		t.Fatalf("expected file %s after extraction: %v", expectedFile, err)
	}
}

func TestExtractZipRejectsZipSlip(t *testing.T) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	fw, _ := w.Create("../../../etc/passwd")
	fw.Write([]byte("malicious"))
	w.Close()

	archive := writeTempFile(t, buf.Bytes(), "test-slip-*.zip")
	dest := t.TempDir()

	if err := extractZip(archive, dest); err == nil {
		t.Fatalf("expected error for zip-slip path")
	}
}

// --- helpers ---

func fakeJavaBinName() string {
	if runtime.GOOS == "windows" {
		return "java.exe"
	}
	return "java"
}

func createFakeJREArchive(t *testing.T, topDirName string) []byte {
	t.Helper()
	if runtime.GOOS == "windows" {
		return createFakeJREZip(t, topDirName)
	}
	return createFakeJRETarGz(t, topDirName)
}

func createFakeJREZip(t *testing.T, topDirName string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	binName := topDirName + "/bin/java.exe"
	fh := &zip.FileHeader{Name: binName, Method: zip.Deflate}
	fh.SetMode(0o755)
	fw, err := w.CreateHeader(fh)
	if err != nil {
		t.Fatalf("failed to create zip entry: %v", err)
	}
	fw.Write([]byte("fake-java"))
	w.Close()
	return buf.Bytes()
}

func createFakeJRETarGz(t *testing.T, topDirName string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	binContent := []byte("#!/bin/sh\necho java")
	tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     topDirName + "/bin/java",
		Mode:     0o755,
		Size:     int64(len(binContent)),
	})
	tw.Write(binContent)
	tw.Close()
	gz.Close()
	return buf.Bytes()
}

func createFakeJavaBin(t *testing.T) (string, string) {
	t.Helper()
	tempDir := t.TempDir()
	name := fakeJavaBinName()
	path := filepath.Join(tempDir, name)
	if err := os.WriteFile(path, []byte("fake"), 0o755); err != nil {
		t.Fatalf("failed to write fake java: %v", err)
	}
	return path, tempDir
}

func writeTempFile(t *testing.T, content []byte, pattern string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), pattern)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	f.Write(content)
	f.Close()
	return f.Name()
}
