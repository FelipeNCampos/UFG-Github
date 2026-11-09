# 🔍 REFERÊNCIA DE MUDANÇAS — O QUE MUDOU EM CADA ARQUIVO

**Data**: 19 de maio de 2026  
**Formato**: Antes vs. Depois para cada arquivo

---

## 📄 ARQUIVO 1: `cmd/common.go`

### ✏️ MUDANÇA 1: Imports

**ANTES**:
```go
import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/felip/runner/assinador-cli/internal/runner"
	"github.com/spf13/cobra"
)
```

**DEPOIS**:
```go
import (
	"fmt"
	"io"
	"net"           // ← NOVO
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"          // ← NOVO

	"github.com/felip/runner/assinador-cli/internal/runner"
	"github.com/spf13/cobra"
)
```

### ✏️ MUDANÇA 2: Novo tipo + funções (após `type runtimeFlags struct`)

**ADICIONADO**:
```go
// ExecutionStrategy define como o assinador.jar será invocado
type ExecutionStrategy string

const (
	StrategyAuto   ExecutionStrategy = "auto"   // Tenta HTTP, fallback para direto
	StrategyHTTP   ExecutionStrategy = "http"   // Força HTTP
	StrategyDirect ExecutionStrategy = "direct" // Força direto
)

func (s ExecutionStrategy) String() string {
	return string(s)
}

// ParseExecutionStrategy converte entrada do usuário para estratégia
// Entrada vazia padrão para StrategyAuto (novo comportamento)
func ParseExecutionStrategy(input string) (ExecutionStrategy, error) {
	normalized := strings.ToLower(strings.TrimSpace(input))

	// Padrão: string vazia mapeia para auto
	if normalized == "" {
		return StrategyAuto, nil
	}

	switch normalized {
	case "auto":
		return StrategyAuto, nil
	case "http":
		return StrategyHTTP, nil
	case "direto", "direct": // Suporta português e inglês
		return StrategyDirect, nil
	default:
		return "", validationError(
			"Estrategia de modo invalida: %s. Use 'auto', 'http' ou 'direto'.",
			input,
		)
	}
}
```

---

## 📄 ARQUIVO 2: `internal/runner/runner.go`

### ✏️ MUDANÇA 1: Imports

**ANTES**:
```go
import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)
```

**DEPOIS**:
```go
import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"      // ← NOVO
	"os"
	"os/exec"
	"strings"
	"time"     // ← NOVO
)
```

### ✏️ MUDANÇA 2: Cinco novas funções (ANTES de `func (c Config) Run()`)

**ADICIONADO**:
```go
// IsServerAvailable verifica se o servidor responde na porta dada
func IsServerAvailable(port int, timeoutSecs int) bool {
	address := fmt.Sprintf("localhost:%d", port)
	timeout := time.Duration(timeoutSecs) * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// DetermineExecutionMode decide o modo final baseado na estratégia
func DetermineExecutionMode(strategy string, port int) (string, error) {
	switch strings.ToLower(strategy) {
	case "direct":
		return "direct", nil
	case "http":
		return "http", nil
	case "auto":
		if IsServerAvailable(port, 1) {
			return "http", nil
		}
		return "direct", nil
	default:
		return "", fmt.Errorf("estrategia desconhecida: %s", strategy)
	}
}

// ApplyExecutionMode modifica os argumentos do comando
func ApplyExecutionMode(args []string, mode string, port int) []string {
	filtered := RemoveModeFromArgs(args)
	filtered = append(filtered, "--mode", mode)
	if mode == "http" {
		hasPort := false
		for j := 0; j < len(filtered); j++ {
			if filtered[j] == "--port" {
				hasPort = true
				break
			}
		}
		if !hasPort {
			filtered = append(filtered, "--port", fmt.Sprintf("%d", port))
		}
	}
	return filtered
}

// RemoveModeFromArgs remove qualquer flag --mode e seu valor
func RemoveModeFromArgs(args []string) []string {
	result := []string{}
	i := 0
	for i < len(args) {
		if args[i] == "--mode" {
			i += 2
			continue
		}
		result = append(result, args[i])
		i++
	}
	return result
}

// RunWithStrategy executa com seleção automática de modo
func (c Config) RunWithStrategy(args []string, strategy string, port int) (Result, error) {
	if err := c.Validate(); err != nil {
		return Result{}, err
	}
	mode, err := DetermineExecutionMode(strategy, port)
	if err != nil {
		return Result{}, err
	}
	finalArgs := ApplyExecutionMode(args, mode, port)
	return c.Run(finalArgs)
}
```

---

## 📄 ARQUIVO 3: `cmd/assinar.go`

### ✏️ MUDANÇA 1: Flag --modo (na função `newAssinarCmd()`)

**ANTES**:
```go
	flags.StringVar(&options.Modo, "modo", "", "Modo de execucao: direto ou http.")
```

**DEPOIS**:
```go
	flags.StringVar(&options.Modo, "modo", "auto",
		"Estrategia de execucao: auto, http ou direto.\n"+
			"  auto   (padrao): Tenta servidor HTTP, fallback para direto se indisponivel.\n"+
			"  http   : Usa apenas servidor HTTP, falha se indisponivel.\n"+
			"  direto : Usa apenas execucao direta (sem servidor).")
```

### ✏️ MUDANÇA 2: Função `run()` completamente refatorizada

**ANTES**:
```go
func (o *assinarOptions) run(command *cobra.Command, _ []string) error {
	if err := ensureValidJSONPaths(o.Entrada, o.Saida); err != nil {
		return err
	}
	if err := ensureValidAlias(o.Alias); err != nil {
		return err
	}
	if err := ensureValidPKCS11Library(o.BibliotecaPKCS11); err != nil {
		return err
	}
	if err := ensureValidPKCS11Slot(o.SlotPKCS11); err != nil {
		return err
	}

	mode, err := modeToJarValue(o.Modo)
	if err != nil {
		return err
	}
	if err := ensurePort(command, mode, o.Porta); err != nil {
		return err
	}

	args := []string{
		"sign",
		"--pathin", o.Entrada,
		"--pathout", o.Saida,
		"--mode", mode,
		"--alias", o.Alias,
	}
	args = appendPortIfNeeded(args, mode, o.Porta)

	if o.BibliotecaPKCS11 != "" {
		args = append(args, "--pkcs11-lib", o.BibliotecaPKCS11)
	}
	if o.SlotPKCS11 != "" {
		args = append(args, "--pkcs11-slot", o.SlotPKCS11)
	}

	result, err := newRunnerConfig(o.runtimeFlags).Run(args)
	if err != nil {
		return wrapRuntimeError(err)
	}

	return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}
```

**DEPOIS**:
```go
func (o *assinarOptions) run(command *cobra.Command, _ []string) error {
	if err := ensureValidJSONPaths(o.Entrada, o.Saida); err != nil {
		return err
	}
	if err := ensureValidAlias(o.Alias); err != nil {
		return err
	}
	if err := ensureValidPKCS11Library(o.BibliotecaPKCS11); err != nil {
		return err
	}
	if err := ensureValidPKCS11Slot(o.SlotPKCS11); err != nil {
		return err
	}

	// ← NOVO: Converte --modo para estratégia
	strategy, err := ParseExecutionStrategy(o.Modo)
	if err != nil {
		return err
	}

	// ← NOVO: Se HTTP foi forçado, valida porta
	if strategy == StrategyHTTP {
		if err := ensurePort(command, "http", o.Porta); err != nil {
			return err
		}
	}

	// ← NOVO: Constrói argumentos para JAR (sem --mode)
	args := []string{
		"sign",
		"--pathin", o.Entrada,
		"--pathout", o.Saida,
		"--alias", o.Alias,
	}

	// ← MANTIDO: Adiciona PKCS#11 se fornecido
	if o.BibliotecaPKCS11 != "" {
		args = append(args, "--pkcs11-lib", o.BibliotecaPKCS11)
	}
	if o.SlotPKCS11 != "" {
		args = append(args, "--pkcs11-slot", o.SlotPKCS11)
	}

	// ← NOVO: Usa RunWithStrategy para seleção automática de modo
	config := newRunnerConfig(o.runtimeFlags)
	result, err := config.RunWithStrategy(args, strategy.String(), o.Porta)
	if err != nil {
		return wrapRuntimeError(err)
	}

	return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}
```

---

## 📄 ARQUIVO 4: `cmd/verificar.go`

### ✏️ MUDANÇA 1: Flag --modo (na função `newVerificarCmd()`)

**ANTES**:
```go
	flags.StringVar(&options.Modo, "modo", "", "Modo de execucao: direto ou http.")
```

**DEPOIS**:
```go
	flags.StringVar(&options.Modo, "modo", "auto",
		"Estrategia de execucao: auto, http ou direto.\n"+
			"  auto   (padrao): Tenta servidor HTTP, fallback para direto se indisponivel.\n"+
			"  http   : Usa apenas servidor HTTP, falha se indisponivel.\n"+
			"  direto : Usa apenas execucao direta (sem servidor).")
```

### ✏️ MUDANÇA 2: Função `run()` refatorizada (sem PKCS#11)

**ANTES**:
```go
func (o *verificarOptions) run(command *cobra.Command, _ []string) error {
	if err := ensureValidJSONPaths(o.Entrada, o.Saida); err != nil {
		return err
	}
	if command.Flags().Changed("alias") || command.Flags().Changed("biblioteca-pkcs11") || command.Flags().Changed("slot-pkcs11") {
		return validationError("O comando verificar nao aceita --alias, --biblioteca-pkcs11 nem --slot-pkcs11.")
	}

	mode, err := modeToJarValue(o.Modo)
	if err != nil {
		return err
	}
	if err := ensurePort(command, mode, o.Porta); err != nil {
		return err
	}

	args := []string{
		"validate",
		"--pathin", o.Entrada,
		"--pathout", o.Saida,
		"--mode", mode,
	}
	args = appendPortIfNeeded(args, mode, o.Porta)

	result, err := newRunnerConfig(o.runtimeFlags).Run(args)
	if err != nil {
		return wrapRuntimeError(err)
	}

	return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}
```

**DEPOIS**:
```go
func (o *verificarOptions) run(command *cobra.Command, _ []string) error {
	if err := ensureValidJSONPaths(o.Entrada, o.Saida); err != nil {
		return err
	}
	if command.Flags().Changed("alias") || command.Flags().Changed("biblioteca-pkcs11") || command.Flags().Changed("slot-pkcs11") {
		return validationError("O comando verificar nao aceita --alias, --biblioteca-pkcs11 nem --slot-pkcs11.")
	}

	// ← NOVO: Converte --modo para estratégia
	strategy, err := ParseExecutionStrategy(o.Modo)
	if err != nil {
		return err
	}

	// ← NOVO: Se HTTP foi forçado, valida porta
	if strategy == StrategyHTTP {
		if err := ensurePort(command, "http", o.Porta); err != nil {
			return err
		}
	}

	// ← NOVO: Constrói argumentos (sem --mode)
	args := []string{
		"validate",
		"--pathin", o.Entrada,
		"--pathout", o.Saida,
	}

	// ← NOVO: Usa RunWithStrategy
	config := newRunnerConfig(o.runtimeFlags)
	result, err := config.RunWithStrategy(args, strategy.String(), o.Porta)
	if err != nil {
		return wrapRuntimeError(err)
	}

	return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}
```

---

## 📄 ARQUIVO 5: `cmd/strategy_test.go` [NOVO]

**CRIADO**: Arquivo completamente novo com 6 testes

```go
package cmd

import (
	"errors"
	"testing"
)

func TestParseExecutionStrategyDefault(t *testing.T) { /* ... */ }
func TestParseExecutionStrategyAuto(t *testing.T) { /* ... */ }
func TestParseExecutionStrategyHTTP(t *testing.T) { /* ... */ }
func TestParseExecutionStrategyDireto(t *testing.T) { /* ... */ }
func TestParseExecutionStrategyInvalid(t *testing.T) { /* ... */ }
func TestParseExecutionStrategyCaseInsensitive(t *testing.T) { /* ... */ }
```

---

## 📄 ARQUIVO 6: `internal/runner/runner_test.go`

### ✏️ MUDANÇA: Imports

**ANTES**:
```go
import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)
```

**DEPOIS**:
```go
import (
	"net"      // ← NOVO
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)
```

### ✏️ MUDANÇA: Adicionar 8 testes ao final do arquivo

**ADICIONADO** (após função `createFakeJava()`):
```go
func TestIsServerAvailableWhenListening(t *testing.T) { /* ... */ }
func TestIsServerAvailableWhenNotListening(t *testing.T) { /* ... */ }
func TestDetermineExecutionModeDirect(t *testing.T) { /* ... */ }
func TestDetermineExecutionModeHTTP(t *testing.T) { /* ... */ }
func TestDetermineExecutionModeAutoWithServer(t *testing.T) { /* ... */ }
func TestDetermineExecutionModeAutoWithoutServer(t *testing.T) { /* ... */ }
func TestApplyExecutionModeAddsMode(t *testing.T) { /* ... */ }
func TestApplyExecutionModeRemovesExistingMode(t *testing.T) { /* ... */ }
func TestRemoveModeFromArgs(t *testing.T) { /* ... */ }
```

---

## 📊 RESUMO DE MUDANÇAS

| Arquivo | Tipo | Mudanças | Linhas |
|:---|:---:|:---|---:|
| common.go | ✏️ Modificado | Imports + tipo + 1 função | +80 |
| runner.go | ✏️ Modificado | Imports + 5 novas funções | +120 |
| assinar.go | ✏️ Modificado | Flag + função run() | ±40 |
| verificar.go | ✏️ Modificado | Flag + função run() | ±40 |
| strategy_test.go | ✨ Novo | 6 testes novos | +80 |
| runner_test.go | ✏️ Modificado | Imports + 8 testes | +80 |
| **TOTAL** | | | **~450** |

---

## ✨ MUDANÇAS POR TIPO

### Imports Adicionados
- ✅ `import "net"`
- ✅ `import "time"`

### Tipos Novos
- ✅ `type ExecutionStrategy string`
- ✅ 3 constantes: StrategyAuto, StrategyHTTP, StrategyDirect

### Funções Novas
- ✅ `ParseExecutionStrategy()` — em common.go
- ✅ `IsServerAvailable()` — em runner.go
- ✅ `DetermineExecutionMode()` — em runner.go
- ✅ `ApplyExecutionMode()` — em runner.go
- ✅ `RemoveModeFromArgs()` — em runner.go
- ✅ `RunWithStrategy()` — em runner.go (método)

### Funções Refatoradas
- ✅ `assinarOptions.run()` — remove manual `--mode`, chama `RunWithStrategy()`
- ✅ `verificarOptions.run()` — remove manual `--mode`, chama `RunWithStrategy()`

### Testes Novos
- ✅ 6 testes em strategy_test.go
- ✅ 8 testes em runner_test.go
- ✅ Total: 14 testes novos

---

## 🔄 FLUXO DE MUDANÇA

```
Antes:
user → assinar.go → run() → modeToJarValue() → Run()

Depois:
user → assinar.go → run() → ParseExecutionStrategy() 
     → RunWithStrategy() 
     → DetermineExecutionMode() 
     → IsServerAvailable()
     → ApplyExecutionMode() 
     → Run()
```

---

**Referência completa de mudanças**  
Use este documento para revisar cada mudança em detalhe.

