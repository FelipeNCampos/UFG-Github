package javaruntime

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	temurinReleaseName = "jdk-17.0.18+8"
	defaultBaseURL     = "https://api.adoptium.net/v3"
	metadataFileName   = ".runner-java.json"
)

type Resolver struct {
	BaseURL    string
	Client     *http.Client
	HomeDir    string
	GOOS       string
	GOARCH     string
	Now        func() time.Time
	HTTPGetter func(*http.Request) (*http.Response, error)
}

type installMetadata struct {
	Distribution string    `json:"distribution"`
	Release      string    `json:"release"`
	OS           string    `json:"os"`
	Arch         string    `json:"arch"`
	InstalledAt  time.Time `json:"installedAt"`
}

func NewDefaultResolver() Resolver {
	return Resolver{
		BaseURL: defaultBaseURL,
		Client: &http.Client{
			Timeout: 2 * time.Minute,
		},
		GOOS:   runtime.GOOS,
		GOARCH: runtime.GOARCH,
		Now:    time.Now,
	}
}

func (r Resolver) Resolve(explicitJavaBin string) (string, error) {
	if strings.TrimSpace(explicitJavaBin) != "" {
		javaBin, err := exec.LookPath(explicitJavaBin)
		if err != nil {
			return "", fmt.Errorf("Java explicito invalido: %s. Verifique o valor de --java-bin.", explicitJavaBin)
		}
		return javaBin, nil
	}

	javaBin, err := r.ensureManagedRuntime()
	if err != nil {
		return "", err
	}
	return javaBin, nil
}

func (r Resolver) ensureManagedRuntime() (string, error) {
	r = r.withDefaults()

	rootDir, err := r.runnerHomeDir()
	if err != nil {
		return "", fmt.Errorf("Nao foi possivel resolver o diretorio local do Runner: %w", err)
	}

	runtimeDir := filepath.Join(rootDir, "java", "temurin-jre-17", r.platformKey())
	javaBin := r.existingJavaBinary(runtimeDir)
	if err := validateJavaBinary(javaBin); err == nil {
		return javaBin, nil
	}

	if err := r.installManagedRuntime(runtimeDir); err != nil {
		return "", err
	}

	javaBin = r.existingJavaBinary(runtimeDir)
	if err := validateJavaBinary(javaBin); err != nil {
		return "", fmt.Errorf("O runtime Java provisionado ficou invalido ou corrompido. Remova %s e tente novamente.", runtimeDir)
	}

	return javaBin, nil
}

func (r Resolver) installManagedRuntime(runtimeDir string) error {
	baseDir := filepath.Dir(runtimeDir)
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return fmt.Errorf("Nao foi possivel preparar o diretorio de runtimes do Runner: %w", err)
	}

	tempDir, err := os.MkdirTemp(baseDir, "install-*")
	if err != nil {
		return fmt.Errorf("Nao foi possivel criar um diretorio temporario para instalar o Java: %w", err)
	}
	defer os.RemoveAll(tempDir)

	archivePath, err := r.downloadRuntimeArchive(tempDir)
	if err != nil {
		return err
	}

	extractDir := filepath.Join(tempDir, "extract")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return fmt.Errorf("Nao foi possivel preparar a extracao do runtime Java: %w", err)
	}

	if err := extractArchive(archivePath, extractDir); err != nil {
		return fmt.Errorf("Falha ao extrair o runtime Java gerenciado: %w", err)
	}

	javaBin, err := findJavaBinary(extractDir, r.GOOS)
	if err != nil {
		return fmt.Errorf("Falha ao localizar o executavel Java extraido: %w", err)
	}

	javaHome := filepath.Dir(filepath.Dir(javaBin))
	if r.GOOS == "darwin" && strings.EqualFold(filepath.Base(javaHome), "Home") {
		javaHome = filepath.Clean(javaHome)
	}

	if err := r.writeMetadata(javaHome); err != nil {
		return err
	}

	if err := os.RemoveAll(runtimeDir); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("Nao foi possivel limpar o runtime Java anterior: %w", err)
	}

	if err := os.Rename(javaHome, runtimeDir); err != nil {
		return fmt.Errorf("Falha ao promover o runtime Java provisionado: %w", err)
	}

	return nil
}

func (r Resolver) downloadRuntimeArchive(tempDir string) (string, error) {
	downloadURL, err := r.binaryURL()
	if err != nil {
		return "", err
	}
	checksumURL, err := r.checksumURL()
	if err != nil {
		return "", err
	}

	checksum, err := r.fetchChecksum(checksumURL)
	if err != nil {
		return "", err
	}

	archivePath := filepath.Join(tempDir, "temurin"+r.archiveExtension())
	if err := r.downloadFile(downloadURL, archivePath); err != nil {
		return "", err
	}

	if err := verifySHA256(archivePath, checksum); err != nil {
		return "", fmt.Errorf("Falha ao validar a integridade do runtime Java baixado: %w", err)
	}

	return archivePath, nil
}

func (r Resolver) fetchChecksum(rawURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("Nao foi possivel preparar a requisicao do checksum do Java: %w", err)
	}

	response, err := r.do(req)
	if err != nil {
		return "", fmt.Errorf("Falha ao baixar o checksum do runtime Java: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Falha ao baixar o checksum do runtime Java: status %d", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Falha ao ler o checksum do runtime Java: %w", err)
	}

	checksum := strings.TrimSpace(string(data))
	if checksum == "" {
		return "", fmt.Errorf("O checksum do runtime Java veio vazio")
	}

	return checksum, nil
}

func (r Resolver) downloadFile(rawURL string, destination string) error {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return fmt.Errorf("Nao foi possivel preparar o download do runtime Java: %w", err)
	}

	response, err := r.do(req)
	if err != nil {
		return fmt.Errorf("Falha ao baixar o runtime Java gerenciado: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Falha ao baixar o runtime Java gerenciado: status %d", response.StatusCode)
	}

	file, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("Nao foi possivel salvar o runtime Java baixado: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		return fmt.Errorf("Falha ao salvar o runtime Java baixado: %w", err)
	}

	return nil
}

func (r Resolver) writeMetadata(javaHome string) error {
	metadata := installMetadata{
		Distribution: "Temurin JRE",
		Release:      temurinReleaseName,
		OS:           r.GOOS,
		Arch:         r.GOARCH,
		InstalledAt:  r.Now().UTC(),
	}

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("Nao foi possivel preparar os metadados do runtime Java: %w", err)
	}

	if err := os.WriteFile(filepath.Join(javaHome, metadataFileName), append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("Nao foi possivel gravar os metadados do runtime Java: %w", err)
	}

	return nil
}

func (r Resolver) runnerHomeDir() (string, error) {
	if strings.TrimSpace(r.HomeDir) != "" {
		return r.HomeDir, nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "runner"), nil
}

func (r Resolver) platformKey() string {
	return fmt.Sprintf("%s-%s", r.GOOS, platformArch(r.GOARCH))
}

func (r Resolver) javaBinaryPath(runtimeDir string) string {
	fileName := "java"
	if r.GOOS == "windows" {
		fileName = "java.exe"
	}
	return filepath.Join(runtimeDir, "bin", fileName)
}

func (r Resolver) existingJavaBinary(runtimeDir string) string {
	for _, candidate := range javaBinaryCandidates(runtimeDir, r.GOOS) {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return r.javaBinaryPath(runtimeDir)
}

func (r Resolver) binaryURL() (string, error) {
	osName, archName, err := r.adoptiumPlatform()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/binary/version/%s/%s/%s/jre/hotspot/normal/eclipse", strings.TrimRight(r.BaseURL, "/"), url.PathEscape(temurinReleaseName), osName, archName), nil
}

func (r Resolver) checksumURL() (string, error) {
	osName, archName, err := r.adoptiumPlatform()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/checksum/version/%s/%s/%s/jre/hotspot/normal/eclipse", strings.TrimRight(r.BaseURL, "/"), url.PathEscape(temurinReleaseName), osName, archName), nil
}

func (r Resolver) adoptiumPlatform() (string, string, error) {
	var osName string
	switch r.GOOS {
	case "windows":
		osName = "windows"
	case "linux":
		osName = "linux"
	case "darwin":
		osName = "mac"
	default:
		return "", "", fmt.Errorf("Sistema operacional nao suportado para provisionamento automatico de Java: %s", r.GOOS)
	}

	switch r.GOARCH {
	case "amd64":
		return osName, "x64", nil
	default:
		return "", "", fmt.Errorf("Arquitetura nao suportada para provisionamento automatico de Java: %s", r.GOARCH)
	}
}

func (r Resolver) archiveExtension() string {
	if r.GOOS == "windows" {
		return ".zip"
	}
	return ".tar.gz"
}

func (r Resolver) do(req *http.Request) (*http.Response, error) {
	if r.HTTPGetter != nil {
		return r.HTTPGetter(req)
	}
	return r.Client.Do(req)
}

func (r Resolver) withDefaults() Resolver {
	if r.Client == nil {
		r.Client = &http.Client{Timeout: 2 * time.Minute}
	}
	if r.BaseURL == "" {
		r.BaseURL = defaultBaseURL
	}
	if r.GOOS == "" {
		r.GOOS = runtime.GOOS
	}
	if r.GOARCH == "" {
		r.GOARCH = runtime.GOARCH
	}
	if r.Now == nil {
		r.Now = time.Now
	}
	return r
}

func validateJavaBinary(javaBin string) error {
	if _, err := os.Stat(javaBin); err != nil {
		return err
	}

	command := exec.Command(javaBin, "-version")
	if err := command.Run(); err != nil {
		return err
	}

	return nil
}

func verifySHA256(path string, expectedChecksum string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	actual := hex.EncodeToString(hash.Sum(nil))
	if !strings.EqualFold(actual, strings.TrimSpace(expectedChecksum)) {
		return fmt.Errorf("checksum esperado %s, obtido %s", expectedChecksum, actual)
	}
	return nil
}

func extractArchive(archivePath string, destination string) error {
	if strings.HasSuffix(strings.ToLower(archivePath), ".zip") {
		return extractZip(archivePath, destination)
	}
	return extractTarGz(archivePath, destination)
}

func extractZip(archivePath string, destination string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		targetPath, err := safeJoin(destination, file.Name)
		if err != nil {
			return err
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return err
		}

		source, err := file.Open()
		if err != nil {
			return err
		}

		target, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
		if err != nil {
			source.Close()
			return err
		}

		_, copyErr := io.Copy(target, source)
		closeErr := errors.Join(target.Close(), source.Close())
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
	}

	return nil
}

func extractTarGz(archivePath string, destination string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	reader := tar.NewReader(gzipReader)
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		targetPath, err := safeJoin(destination, header.Name)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				return err
			}
			target, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(target, reader); err != nil {
				target.Close()
				return err
			}
			if err := target.Close(); err != nil {
				return err
			}
		}
	}
}

func safeJoin(baseDir string, relativePath string) (string, error) {
	cleanPath := filepath.Clean(relativePath)
	if cleanPath == "." {
		return baseDir, nil
	}

	targetPath := filepath.Join(baseDir, cleanPath)
	relative, err := filepath.Rel(baseDir, targetPath)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(relative, "..") {
		return "", fmt.Errorf("o arquivo compactado tentou escrever fora do diretorio de destino: %s", relativePath)
	}
	return targetPath, nil
}

func findJavaBinary(rootDir string, goos string) (string, error) {
	var found string
	walkErr := filepath.WalkDir(rootDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if isJavaBinaryName(entry.Name(), goos) && strings.EqualFold(filepath.Base(filepath.Dir(path)), "bin") {
			found = path
			return io.EOF
		}
		return nil
	})
	if walkErr != nil && !errors.Is(walkErr, io.EOF) {
		return "", walkErr
	}
	if found == "" {
		return "", fmt.Errorf("nenhum executavel Java foi encontrado no runtime extraido")
	}
	return found, nil
}

func platformArch(goarch string) string {
	switch goarch {
	case "amd64":
		return "x64"
	default:
		return goarch
	}
}

func javaBinaryCandidates(runtimeDir string, goos string) []string {
	binDir := filepath.Join(runtimeDir, "bin")
	if goos == "windows" {
		return []string{
			filepath.Join(binDir, "java.exe"),
			filepath.Join(binDir, "java.cmd"),
		}
	}
	return []string{filepath.Join(binDir, "java")}
}

func isJavaBinaryName(name string, goos string) bool {
	if goos == "windows" {
		return strings.EqualFold(name, "java.exe") || strings.EqualFold(name, "java.cmd")
	}
	return name == "java"
}
