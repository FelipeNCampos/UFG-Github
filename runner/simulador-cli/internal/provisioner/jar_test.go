package provisioner

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// startMockServers sobe dois httptest.Server: um para release.json, outro para o JAR.
// Retorna os dois servidores e o conteúdo de release.json apontando para o servidor do JAR.
func startMockServers(t *testing.T, version string, jarContent []byte) (releaseSrv, jarSrv *httptest.Server) {
	t.Helper()

	jarSrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/java-archive")
		w.Write(jarContent)
	}))

	release := ReleaseInfo{Version: version, URL: jarSrv.URL + "/simulador.jar"}
	payload, _ := json.Marshal(release)

	releaseSrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	}))

	return releaseSrv, jarSrv
}

func TestEnsureJarDownloadsWhenAbsent(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "simulador.jar")
	fakeJarContent := []byte("fake-jar-content")

	releaseSrv, jarSrv := startMockServers(t, "1.0.0", fakeJarContent)
	defer releaseSrv.Close()
	defer jarSrv.Close()

	t.Setenv("SIMULADOR_RELEASE_URL", releaseSrv.URL)

	var out bytes.Buffer
	if err := EnsureJar(jarPath, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(jarPath)
	if err != nil {
		t.Fatalf("expected JAR to be downloaded, got error: %v", err)
	}
	if string(content) != string(fakeJarContent) {
		t.Fatalf("unexpected JAR content: %q", content)
	}
	if !strings.Contains(out.String(), "1.0.0") {
		t.Fatalf("expected version in output, got: %s", out.String())
	}
}

func TestEnsureJarSkipsDownloadWhenUpToDate(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "simulador.jar")
	if err := os.WriteFile(jarPath, []byte("original-content"), 0o644); err != nil {
		t.Fatalf("failed to write jar: %v", err)
	}
	if err := saveVersion(jarPath, "1.0.0"); err != nil {
		t.Fatalf("failed to save version: %v", err)
	}

	downloadCalled := false
	jarSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downloadCalled = true
		w.Write([]byte("new-content"))
	}))
	defer jarSrv.Close()

	release := ReleaseInfo{Version: "1.0.0", URL: jarSrv.URL + "/simulador.jar"}
	payload, _ := json.Marshal(release)
	releaseSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(payload)
	}))
	defer releaseSrv.Close()

	t.Setenv("SIMULADOR_RELEASE_URL", releaseSrv.URL)

	var out bytes.Buffer
	if err := EnsureJar(jarPath, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if downloadCalled {
		t.Fatalf("expected no download when JAR is up to date")
	}
	content, _ := os.ReadFile(jarPath)
	if string(content) != "original-content" {
		t.Fatalf("JAR should not have been replaced: %q", content)
	}
}

func TestEnsureJarDownloadsWhenOutdated(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "simulador.jar")
	if err := os.WriteFile(jarPath, []byte("old-content"), 0o644); err != nil {
		t.Fatalf("failed to write jar: %v", err)
	}
	if err := saveVersion(jarPath, "0.9.0"); err != nil {
		t.Fatalf("failed to save version: %v", err)
	}

	newContent := []byte("new-jar-content")
	releaseSrv, jarSrv := startMockServers(t, "1.0.0", newContent)
	defer releaseSrv.Close()
	defer jarSrv.Close()

	t.Setenv("SIMULADOR_RELEASE_URL", releaseSrv.URL)

	var out bytes.Buffer
	if err := EnsureJar(jarPath, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, _ := os.ReadFile(jarPath)
	if string(content) != string(newContent) {
		t.Fatalf("expected JAR to be updated, got: %q", content)
	}

	versionData, _ := os.ReadFile(versionFilePath(jarPath))
	if !strings.Contains(string(versionData), "1.0.0") {
		t.Fatalf("expected version cache to be updated, got: %s", versionData)
	}
}

func TestEnsureJarContinuesWithWarningOnNetworkFailureWhenJarExists(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "simulador.jar")
	if err := os.WriteFile(jarPath, []byte("existing-jar"), 0o644); err != nil {
		t.Fatalf("failed to write jar: %v", err)
	}

	// Aponta para URL inexistente para simular falha de rede
	t.Setenv("SIMULADOR_RELEASE_URL", "http://localhost:1")

	var out bytes.Buffer
	if err := EnsureJar(jarPath, &out); err != nil {
		t.Fatalf("expected no error when jar exists and network fails, got: %v", err)
	}

	if !strings.Contains(out.String(), "Aviso") {
		t.Fatalf("expected warning in output, got: %s", out.String())
	}
	if !strings.Contains(out.String(), "versao local") {
		t.Fatalf("expected 'versao local' in output, got: %s", out.String())
	}
}

func TestEnsureJarErrorsOnNetworkFailureWhenJarAbsent(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "simulador.jar")

	t.Setenv("SIMULADOR_RELEASE_URL", "http://localhost:1")

	var out bytes.Buffer
	err := EnsureJar(jarPath, &out)
	if err == nil {
		t.Fatalf("expected error when network fails and jar is absent")
	}
	if !strings.Contains(err.Error(), "Verifique sua conexao") {
		t.Fatalf("expected actionable error, got: %v", err)
	}
}

func TestEnsureJarErrorsWhenReleaseJSONIsMalformed(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "simulador.jar")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	t.Setenv("SIMULADOR_RELEASE_URL", srv.URL)

	var out bytes.Buffer
	err := EnsureJar(jarPath, &out)
	if err == nil {
		t.Fatalf("expected error on malformed release.json")
	}
	if !strings.Contains(err.Error(), "release.json invalido") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEnsureJarErrorsWhenReleaseJSONMissingFields(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "simulador.jar")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"version":"1.0.0"}`)) // sem "url"
	}))
	defer srv.Close()

	t.Setenv("SIMULADOR_RELEASE_URL", srv.URL)

	var out bytes.Buffer
	err := EnsureJar(jarPath, &out)
	if err == nil {
		t.Fatalf("expected error when release.json is missing 'url'")
	}
	if !strings.Contains(err.Error(), "incompleto") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEnsureJarErrorsWhenDownloadFails(t *testing.T) {
	tempDir := t.TempDir()
	jarPath := filepath.Join(tempDir, "simulador.jar")

	// Release aponta para URL de download que retorna 500
	downloadSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer downloadSrv.Close()

	release := ReleaseInfo{Version: "1.0.0", URL: downloadSrv.URL + "/simulador.jar"}
	payload, _ := json.Marshal(release)
	releaseSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(payload)
	}))
	defer releaseSrv.Close()

	t.Setenv("SIMULADOR_RELEASE_URL", releaseSrv.URL)

	var out bytes.Buffer
	err := EnsureJar(jarPath, &out)
	if err == nil {
		t.Fatalf("expected error when download fails")
	}
	if !strings.Contains(err.Error(), "Falha ao baixar") {
		t.Fatalf("unexpected error: %v", err)
	}
}
