package provisioner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// DefaultReleaseURL é sobrescrito via SIMULADOR_RELEASE_URL para testes e ambientes customizados.
const DefaultReleaseURL = "https://raw.githubusercontent.com/ufgfelipencampos/runner/main/release.json"

type ReleaseInfo struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

type localVersion struct {
	Version string `json:"version"`
}

// EnsureJar garante que o simulador.jar em jarPath existe e está na versão mais recente.
//
// Comportamento:
//   - JAR atualizado → nada é feito.
//   - JAR ausente ou desatualizado → download automático.
//   - Falha de rede + JAR presente → aviso, continua com versão local.
//   - Falha de rede + JAR ausente → erro acionável.
func EnsureJar(jarPath string, out io.Writer) error {
	remote, err := fetchRelease(releaseURL())
	if err != nil {
		if _, statErr := os.Stat(jarPath); statErr == nil {
			fmt.Fprintf(out, "Aviso: nao foi possivel verificar atualizacao do simulador.jar (%v). Usando versao local.\n", err)
			return nil
		}
		return fmt.Errorf("Falha ao verificar versao do simulador.jar: %v. Verifique sua conexao e tente novamente.", err)
	}

	if isUpToDate(jarPath, remote.Version) {
		return nil
	}

	fmt.Fprintf(out, "Baixando simulador.jar v%s...\n", remote.Version)
	if err := downloadFile(remote.URL, jarPath); err != nil {
		return fmt.Errorf("Falha ao baixar simulador.jar: %v. Verifique sua conexao e tente novamente.", err)
	}

	if err := saveVersion(jarPath, remote.Version); err != nil {
		fmt.Fprintf(out, "Aviso: nao foi possivel salvar informacao de versao: %v\n", err)
	}

	fmt.Fprintf(out, "simulador.jar v%s pronto.\n", remote.Version)
	return nil
}

func releaseURL() string {
	if url := os.Getenv("SIMULADOR_RELEASE_URL"); url != "" {
		return url
	}
	return DefaultReleaseURL
}

func fetchRelease(url string) (*ReleaseInfo, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("servidor retornou status %d ao buscar release.json", resp.StatusCode)
	}

	var info ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("release.json invalido: %v", err)
	}
	if info.Version == "" || info.URL == "" {
		return nil, fmt.Errorf("release.json incompleto: campos 'version' e 'url' sao obrigatorios")
	}
	return &info, nil
}

func isUpToDate(jarPath, remoteVersion string) bool {
	data, err := os.ReadFile(versionFilePath(jarPath))
	if err != nil {
		return false
	}
	var local localVersion
	if err := json.Unmarshal(data, &local); err != nil {
		return false
	}
	if local.Version != remoteVersion {
		return false
	}
	_, err = os.Stat(jarPath)
	return err == nil
}

func downloadFile(url, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("nao foi possivel criar diretorio de destino: %v", err)
	}

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download retornou status %d", resp.StatusCode)
	}

	tmp := destPath + ".download"
	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("nao foi possivel criar arquivo temporario: %v", err)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("falha durante download: %v", err)
	}
	f.Close()

	return os.Rename(tmp, destPath)
}

func saveVersion(jarPath, version string) error {
	data, err := json.Marshal(localVersion{Version: version})
	if err != nil {
		return err
	}
	return os.WriteFile(versionFilePath(jarPath), data, 0o644)
}

func versionFilePath(jarPath string) string {
	return filepath.Join(filepath.Dir(jarPath), "simulador-version.json")
}
