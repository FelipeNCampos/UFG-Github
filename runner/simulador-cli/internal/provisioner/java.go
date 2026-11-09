package provisioner

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const adoptiumBaseURL = "https://api.adoptium.net"

// EnsureJava garante que um executável Java compatível está disponível e retorna
// o caminho resolvido para ele. Ordem de busca:
//  1. javaBin via LookPath (resolve PATH e caminhos absolutos)
//  2. Se javaBin tem separador de caminho e não foi encontrado → erro (usuário deu path explícito inválido)
//  3. JRE já provisionado em ~/.hubsaude/jre/
//  4. Download automático do Temurin 21 via Adoptium API
func EnsureJava(javaBin string, out io.Writer) (string, error) {
	if path, err := exec.LookPath(javaBin); err == nil {
		return path, nil
	}

	if strings.ContainsAny(javaBin, `/\`) {
		return "", fmt.Errorf("Java nao encontrado em: %s. Verifique o caminho informado em --java-bin.", javaBin)
	}

	jreBin := provisionedJREBin()
	if _, err := os.Stat(jreBin); err == nil {
		return jreBin, nil
	}

	fmt.Fprintln(out, "Java nao encontrado. Baixando JRE Temurin 21...")
	if err := downloadJRE(jreBin, out); err != nil {
		return "", fmt.Errorf("Falha ao provisionar Java: %v. Instale Java manualmente ou tente novamente.", err)
	}
	return jreBin, nil
}

func provisionedJREBin() string {
	jreHome := os.Getenv("JAVA_JRE_HOME")
	if jreHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "java"
		}
		jreHome = filepath.Join(home, ".hubsaude", "jre")
	}
	if runtime.GOOS == "windows" {
		return filepath.Join(jreHome, "bin", "java.exe")
	}
	return filepath.Join(jreHome, "bin", "java")
}

func downloadJRE(jreBin string, out io.Writer) error {
	goos, arch, ext, err := currentPlatform()
	if err != nil {
		return err
	}

	base := adoptiumBaseURL
	if u := os.Getenv("JAVA_ADOPTIUM_BASE_URL"); u != "" {
		base = u
	}

	url := fmt.Sprintf("%s/v3/binary/latest/21/ga/%s/%s/jre/hotspot/normal/eclipse", base, goos, arch)
	fmt.Fprintf(out, "Baixando Temurin 21 JRE para %s/%s...\n", goos, arch)

	archivePath, err := fetchArchive(url, ext)
	if err != nil {
		return err
	}
	defer os.Remove(archivePath)

	// jreBin = ~/.hubsaude/jre/bin/java → jreDir = ~/.hubsaude/jre
	jreDir := filepath.Dir(filepath.Dir(jreBin))
	if err := extractJRE(archivePath, jreDir, ext); err != nil {
		return err
	}

	fmt.Fprintln(out, "JRE Temurin 21 provisionado com sucesso.")
	return nil
}

func fetchArchive(url, ext string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download retornou status %d", resp.StatusCode)
	}

	tmp, err := os.CreateTemp("", "temurin-*."+ext)
	if err != nil {
		return "", fmt.Errorf("nao foi possivel criar arquivo temporario: %v", err)
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("falha durante download: %v", err)
	}
	return tmp.Name(), nil
}

func extractJRE(archivePath, jreDir, ext string) error {
	tmp, err := os.MkdirTemp("", "jre-extract-*")
	if err != nil {
		return fmt.Errorf("nao foi possivel criar diretorio temporario: %v", err)
	}
	defer os.RemoveAll(tmp)

	if ext == "zip" {
		if err := extractZip(archivePath, tmp); err != nil {
			return err
		}
	} else {
		if err := extractTarGz(archivePath, tmp); err != nil {
			return err
		}
	}

	entries, err := os.ReadDir(tmp)
	if err != nil || len(entries) == 0 {
		return fmt.Errorf("arquivo extraido vazio ou invalido")
	}
	extracted := filepath.Join(tmp, entries[0].Name())

	if err := os.RemoveAll(jreDir); err != nil {
		return fmt.Errorf("nao foi possivel limpar diretorio de destino: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(jreDir), 0o755); err != nil {
		return fmt.Errorf("nao foi possivel criar diretorio pai: %v", err)
	}
	return os.Rename(extracted, jreDir)
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("nao foi possivel abrir arquivo zip: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, filepath.FromSlash(f.Name))
		if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("caminho inseguro no zip: %s", f.Name)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}
		_, copyErr := io.Copy(out, rc)
		out.Close()
		rc.Close()
		if copyErr != nil {
			return copyErr
		}
	}
	return nil
}

func extractTarGz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("nao foi possivel ler gzip: %v", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erro ao ler tar: %v", err)
		}

		path := filepath.Join(dest, filepath.FromSlash(hdr.Name))
		if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("caminho inseguro no tar: %s", hdr.Name)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(path, os.FileMode(hdr.Mode))
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			out, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(out, tr)
			out.Close()
			if copyErr != nil {
				return copyErr
			}
		}
	}
	return nil
}

func currentPlatform() (string, string, string, error) {
	switch runtime.GOOS {
	case "windows":
		return "windows", "x64", "zip", nil
	case "linux":
		return "linux", "x64", "tar.gz", nil
	case "darwin":
		return "mac", "x64", "tar.gz", nil
	default:
		return "", "", "", fmt.Errorf("plataforma nao suportada: %s/%s", runtime.GOOS, runtime.GOARCH)
	}
}
