# Quick Start — Implementação US-01 em 20 Passos
## Guia Passo-a-Passo para Desenvolvedor

---

## PREPARAÇÃO (5 minutos)

### Passo 1: Clonar e preparar branch
```bash
cd c:\Users\Aluno\Downloads\apagar\runner
git checkout -b feature/us01-auto-mode
```

### Passo 2: Revisar estrutura
```bash
cd assinador-cli
ls -la cmd/      # Verificar arquivos: assinar.go, verificar.go, comum.go, servidor.go, root.go
ls -la internal/runner/  # Verificar: runner.go, runner_test.go
```

### Passo 3: Revisar testes existentes
```bash
go test ./... -v  # Executar testes base (devem passar)
```

---

## FASE 1: ADICIONAR EXECUTIONSTRATEGY (30 minutos)

### Passo 4: Adicionar tipo em `cmd/common.go`

**Abrir**: `assinador-cli/cmd/common.go`

**Localizar**: Encontrar `type runtimeFlags struct { ... }`

**Adicionar após existentes (antes das funções)**:
```go
// ExecutionStrategy defines how assinador.jar should be invoked
type ExecutionStrategy string

const (
    StrategyAuto   ExecutionStrategy = "auto"
    StrategyHTTP   ExecutionStrategy = "http"
    StrategyDirect ExecutionStrategy = "direct"
)

func (s ExecutionStrategy) String() string {
    return string(s)
}

// ParseExecutionStrategy converts user input to strategy
func ParseExecutionStrategy(input string) (ExecutionStrategy, error) {
    normalized := strings.ToLower(strings.TrimSpace(input))
    
    if normalized == "" {
        return StrategyAuto, nil
    }
    
    switch normalized {
    case "auto":
        return StrategyAuto, nil
    case "http":
        return StrategyHTTP, nil
    case "direto", "direct":
        return StrategyDirect, nil
    default:
        return "", validationError(
            "Estrategia de modo invalida: %s. Use 'auto', 'http' ou 'direto'.",
            input,
        )
    }
}
```

### Passo 5: Adicionar helpers em `cmd/common.go`

**Adicionar ao final da seção de tipos/constantes**:
```go
// RemoveModeFromArgs removes any existing --mode flag and value
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
```

### Passo 6: Atualizar flag em `cmd/assinar.go`

**Localizar**: Função `newAssinarCmd()`, linha com `flags.StringVar(&options.Modo...`

**Substituir**:
```go
// OLD:
// flags.StringVar(&options.Modo, "modo", "", "Modo de execucao: direto ou http.")

// NEW:
flags.StringVar(&options.Modo, "modo", "auto", 
    "Estrategia de execucao: auto, http ou direto.\n"+
    "  auto   (padrao): Tenta servidor HTTP, fallback para direto.\n"+
    "  http   : Usa apenas servidor HTTP.\n"+
    "  direto : Usa apenas execucao direta.")
```

### Passo 7: Atualizar flag em `cmd/verificar.go`

**Localizar**: Função `newVerificarCmd()`, linha com `flags.StringVar(&options.Modo...`

**Substituir**: Mesmo que passo anterior

### Passo 8: Testar parsing
```bash
cd assinador-cli
go test -run TestParseExecutionStrategy ./cmd -v
# Esperado: 0 testes (TestParseExecutionStrategy não existe, será criado depois)
go build .  # Deve compilar sem erros
```

---

## FASE 2: ADICIONAR DETECÇÃO DE SERVIDOR (40 minutos)

### Passo 9: Adicionar imports em `internal/runner/runner.go`

**Localizar**: Section de imports no topo do arquivo

**Adicionar**:
```go
import (
    // ... existing imports ...
    "net"
    "time"
)
```

### Passo 10: Adicionar funções de detecção em `internal/runner/runner.go`

**Adicionar antes de `func (c Config) Run(...)`**:
```go
// IsServerAvailable checks if the server is responding on given port
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

// DetermineExecutionMode decides final execution mode based on strategy
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

// ApplyExecutionMode modifies args to set --mode and --port
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
```

### Passo 11: Adicionar RunWithStrategy em `internal/runner/runner.go`

**Adicionar após `func (c Config) StartServer(...)`**:
```go
// RunWithStrategy executes with automatic mode selection
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

### Passo 12: Testar compilação
```bash
go build .
# Esperado: Sem erros
```

---

## FASE 3: INTEGRAR EM ASSINAR (20 minutos)

### Passo 13: Atualizar `cmd/assinar.go` função `run()`

**Localizar**: Função `(o *assinarOptions) run(command *cobra.Command, _ []string) error`

**Encontrar linha**: `mode, err := modeToJarValue(o.Modo)`

**Substituir TODO O BLOCO** que começa com `mode, err := modeToJarValue`:
```go
// REMOVER:
// mode, err := modeToJarValue(o.Modo)
// if err != nil {
//     return err
// }
// if err := ensurePort(command, mode, o.Porta); err != nil {
//     return err
// }
// 
// args := []string{...}
// args = appendPortIfNeeded(args, mode, o.Porta)
// 
// result, err := newRunnerConfig(o.runtimeFlags).Run(args)

// ADICIONAR:
strategy, err := ParseExecutionStrategy(o.Modo)
if err != nil {
    return err
}

if strategy == StrategyHTTP {
    if err := ensurePort(command, "http", o.Porta); err != nil {
        return err
    }
}

args := []string{
    "sign",
    "--pathin", o.Entrada,
    "--pathout", o.Saida,
    "--alias", o.Alias,
}
args = appendPortIfNeeded(args, "", o.Porta)  // Não adiciona --mode aqui

if o.BibliotecaPKCS11 != "" {
    args = append(args, "--pkcs11-lib", o.BibliotecaPKCS11)
}
if o.SlotPKCS11 != "" {
    args = append(args, "--pkcs11-slot", o.SlotPKCS11)
}

config := newRunnerConfig(o.runtimeFlags)
result, err := config.RunWithStrategy(args, strategy.String(), o.Porta)
```

### Passo 14: Atualizar `cmd/verificar.go` função `run()`

**Localizar**: Função `(o *verificarOptions) run(...)`

**Encontrar linha**: `mode, err := modeToJarValue(o.Modo)`

**Substituir TODO O BLOCO** (similar a passo 13):
```go
strategy, err := ParseExecutionStrategy(o.Modo)
if err != nil {
    return err
}

if strategy == StrategyHTTP {
    if err := ensurePort(command, "http", o.Porta); err != nil {
        return err
    }
}

args := []string{
    "validate",
    "--pathin", o.Entrada,
    "--pathout", o.Saida,
}
args = appendPortIfNeeded(args, "", o.Porta)

config := newRunnerConfig(o.runtimeFlags)
result, err := config.RunWithStrategy(args, strategy.String(), o.Porta)
```

### Passo 15: Testar compilação
```bash
go build .
# Esperado: Sem erros
```

---

## FASE 4: CRIAR TESTES (30 minutos)

### Passo 16: Criar `cmd/strategy_test.go`

**Novo arquivo**: `assinador-cli/cmd/strategy_test.go`

**Conteúdo**:
```go
package cmd

import (
    "errors"
    "testing"
)

func TestParseExecutionStrategyDefault(t *testing.T) {
    strategy, err := ParseExecutionStrategy("")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if strategy != StrategyAuto {
        t.Fatalf("expected StrategyAuto, got %v", strategy)
    }
}

func TestParseExecutionStrategyAuto(t *testing.T) {
    strategy, err := ParseExecutionStrategy("auto")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if strategy != StrategyAuto {
        t.Fatalf("expected StrategyAuto, got %v", strategy)
    }
}

func TestParseExecutionStrategyHTTP(t *testing.T) {
    strategy, err := ParseExecutionStrategy("http")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if strategy != StrategyHTTP {
        t.Fatalf("expected StrategyHTTP, got %v", strategy)
    }
}

func TestParseExecutionStrategyDireto(t *testing.T) {
    strategy, err := ParseExecutionStrategy("direto")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if strategy != StrategyDirect {
        t.Fatalf("expected StrategyDirect, got %v", strategy)
    }
}

func TestParseExecutionStrategyInvalid(t *testing.T) {
    _, err := ParseExecutionStrategy("invalid")
    if err == nil {
        t.Fatalf("expected error for invalid strategy")
    }
    var exitErr *ExitError
    if !errors.As(err, &exitErr) || exitErr.Code != 2 {
        t.Fatalf("expected validation error")
    }
}
```

### Passo 17: Adicionar testes a `internal/runner/runner_test.go`

**Adicionar ao final**:
```go
func TestIsServerAvailableWhenListening(t *testing.T) {
    listener, _ := net.Listen("tcp", "localhost:0")
    defer listener.Close()
    port := listener.Addr().(*net.TCPAddr).Port
    
    if !IsServerAvailable(port, 1) {
        t.Fatalf("expected server to be available")
    }
}

func TestIsServerAvailableWhenNotListening(t *testing.T) {
    if IsServerAvailable(59999, 1) {
        t.Fatalf("expected server to be unavailable")
    }
}

func TestDetermineExecutionModeDirect(t *testing.T) {
    mode, _ := DetermineExecutionMode("direct", 8080)
    if mode != "direct" {
        t.Fatalf("expected direct")
    }
}

func TestDetermineExecutionModeHTTP(t *testing.T) {
    mode, _ := DetermineExecutionMode("http", 8080)
    if mode != "http" {
        t.Fatalf("expected http")
    }
}

func TestDetermineExecutionModeAutoWithServer(t *testing.T) {
    listener, _ := net.Listen("tcp", "localhost:0")
    defer listener.Close()
    port := listener.Addr().(*net.TCPAddr).Port
    
    mode, _ := DetermineExecutionMode("auto", port)
    if mode != "http" {
        t.Fatalf("expected http when server available")
    }
}

func TestDetermineExecutionModeAutoWithoutServer(t *testing.T) {
    mode, _ := DetermineExecutionMode("auto", 59999)
    if mode != "direct" {
        t.Fatalf("expected direct when server unavailable")
    }
}

func TestApplyExecutionModeAddsMode(t *testing.T) {
    args := []string{"sign", "--pathin", "e.json"}
    result := ApplyExecutionMode(args, "http", 8080)
    
    hasMode := false
    for i := 0; i < len(result)-1; i++ {
        if result[i] == "--mode" && result[i+1] == "http" {
            hasMode = true
        }
    }
    if !hasMode {
        t.Fatalf("--mode not added: %v", result)
    }
}
```

### Passo 18: Executar todos os testes
```bash
cd assinador-cli
go test ./...
# Esperado: Todos os testes passam (antigos + novos)
```

### Passo 19: Verificar cobertura
```bash
go test -cover ./...
# Esperado: common.go e runner.go com ≥80%
```

---

## FASE 5: VALIDAÇÃO MANUAL (20 minutos)

### Passo 20: Testar 5 cenários

**Teste 1: Modo automático (novo padrão)**
```bash
assinador-cli assinar --help | grep -A2 "modo"
# Esperado: mostra "auto" como padrão

assinador-cli assinar --entrada entrada.json --saida s.json --alias demo
# Esperado: Não reclamação de modo obrigatório
```

**Teste 2: Forçar modo HTTP**
```bash
assinador-cli assinar --entrada entrada.json --saida s.json --alias demo --modo http
# Esperado: Tenta HTTP
```

**Teste 3: Forçar modo direto**
```bash
assinador-cli assinar --entrada entrada.json --saida s.json --alias demo --modo direto
# Esperado: Usa direto
```

**Teste 4: Verificar compatibilidade reversa**
```bash
assinador-cli assinar --entrada entrada.json --saida s.json --alias demo --modo direto
# Esperado: Mesma saída que antes
```

**Teste 5: Servidor com timeout**
```bash
assinador-cli servidor iniciar --porta 8080 --timeout 30
# Esperado: Inicia normalmente
```

---

## CHECKLIST FINAL

- [ ] Passo 4-7: ExecutionStrategy adicionado e compilação OK
- [ ] Passo 8: Build sem erros
- [ ] Passo 9-12: Detecção de servidor adicionada e compilação OK
- [ ] Passo 13-15: Integração em assinar/verificar e compilação OK
- [ ] Passo 16-17: Testes criados
- [ ] Passo 18: `go test ./...` passa 100%
- [ ] Passo 19: Cobertura ≥80%
- [ ] Passo 20: 5 testes manuais validam

---

## TROUBLESHOOTING

| Erro | Solução |
|:---|:---|
| `undefined: StrategyAuto` | Verificar que type `ExecutionStrategy` foi adicionado em `common.go` |
| `undefined: ParseExecutionStrategy` | Verificar que função foi adicionada em `common.go` |
| `undefined: IsServerAvailable` | Verificar que função foi adicionada em `internal/runner/runner.go` |
| `missing "net" import` | Adicionar `"net"` aos imports em `internal/runner/runner.go` |
| Teste `TestParseExecutionStrategyDefault` falha | Verificar que default de `--modo` é `"auto"` não `""` |
| Build falha | Executar `go mod tidy` e `go get -u` |

---

## PRÓXIMO: CODE REVIEW

1. Commit mudanças: `git commit -am "US-01: Auto-detection e fallback de modo"`
2. Push branch: `git push origin feature/us01-auto-mode`
3. Abrir Pull Request no GitHub
4. Solicitar review
5. Executar testes em CI/CD
6. Merge em main após aprovação

---

**Tempo esperado**: 2-3 horas  
**Complexidade**: Média  
**Risco**: Baixo (testes cobrem regressões)  
**Status**: Pronto para implementação
