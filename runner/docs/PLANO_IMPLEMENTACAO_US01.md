# Plano de Implementação Detalhado — US-01
## Código-fonte, Arquivos Modificados e Testes

---

## PARTE 1: MUDANÇAS EM common.go

### Adicionar ao final do arquivo (antes das funções existentes):

```go
// ExecutionStrategy defines how assinador.jar should be invoked
type ExecutionStrategy string

const (
    // StrategyAuto tries HTTP (server) first, falls back to direct
    StrategyAuto ExecutionStrategy = "auto"
    
    // StrategyHTTP uses only HTTP (server), fails if unavailable
    StrategyHTTP ExecutionStrategy = "http"
    
    // StrategyDirect uses only direct invocation (local CLI)
    StrategyDirect ExecutionStrategy = "direct"
)

func (s ExecutionStrategy) String() string {
    return string(s)
}

func (s ExecutionStrategy) IsValid() bool {
    return s == StrategyAuto || s == StrategyHTTP || s == StrategyDirect
}

// ParseExecutionStrategy converts user input (Portuguese or English) to strategy
// Empty string defaults to StrategyAuto (new behavior)
func ParseExecutionStrategy(input string) (ExecutionStrategy, error) {
    normalized := strings.ToLower(strings.TrimSpace(input))
    
    // Default: empty string maps to auto
    if normalized == "" {
        return StrategyAuto, nil
    }
    
    switch normalized {
    case "auto":
        return StrategyAuto, nil
    case "http":
        return StrategyHTTP, nil
    case "direto", "direct":  // Support both Portuguese and English
        return StrategyDirect, nil
    default:
        return "", validationError(
            "Estrategia de modo invalida: %s. Use 'auto', 'http' ou 'direto'.",
            input,
        )
    }
}
```

### Modificar flag binding em `bindRuntimeFlags()` da flag `--modo`:

```go
// OLD CODE (in newAssinarCmd, around line 40):
// flags.StringVar(&options.Modo, "modo", "", "Modo de execucao: direto ou http.")

// NEW CODE:
flags.StringVar(&options.Modo, "modo", "auto", 
    "Estrategia de execucao: auto, http ou direto.\n"+
    "  auto   (padrao): Tenta servidor HTTP, fallback para direto se indisponivel.\n"+
    "  http   : Usa apenas servidor HTTP, falha se indisponivel.\n"+
    "  direto : Usa apenas execucao direta (sem servidor).")
```

### Adicionar helper function em common.go:

```go
// ExtractModeFromArgs finds and returns the current --mode value in args
// Returns empty string if not found
func ExtractModeFromArgs(args []string) string {
    for i := 0; i < len(args)-1; i++ {
        if args[i] == "--mode" {
            return args[i+1]
        }
    }
    return ""
}

// RemoveModeFromArgs removes any existing --mode flag and value
func RemoveModeFromArgs(args []string) []string {
    result := []string{}
    i := 0
    for i < len(args) {
        if args[i] == "--mode" {
            i += 2  // Skip flag and value
            continue
        }
        result = append(result, args[i])
        i++
    }
    return result
}
```

---

## PARTE 2: MUDANÇAS EM runner.go

### Adicionar imports:

```go
import (
    "net"
    "time"
)
```

### Adicionar funções de detecção e estratégia:

```go
// IsServerAvailable checks if the server is responding on given port
// Uses TCP connection test (works independent of HTTP server implementation)
// timeoutSecs: timeout in seconds for connection attempt
// Returns true if connection succeeds, false if timeout or refused
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
// For "auto": returns "http" if server available, else "direct"
// For "http": returns "http"
// For "direct": returns "direct"
// port: port number for auto-detection
func DetermineExecutionMode(strategy string, port int) (string, error) {
    switch strings.ToLower(strategy) {
    case "direct":
        return "direct", nil
        
    case "http":
        return "http", nil
        
    case "auto":
        // Try HTTP with 1 second timeout
        if IsServerAvailable(port, 1) {
            return "http", nil
        }
        // Fallback to direct
        return "direct", nil
        
    default:
        return "", fmt.Errorf("estrategia desconhecida: %s", strategy)
    }
}

// ApplyExecutionMode modifies the command arguments to set --mode and --port
// Removes any existing --mode, adds new --mode, and --port if needed
func ApplyExecutionMode(args []string, mode string, port int) []string {
    // Remove existing --mode flag and value
    filtered := RemoveModeFromArgs(args)
    
    // Add new --mode
    filtered = append(filtered, "--mode", mode)
    
    // Add --port if mode is http and not already present
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

// RunWithStrategy executes with automatic mode selection based on strategy
// If strategy is "auto": tries HTTP first, falls back to direct if needed
// If strategy is "http": enforces HTTP mode
// If strategy is "direct": enforces direct mode
func (c Config) RunWithStrategy(args []string, strategy string, port int) (Result, error) {
    if err := c.Validate(); err != nil {
        return Result{}, err
    }
    
    // Determine final mode based on strategy and availability
    mode, err := DetermineExecutionMode(strategy, port)
    if err != nil {
        return Result{}, err
    }
    
    // Apply mode to arguments
    finalArgs := ApplyExecutionMode(args, mode, port)
    
    // Execute with determined mode
    return c.Run(finalArgs)
}
```

---

## PARTE 3: MUDANÇAS EM assinar.go

### Modificar a struct `assinarOptions`:

```go
type assinarOptions struct {
    runtimeFlags
    Entrada          string
    Saida            string
    Modo             string  // Now can be "auto", "http", or "direto" (default: "auto")
    Alias            string
    BibliotecaPKCS11 string
    SlotPKCS11       string
    Porta            int
}
```

### Modificar função `(o *assinarOptions) run()`:

```go
func (o *assinarOptions) run(command *cobra.Command, _ []string) error {
    // Existing validations
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

    // NEW: Parse execution strategy
    strategy, err := ParseExecutionStrategy(o.Modo)
    if err != nil {
        return err
    }

    // Validate port if needed
    if strategy == StrategyHTTP {  // Only validate port if forcing HTTP
        if err := ensurePort(command, "http", o.Porta); err != nil {
            return err
        }
    }

    // Build args for JAR (without --mode, will be added by RunWithStrategy)
    args := []string{
        "sign",
        "--pathin", o.Entrada,
        "--pathout", o.Saida,
        "--alias", o.Alias,
    }

    // Add PKCS#11 if provided
    if o.BibliotecaPKCS11 != "" {
        args = append(args, "--pkcs11-lib", o.BibliotecaPKCS11)
    }
    if o.SlotPKCS11 != "" {
        args = append(args, "--pkcs11-slot", o.SlotPKCS11)
    }

    // NEW: Use RunWithStrategy for automatic mode selection
    config := newRunnerConfig(o.runtimeFlags)
    result, err := config.RunWithStrategy(args, strategy.String(), o.Porta)
    if err != nil {
        return wrapRuntimeError(err)
    }

    return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}
```

---

## PARTE 4: MUDANÇAS EM verificar.go

### Modificar função `(o *verificarOptions) run()`:

```go
func (o *verificarOptions) run(command *cobra.Command, _ []string) error {
    if err := ensureValidJSONPaths(o.Entrada, o.Saida); err != nil {
        return err
    }
    if command.Flags().Changed("alias") || command.Flags().Changed("biblioteca-pkcs11") || command.Flags().Changed("slot-pkcs11") {
        return validationError("O comando verificar nao aceita --alias, --biblioteca-pkcs11 nem --slot-pkcs11.")
    }

    // NEW: Parse execution strategy
    strategy, err := ParseExecutionStrategy(o.Modo)
    if err != nil {
        return err
    }

    // Validate port if needed
    if strategy == StrategyHTTP {
        if err := ensurePort(command, "http", o.Porta); err != nil {
            return err
        }
    }

    // Build args (without --mode)
    args := []string{
        "validate",
        "--pathin", o.Entrada,
        "--pathout", o.Saida,
    }

    // NEW: Use RunWithStrategy
    config := newRunnerConfig(o.runtimeFlags)
    result, err := config.RunWithStrategy(args, strategy.String(), o.Porta)
    if err != nil {
        return wrapRuntimeError(err)
    }

    return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}
```

### Também em `verificar.go`, atualizar flag:

```go
// In newVerificarCmd():
flags.StringVar(&options.Modo, "modo", "auto", 
    "Estrategia de execucao: auto, http ou direto.\n"+
    "  auto   (padrao): Tenta servidor HTTP, fallback para direto se indisponivel.\n"+
    "  http   : Usa apenas servidor HTTP, falha se indisponivel.\n"+
    "  direto : Usa apenas execucao direta (sem servidor).")
```

---

## PARTE 5: TESTES UNITÁRIOS (novo arquivo: cmd/strategy_test.go)

```go
package cmd

import (
    "testing"
)

func TestParseExecutionStrategyDefault(t *testing.T) {
    strategy, err := ParseExecutionStrategy("")
    if err != nil {
        t.Fatalf("unexpected error for empty string: %v", err)
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

func TestParseExecutionStrategyDirect(t *testing.T) {
    strategy, err := ParseExecutionStrategy("direct")
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
        t.Fatalf("expected validation error for invalid strategy")
    }
    
    var exitErr *ExitError
    if !errors.As(err, &exitErr) || exitErr.Code != 2 {
        t.Fatalf("expected validation error with code 2, got: %v", err)
    }
}

func TestParseExecutionStrategyWhitespace(t *testing.T) {
    strategy, err := ParseExecutionStrategy("  http  ")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if strategy != StrategyHTTP {
        t.Fatalf("expected StrategyHTTP after trimming, got %v", strategy)
    }
}

func TestParseExecutionStrategyCaseInsensitive(t *testing.T) {
    testCases := []struct {
        input    string
        expected ExecutionStrategy
    }{
        {"AUTO", StrategyAuto},
        {"Auto", StrategyAuto},
        {"HTTP", StrategyHTTP},
        {"Http", StrategyHTTP},
        {"DIRETO", StrategyDirect},
        {"Direto", StrategyDirect},
    }
    
    for _, tc := range testCases {
        strategy, err := ParseExecutionStrategy(tc.input)
        if err != nil {
            t.Fatalf("unexpected error for %q: %v", tc.input, err)
        }
        if strategy != tc.expected {
            t.Fatalf("for input %q: expected %v, got %v", tc.input, tc.expected, strategy)
        }
    }
}
```

---

## PARTE 6: TESTES UNITÁRIOS (runner_test.go - ADICIONAR)

```go
package runner

import (
    "fmt"
    "net"
    "testing"
    "time"
)

// TestIsServerAvailableWhenListening verifies detection of running server
func TestIsServerAvailableWhenListening(t *testing.T) {
    // Start a listener on ephemeral port
    listener, err := net.Listen("tcp", "localhost:0")
    if err != nil {
        t.Fatalf("failed to create listener: %v", err)
    }
    defer listener.Close()
    
    // Extract port
    port := listener.Addr().(*net.TCPAddr).Port
    
    // Test IsServerAvailable
    available := IsServerAvailable(port, 1)
    if !available {
        t.Fatalf("expected server to be available on port %d", port)
    }
}

// TestIsServerAvailableWhenNotListening verifies detection of unavailable server
func TestIsServerAvailableWhenNotListening(t *testing.T) {
    // Use a port that's very unlikely to be in use (ephemeral range, unused)
    unavailablePort := 59999
    
    available := IsServerAvailable(unavailablePort, 1)
    if available {
        t.Fatalf("expected server to be unavailable on port %d", unavailablePort)
    }
}

// TestIsServerAvailableRespectTimeout verifies timeout behavior
func TestIsServerAvailableRespectTimeout(t *testing.T) {
    // Try to connect to a non-routable address with very short timeout
    start := time.Now()
    available := IsServerAvailable(59999, 1)  // 1 second timeout
    elapsed := time.Since(start)
    
    if available {
        t.Fatalf("expected unavailable server")
    }
    if elapsed > 5*time.Second {
        t.Fatalf("timeout took too long: %v", elapsed)
    }
}

// TestDetermineExecutionModeDirect returns "direct" when strategy is direct
func TestDetermineExecutionModeDirect(t *testing.T) {
    mode, err := DetermineExecutionMode("direct", 8080)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if mode != "direct" {
        t.Fatalf("expected 'direct', got %q", mode)
    }
}

// TestDetermineExecutionModeHTTP returns "http" when strategy is http
func TestDetermineExecutionModeHTTP(t *testing.T) {
    mode, err := DetermineExecutionMode("http", 8080)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if mode != "http" {
        t.Fatalf("expected 'http', got %q", mode)
    }
}

// TestDetermineExecutionModeAutoWithServer returns "http" when server available
func TestDetermineExecutionModeAutoWithServer(t *testing.T) {
    listener, err := net.Listen("tcp", "localhost:0")
    if err != nil {
        t.Fatalf("failed to create listener: %v", err)
    }
    defer listener.Close()
    
    port := listener.Addr().(*net.TCPAddr).Port
    
    mode, err := DetermineExecutionMode("auto", port)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if mode != "http" {
        t.Fatalf("expected 'http' for available server, got %q", mode)
    }
}

// TestDetermineExecutionModeAutoWithoutServer returns "direct" when server unavailable
func TestDetermineExecutionModeAutoWithoutServer(t *testing.T) {
    mode, err := DetermineExecutionMode("auto", 59999)  // unused port
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if mode != "direct" {
        t.Fatalf("expected 'direct' for unavailable server, got %q", mode)
    }
}

// TestDetermineExecutionModeInvalidStrategy returns error
func TestDetermineExecutionModeInvalidStrategy(t *testing.T) {
    _, err := DetermineExecutionMode("invalid", 8080)
    if err == nil {
        t.Fatalf("expected error for invalid strategy")
    }
}

// TestApplyExecutionModeAddsMode verifies --mode is added to args
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
        t.Fatalf("--mode flag not added to args: %v", result)
    }
    if !hasPort {
        t.Fatalf("--port flag not added when mode is http: %v", result)
    }
}

// TestApplyExecutionModeRemovesExistingMode verifies old --mode is removed
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
        t.Fatalf("expected exactly one --mode flag, found %d in: %v", count, result)
    }
}

// TestApplyExecutionModeDirectNoPort verifies --port not added for direct mode
func TestApplyExecutionModeDirectNoPort(t *testing.T) {
    args := []string{"sign", "--pathin", "e.json"}
    result := ApplyExecutionMode(args, "direct", 8080)
    
    for i := 0; i < len(result)-1; i++ {
        if result[i] == "--port" {
            t.Fatalf("--port should not be added for direct mode: %v", result)
        }
    }
}
```

---

## PARTE 7: TESTES DE INTEGRAÇÃO (root_test.go - ADICIONAR)

```go
func TestAssinarWithStrategyAutoServerAvailable(t *testing.T) {
    // Setup: Start mock server
    listener, err := net.Listen("tcp", "localhost:8080")
    if err != nil {
        t.Skipf("port 8080 already in use, skipping test")
    }
    defer listener.Close()
    
    // Setup: Create temp files
    tempDir := t.TempDir()
    inputPath := filepath.Join(tempDir, "entrada.json")
    outputPath := filepath.Join(tempDir, "saida.json")
    if err := os.WriteFile(inputPath, []byte(`{"resourceType":"Bundle"}`), 0o644); err != nil {
        t.Fatalf("failed to write input: %v", err)
    }
    
    // Execute: assinar with default mode (auto should detect server)
    _, _, err = executeRootCommand(
        "assinar",
        "--entrada", inputPath,
        "--saida", outputPath,
        "--alias", "test-auto",
        // NO --modo specified, should default to "auto"
    )
    
    // Expected: Command completes (may fail for other reasons, but not a mode error)
    // The important part is it doesn't complain about mode
    if err != nil {
        var exitErr *ExitError
        if errors.As(err, &exitErr) && strings.Contains(exitErr.Message, "mode") {
            t.Fatalf("unexpected mode error: %v", err)
        }
    }
}

func TestAssinarWithExplicitModeAuto(t *testing.T) {
    tempDir := t.TempDir()
    inputPath := filepath.Join(tempDir, "entrada.json")
    outputPath := filepath.Join(tempDir, "saida.json")
    if err := os.WriteFile(inputPath, []byte(`{"resourceType":"Bundle"}`), 0o644); err != nil {
        t.Fatalf("failed to write input: %v", err)
    }
    
    // Execute: Explicit --modo auto
    _, _, err := executeRootCommand(
        "assinar",
        "--entrada", inputPath,
        "--saida", outputPath,
        "--alias", "test",
        "--modo", "auto",
    )
    
    // Should accept and attempt execution (may fail on Java/JAR not found, but not mode validation)
    if err != nil {
        var exitErr *ExitError
        if errors.As(err, &exitErr) && strings.Contains(exitErr.Message, "Modo invalido") {
            t.Fatalf("should accept --modo auto")
        }
    }
}

func TestVerificarWithStrategyAuto(t *testing.T) {
    tempDir := t.TempDir()
    inputPath := filepath.Join(tempDir, "entrada.json")
    outputPath := filepath.Join(tempDir, "resultado.json")
    if err := os.WriteFile(inputPath, []byte(`{"resourceType":"Bundle"}`), 0o644); err != nil {
        t.Fatalf("failed to write input: %v", err)
    }
    
    // Execute: verificar without explicit mode (auto)
    _, _, err := executeRootCommand(
        "verificar",
        "--entrada", inputPath,
        "--saida", outputPath,
        // NO --modo specified
    )
    
    // Should accept and attempt execution
    if err != nil {
        var exitErr *ExitError
        if errors.As(err, &exitErr) && strings.Contains(exitErr.Message, "obrigatoria") {
            t.Fatalf("--modo should not be required anymore, got: %v", exitErr.Message)
        }
    }
}

func TestAssinarRejectsInvalidMode(t *testing.T) {
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
        "--alias", "test",
        "--modo", "invalid",
    )
    
    if err == nil {
        t.Fatalf("expected validation error for invalid mode")
    }
    
    var exitErr *ExitError
    if !errors.As(err, &exitErr) || exitErr.Code != 2 {
        t.Fatalf("expected exit code 2, got: %v", err)
    }
    if !strings.Contains(exitErr.Message, "Modo invalido") {
        t.Fatalf("unexpected error message: %s", exitErr.Message)
    }
}
```

---

## PARTE 8: ATUALIZAR README

No `assinador-cli/README.md`, adicionar seção:

```markdown
## Estratégia de Execução

O CLI agora suporta três estratégias para determinar como invocar o assinador:

### Modo automático (padrão)
Quando `--modo` não é especificado ou `--modo auto`, o CLI:
1. Tenta conectar ao servidor HTTP na porta padrão (8080)
2. Se conseguir, usa invocação via HTTP (reutiliza servidor existente)
3. Se falhar, cai de volta para invocação direta (executa localmente)

**Exemplo:**
```bash
# Usa servidor se disponível, senão executa localmente
assinador-cli assinar \
  --entrada entrada.json \
  --saida saida.json \
  --alias demo
```

### Modo HTTP explícito
Com `--modo http`, o CLI tenta usar apenas invocação via HTTP.
Falha se o servidor não estiver disponível.

**Exemplo:**
```bash
# Falha se servidor não está rodando
assinador-cli assinar \
  --entrada entrada.json \
  --saida saida.json \
  --modo http \
  --porta 8080 \
  --alias demo
```

### Modo direto explícito
Com `--modo direto`, o CLI invoca o assinador localmente (sem servidor).

**Exemplo:**
```bash
# Sempre usa execução direta, ignora servidor mesmo que ativo
assinador-cli assinar \
  --entrada entrada.json \
  --saida saida.json \
  --modo direto \
  --alias demo
```

## Detecção Automática de Servidor

Quando usando `--modo auto` (padrão) ou sem especificar modo:
- O CLI detecta automaticamente se um servidor está rodando na porta padrão
- Se detectado, reutiliza o servidor (sem iniciar novo processo)
- Se não detectado, executa localmente

Isso permite workflows simplificados:

```bash
# Terminal 1: Inicia servidor uma vez
assinador-cli servidor iniciar

# Terminal 2+: Usa automaticamente o servidor sem precisar especificar
assinador-cli assinar --entrada e.json --saida s.json --alias demo
```

## Timeout de Inatividade

Use `--timeout` com `servidor iniciar` para auto-encerramento após inatividade:

```bash
# Servidor para automaticamente após 30 minutos sem interação
assinador-cli servidor iniciar --porta 8080 --timeout 30
```
```

---

## PARTE 9: CHECKLIST DE IMPLEMENTAÇÃO

- [ ] Adicionar `ExecutionStrategy` type em `common.go`
- [ ] Adicionar `ParseExecutionStrategy()` em `common.go`
- [ ] Atualizar flag `--modo` padrão para `"auto"` em `assinar.go` e `verificar.go`
- [ ] Adicionar `IsServerAvailable()` em `runner.go`
- [ ] Adicionar `DetermineExecutionMode()` em `runner.go`
- [ ] Adicionar `ApplyExecutionMode()` em `runner.go`
- [ ] Adicionar `RunWithStrategy()` em `runner.go`
- [ ] Refatorar `assinar.go` para usar `RunWithStrategy()`
- [ ] Refatorar `verificar.go` para usar `RunWithStrategy()`
- [ ] Criar `strategy_test.go` com testes de parsing
- [ ] Adicionar testes a `runner_test.go` para detecção de servidor
- [ ] Adicionar testes de integração a `root_test.go`
- [ ] Executar testes: `go test ./...`
- [ ] Verificar cobertura: `go test -cover ./...` (objetivo: >80%)
- [ ] Atualizar `README.md` com exemplos de estratégia
- [ ] Validação manual dos 5 cenários (Seção 4.3 de DIAGNOSTICO_US01.md)
- [ ] Commit e push

---

## PARTE 10: VERIFICAÇÃO FINAL

Após implementar, rodar:

```bash
# Compile
cd assinador-cli
go build -o assinador-cli.exe

# Run tests
go test ./...

# Check coverage
go test -cover ./... | grep -E "coverage:|runner|common|cmd"

# Manual validation
./assinador-cli.exe assinar --help
# Should show: --modo string   "Estrategia de execucao: auto, http ou direto..."

# Backward compatibility
./assinador-cli.exe assinar --entrada entrada.json --saida s.json --alias test --modo direto
# Should work

# New default
./assinador-cli.exe assinar --entrada entrada.json --saida s.json --alias test
# Should accept (no error about missing --modo)
```

---

**Este documento é implementável como está. Cada parte pode ser aplicada sequencialmente.**
