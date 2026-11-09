# Diagnóstico e Plano de Implementação — US-01
## Sistema Runner | CLI Assinador

**Data**: 19/05/2026  
**Contexto**: Análise para aderência total do repositório `ufgfelipencampos/runner` à US-01 da especificação oficial.

---

## 1. DIAGNÓSTICO DO GAP ATUAL

### 1.1 Situação Positiva (Já Implementado)

| Aspecto | Status | Evidência |
|---------|--------|-----------|
| CLI em português | ✅ Completo | Comandos `assinar`, `verificar`, `servidor` em pt-BR |
| Mapeamento de flags | ✅ Completo | `--modo direto|http` → `--mode direct|http` no JAR |
| Modo direto | ✅ Funcional | `Run()` executa `java -jar assinador.jar sign ...` |
| Modo servidor | ✅ Funcional | `StartServer()` inicia e monitora processo em background |
| Validações básicas | ✅ Funcional | Validação de caminhos JSON, alias, PKCS#11 |
| Tratamento de erros | ✅ Funcional | Códigos de saída coerentes (0=sucesso, 2=validação, 1=erro) |
| Testes unitários | ✅ Parcial | Testes em Go e Java existem, cobertura aceitável |

### 1.2 Gaps Críticos para US-01

#### **GAP 1: Modo servidor NÃO é padrão**

**Requisito (Especificação)**:
> "O CLI deve preferir o modo servidor quando não orientado para usar modo local."

**Situação Atual**:
```go
// Em assinar.go, linha ~50
flags.StringVar(&options.Modo, "modo", "", "Modo de execucao: direto ou http.")
```
- A flag `--modo` é **obrigatória** (validação em `common.go`)
- Não existe **modo padrão**
- Não existe estratégia de **fallback automático**

**Impacto**: 
- Usuário DEVE especificar `--modo http` ou `--modo direto` sempre
- Contra o requisito: "preferir modo servidor" implica ser o padrão

**Evidência de Falta**:
```bash
# Comando exigido atualmente:
assinador-cli assinar --entrada entrada.json --saida saida.json --modo http --alias demo

# Esperado (com modo padrão):
assinador-cli assinar --entrada entrada.json --saida saida.json --alias demo
# → Deveria tentar modo http primeiro
```

---

#### **GAP 2: Detecção de instância ativa NÃO é comprovada**

**Requisito (Especificação)**:
> "O CLI deve detectar se já existe instância do assinador.jar em execução no modo servidor e reutilizá-la, quando apropriado."

**Situação Atual**:
- Existe comando `servidor status --porta 8080` que verifica status
- **MAS**: Não há verificação **automática antes de assinar/verificar** em modo HTTP
- Se usuário invoca `assinar --modo http`, o CLI assume que servidor está pronto
- Se servidor não está rodando, falha com erro genérico

**Código Atual** (`assinar.go`, linha ~60):
```go
func (o *assinarOptions) run(command *cobra.Command, _ []string) error {
    // ... validações ...
    result, err := newRunnerConfig(o.runtimeFlags).Run(args)
    // Não há verificação prévia de disponibilidade do servidor!
    if err != nil {
        return wrapRuntimeError(err)
    }
```

**Impacto**:
- Usuário precisa iniciar servidor manualmente: `assinador-cli servidor iniciar`
- Não há detecção silenciosa de instância existente
- Cenário esperado NÃO funciona:
  ```bash
  # Terminal 1: iniciar servidor
  assinador-cli servidor iniciar --porta 8080
  
  # Terminal 2: assinar (deveria reutilizar servidor do terminal 1)
  assinador-cli assinar --entrada e.json --saida s.json --modo http
  # Hoje não há verificação se servidor já está rodando
  ```

**Evidência de Falta**:
- Não existe função `CheckServerStatus()` ou `isServerRunning()`
- Não existe lógica de retry ou fallback em `runner.go`

---

#### **GAP 3: Timeout de inatividade NOT YET TESTED**

**Requisito (Especificação)**:
> "O CLI deve permitir interrupção programada após minutos sem interação."

**Situação Atual**:
```go
// Em servidor.go, linha ~75
flags.IntVar(&options.Timeout, "timeout", 0, "Tempo limite de inatividade, em minutos.")
```
- Flag existe e é mapeada para `--timeout` do JAR
- **MAS**: Comportamento NÃO está documentado nem testado end-to-end
- Não há confirmação de que JAR implementa shutdown por inatividade

**Impacto**:
- Feature existe mas não foi validada
- Risco: flag mapeia para JAR mas JAR pode não implementar completamente

**Recomendação**:
- Testar comportamento: iniciar servidor com `--timeout 2`, aguardar 3 min, verificar se processo foi interrompido

---

#### **GAP 4: Definição de estratégia de modo ambígua**

**Situação Atual**:
- Usuário escolhe modo explicitamente: `--modo direto` OU `--modo http`
- Não há "tentativa inteligente" de qual modo usar
- Não há fallback se modo preferido falhar

**Requisito Implícito**:
- "Preferir servidor" sugere modo HTTP como padrão
- "Suportar modo local" sugere fallback se servidor indisponível

**Proposta de Solução**:
```
--modo auto     (padrão)   → Tenta HTTP, fallback para direto
--modo http               → Só HTTP, falha se indisponível  
--modo direto             → Só direto, sem servidor
```

---

### 1.3 Quadro Resumido de Conformidade

| Critério de Aceitação US-01 | Status | Gap |
|:---|:---:|:---:|
| Aceitar comandos para criação/validação de assinatura | ✅ | Nenhum |
| Invocar Assinador com parâmetros fornecidos | ✅ | Nenhum |
| Suportar invocação direta (modo local) | ✅ | Nenhum |
| Suportar invocação via HTTP (modo servidor) | ✅ | Nenhum |
| Exibir resultado de forma legível | ✅ | Nenhum |
| **[FALTA]** Modo servidor é padrão | ❌ | GAP 1 |
| **[FALTA]** Detectar e reutilizar instância ativa | ❌ | GAP 2 |
| **[FALTA]** Interrupção por inatividade (comprovado) | ⚠️ | GAP 3 |
| Iniciar servidor na porta padrão | ✅ | Nenhum |
| Permitir parar servidor | ✅ | Nenhum |

---

## 2. PLANO DE IMPLEMENTAÇÃO MÍNIMO

### 2.1 Decisões de Projeto Recomendadas

#### **Decisão 1: Modo Padrão**
```
RECOMENDAÇÃO: --modo auto (padrão)
- Tenta conectar ao servidor HTTP na porta default (8080)
- Se falhar (timeout curto ~2s), executa modo direto (compatível com fallback)
- Permite `--modo http` (força servidor, falha se indisponível)
- Permite `--modo direto` (força modo local)

JUSTIFICATIVA:
✓ Atende "preferir servidor"
✓ Compatível com "invocação local" como fallback
✓ Sem quebra compatibilidade com scripts existentes
✓ Usuário pode forçar modo se desejar
```

#### **Decisão 2: Detecção de Instância Ativa**
```
RECOMENDAÇÃO: Verificação silenciosa com timeout curto
- Antes de assinar/verificar em modo HTTP, verificar se porta está acessível
- Timeout = 1-2 segundos (rápido para falhar rápido)
- Se porta responde: reutiliza servidor
- Se porta não responde: fallback para direto (decisão 1)

IMPLEMENTAÇÃO:
1. Criar função IsServerAvailable(port int) bool em runner.go
2. Chamar antes de executar operação em modo HTTP
3. Retornar bool, sem exception
```

#### **Decisão 3: Validação de Timeout**
```
RECOMENDAÇÃO: Testar e documentar inatividade
- Confirmar que JAR Java implementa shutdown por inatividade
- Adicionar teste de integração: init server, aguardar timeout, verificar morte
- Documentar timeout > 0 em --timeout como obrigatório se desejar auto-shutdown
```

---

### 2.2 Tarefas Técnicas Ordenadas (Etapa 1)

#### **TAREFA 1: Refatorar modo com padrão e estratégia**

**Arquivo**: `assinador-cli/cmd/common.go`

**Mudanças**:
1. Criar novo tipo enum para estratégia de modo:
```go
type ExecutionStrategy string

const (
    StrategyAuto   ExecutionStrategy = "auto"   // Default: try HTTP, fallback to direct
    StrategyHTTP   ExecutionStrategy = "http"   // Force HTTP only
    StrategyDirect ExecutionStrategy = "direct" // Force direct only
)

func (s ExecutionStrategy) Valid() error {
    if s == StrategyAuto || s == StrategyHTTP || s == StrategyDirect {
        return nil
    }
    return validationError("Estrategia de modo invalida: %s", s)
}
```

2. Manter compatibilidade com alias "modo":
```go
func parseExecutionStrategy(rawInput string) (ExecutionStrategy, error) {
    switch strings.ToLower(strings.TrimSpace(rawInput)) {
    case "", "auto":      // empty = auto (new default)
        return StrategyAuto, nil
    case "http":          // http → http
        return StrategyHTTP, nil
    case "direto":        // direto → direct
        return StrategyDirect, nil
    case "direct":        // direct (for JAR compatibility)
        return StrategyDirect, nil
    default:
        return "", validationError("Modo invalido: %s. Use --modo auto|http|direto.", rawInput)
    }
}
```

3. Update flag binding:
```go
// OLD: flags.StringVar(&options.Modo, "modo", "", "...")
// NEW: 
flags.StringVar(&options.Modo, "modo", "auto", "Estrategia de execucao: auto, http ou direto. Padrao: auto (tenta servidor, fallback para direto).")
```

**Entregável**: 
- Novo tipo `ExecutionStrategy` 
- Função de parsing com suporte a "auto"
- Compatibilidade com "direto" → "direct" existente

---

#### **TAREFA 2: Implementar detecção de servidor disponível**

**Arquivo**: `assinador-cli/internal/runner/runner.go`

**Mudanças**:
1. Adicionar função de health check:
```go
import "net/http"

// IsServerAvailable checks if server is responding on given port within timeout
func IsServerAvailable(port int, timeoutSeconds int) bool {
    url := fmt.Sprintf("http://localhost:%d/health", port)  // Assumir JAR expõe /health ou endpoint similar
    
    client := &http.Client{
        Timeout: time.Duration(timeoutSeconds) * time.Second,
    }
    
    resp, err := client.Get(url)
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    
    return resp.StatusCode == http.StatusOK
}

// Alternativa: tentar conexão TCP se HTTP não disponível
func IsServerAvailable(port int, timeoutSeconds int) bool {
    conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Duration(timeoutSeconds)*time.Second)
    if err != nil {
        return false
    }
    defer conn.Close()
    return true
}
```

2. Refatorar `Run()` para suportar estratégia:
```go
func (c Config) RunWithStrategy(args []string, strategy ExecutionStrategy) (Result, error) {
    if err := c.Validate(); err != nil {
        return Result{}, err
    }
    
    // Determinar modo final baseado em estratégia
    finalMode, err := c.determineMode(args, strategy)
    if err != nil {
        return Result{}, err
    }
    
    // Substituir --mode nos args
    args = replaceMode(args, finalMode)
    
    return c.Run(args)
}

func (c Config) determineMode(args []string, strategy ExecutionStrategy) (string, error) {
    currentMode := extractMode(args)  // Extrair --mode current
    
    switch strategy {
    case StrategyDirect:
        return "direct", nil  // Sempre direct
    case StrategyHTTP:
        return "http", nil    // Sempre http
    case StrategyAuto:
        // Tentar HTTP primeiro
        if IsServerAvailable(8080, 1) {  // 1s timeout
            return "http", nil
        }
        // Fallback para direto
        return "direct", nil
    default:
        return "", fmt.Errorf("estrategia desconhecida: %v", strategy)
    }
}
```

**Entregável**:
- Função `IsServerAvailable(port, timeoutSecs) bool`
- Função `RunWithStrategy(args, strategy) Result`
- Suporte a fallback HTTP → Direct

---

#### **TAREFA 3: Atualizar comandos assinar/verificar para usar estratégia**

**Arquivo**: `assinador-cli/cmd/assinar.go` + `verificar.go`

**Mudanças**:
```go
type assinarOptions struct {
    runtimeFlags
    Entrada   string
    Saida     string
    Modo      string  // Mantém compatibilidade
    Alias     string
    // ... outros campos ...
}

func (o *assinarOptions) run(command *cobra.Command, _ []string) error {
    // ... validações existentes ...
    
    // NOVO: Converter modo para estratégia
    strategy, err := parseExecutionStrategy(o.Modo)
    if err != nil {
        return err
    }
    
    // Construir args para JAR
    args := []string{"sign", "--pathin", o.Entrada, "--pathout", o.Saida}
    // Nota: NOT adding --mode here, será adicionado por RunWithStrategy
    
    // NOVO: Usar RunWithStrategy
    config := newRunnerConfig(o.runtimeFlags)
    result, err := config.RunWithStrategy(args, strategy)
    if err != nil {
        return wrapRuntimeError(err)
    }
    
    return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}
```

**Entregável**:
- Atualizar `assinar.go` para usar `RunWithStrategy()`
- Atualizar `verificar.go` para usar `RunWithStrategy()`
- Manter flag `--modo` existente com novo default `auto`

---

#### **TAREFA 4: Validar e testar comportamento de timeout**

**Arquivo**: `assinador-cli/cmd/servidor.go` + testes

**Validação**:
1. Confirmar JAR implementa `--timeout` para `server start`:
   ```bash
   # Testar manualmente
   java -jar assinador-verificador.jar server start --port 8080 --timeout 2
   # Aguardar 3 minutos
   # Processo deve terminar após inatividade
   ```

2. Se JAR NÃO implementa: adicionar suporte no CLI Go via goroutine + kill:
   ```go
   func (c Config) StartServerWithInactivityTimeout(args []string, timeoutMinutes int) (Result, error) {
       // ... inicia processo ...
       
       if timeoutMinutes > 0 {
           go func() {
               time.Sleep(time.Duration(timeoutMinutes) * time.Minute)
               command.Process.Kill()  // Ou enviar SIGTERM
           }()
       }
       // ...
   }
   ```

**Testes** (`cmd/servidor_test.go`):
```go
func TestServerStartWithTimeoutShutdown(t *testing.T) {
    // 1. Iniciar servidor com --timeout 1
    // 2. Verificar PID processo
    // 3. Aguardar ~2 min (tempo de teste > timeout)
    // 4. Verificar se processo foi terminado
    // 5. Tentar conexão deve falhar
}
```

**Entregável**:
- Teste de integração de timeout
- Documentação no README sobre `--timeout`
- Validação de aderência

---

### 2.3 Estrutura de Arquivos Modificados

```
assinador-cli/
├── cmd/
│   ├── common.go              [MODIFICAR] Adicionar ExecutionStrategy
│   ├── assinar.go             [MODIFICAR] Usar RunWithStrategy
│   ├── verificar.go           [MODIFICAR] Usar RunWithStrategy
│   ├── servidor.go            [VALIDAR] Testar --timeout
│   └── root_test.go           [ADICIONAR] Testes integração estratégia
│
├── internal/runner/
│   ├── runner.go              [MODIFICAR] IsServerAvailable + RunWithStrategy
│   └── runner_test.go         [ADICIONAR] Testes unit/integ server detection
│
└── go.mod                      [SEM MUDANÇA] Sem novas dependências
```

---

## 3. PROPOSTA DE REFATORAÇÃO: EXEMPLOS DE CÓDIGO

### 3.1 Novo tipo ExecutionStrategy (common.go)

```go
package cmd

import (
    "fmt"
    "strings"
)

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

// ToJarMode converts ExecutionStrategy to JAR command-line value
// Note: For StrategyAuto, this is determined at runtime
func (s ExecutionStrategy) ToJarMode() string {
    switch s {
    case StrategyHTTP:
        return "http"
    case StrategyDirect:
        return "direct"
    case StrategyAuto:
        // Determined at runtime, not here
        return ""
    default:
        return ""
    }
}
```

### 3.2 Função de Detecção de Servidor (runner.go)

```go
package runner

import (
    "fmt"
    "net"
    "strings"
    "time"
)

// IsServerAvailable checks if the server is responding on given port
// Uses TCP connection test (works even if /health endpoint not available)
// Returns true if connection succeeds, false if timeout or refused
func IsServerAvailable(port int, timeoutSeconds int) bool {
    address := fmt.Sprintf("localhost:%d", port)
    timeout := time.Duration(timeoutSeconds) * time.Second
    
    conn, err := net.DialTimeout("tcp", address, timeout)
    if err != nil {
        return false
    }
    defer conn.Close()
    
    return true
}

// DetermineExecutionMode decides final execution mode based on strategy
// For StrategyAuto: returns "http" if server available, else "direct"
// For StrategyHTTP: returns "http"
// For StrategyDirect: returns "direct"
func DetermineExecutionMode(strategy string, port int) (string, error) {
    switch strategy {
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

// ApplyExecutionMode modifies the command arguments to set --mode
func ApplyExecutionMode(args []string, mode string, port int) []string {
    // Remove any existing --mode flag
    filtered := []string{}
    i := 0
    for i < len(args) {
        if args[i] == "--mode" {
            i += 2  // Skip flag and value
            continue
        }
        filtered = append(filtered, args[i])
        i++
    }
    
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
```

### 3.3 Refatorar Run() para suportar estratégia (runner.go)

```go
// Run executes assinador.jar with given arguments (direct or HTTP)
// Strategy is determined by --mode in args; if not present, uses direct
func (c Config) Run(args []string) (Result, error) {
    if err := c.Validate(); err != nil {
        return Result{}, err
    }
    
    command := exec.Command(c.JavaBin, c.jarArgs(args)...)
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    command.Stdout = &stdout
    command.Stderr = &stderr
    
    err := command.Run()
    result := Result{
        Stdout:   stdout.String(),
        Stderr:   stderr.String(),
        ExitCode: 0,
    }
    
    if err == nil {
        return result, nil
    }
    
    var exitErr *exec.ExitError
    if errors.As(err, &exitErr) {
        result.ExitCode = exitErr.ExitCode()
        return result, nil
    }
    
    return Result{}, fmt.Errorf(
        "Falha ao executar Java. Verifique --java-bin. Erro: %w", err,
    )
}

// RunWithStrategy executes with automatic mode selection based on strategy
// If strategy is "auto": tries HTTP first, falls back to direct if needed
// If strategy is "http": enforces HTTP mode
// If strategy is "direct": enforces direct mode
func (c Config) RunWithStrategy(args []string, strategy string, port int) (Result, error) {
    if err := c.Validate(); err != nil {
        return Result{}, err
    }
    
    // Determine final mode
    mode, err := DetermineExecutionMode(strategy, port)
    if err != nil {
        return Result{}, err
    }
    
    // Apply mode to arguments
    finalArgs := ApplyExecutionMode(args, mode, port)
    
    // Execute
    return c.Run(finalArgs)
}
```

### 3.4 Atualizar assinar.go

```go
// Exemplo de mudança em assinar.go:

func (o *assinarOptions) run(command *cobra.Command, _ []string) error {
    // Validações existentes
    if err := ensureValidJSONPaths(o.Entrada, o.Saida); err != nil {
        return err
    }
    if err := ensureValidAlias(o.Alias); err != nil {
        return err
    }
    
    // NOVO: Parse execution strategy
    strategy, err := ParseExecutionStrategy(o.Modo)  // Default: "auto"
    if err != nil {
        return err
    }
    
    // NOVO: Validar se strategy HTTP foi forçado mas porta é inválida
    if strategy == "http" || (strategy == "auto" && o.Modo != "") {
        if o.Porta < 1 || o.Porta > 65535 {
            return validationError("Porta invalida: %d", o.Porta)
        }
    }
    
    // Build args for JAR (sem --mode, será adicionado por RunWithStrategy)
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
    
    // NOVO: Use RunWithStrategy for automatic mode selection
    config := newRunnerConfig(o.runtimeFlags)
    result, err := config.RunWithStrategy(args, strategy.String(), o.Porta)
    if err != nil {
        return wrapRuntimeError(err)
    }
    
    return emitResult(command.OutOrStdout(), command.ErrOrStderr(), result)
}
```

---

## 4. LISTA DE TESTES ESPERADOS

### 4.1 Testes Unitários (runner_test.go)

```go
// Test execution strategy parsing and mode determination

func TestStrategyAutoDetectsAvailableServer(t *testing.T) {
    // Setup: mock server on port 8080
    // Call: DetermineExecutionMode("auto", 8080)
    // Expect: mode == "http"
}

func TestStrategyAutoFallsBackToDirectWhenServerUnavailable(t *testing.T) {
    // Setup: no server on port 8080
    // Call: DetermineExecutionMode("auto", 8080)
    // Expect: mode == "direct"
}

func TestStrategyHTTPEnforcesServerMode(t *testing.T) {
    // Call: DetermineExecutionMode("http", 9999)
    // Expect: mode == "http" (regardless of server state)
}

func TestStrategyDirectEnforcesDirect(t *testing.T) {
    // Call: DetermineExecutionMode("direct", 8080)
    // Expect: mode == "direct"
}

func TestIsServerAvailableReturnsTrue(t *testing.T) {
    // Setup: listening server on ephemeral port
    // Call: IsServerAvailable(port, 1)
    // Expect: true
}

func TestIsServerAvailableReturnsFalse(t *testing.T) {
    // Call: IsServerAvailable(9999, 1)  // unused port
    // Expect: false
}

func TestApplyExecutionModeUpdatesArgs(t *testing.T) {
    // Input: ["sign", "--pathin", "e.json", "--pathout", "s.json"]
    // Call: ApplyExecutionMode(..., "http", 8080)
    // Expect: contains "--mode http --port 8080"
}
```

### 4.2 Testes de Integração (root_test.go)

```go
// Test full workflows with strategy

func TestAssinarWithStrategyAuto_ServerAvailable(t *testing.T) {
    // Setup:
    //   - Start real mock server on 8080
    //   - Write valid entrada.json
    // Call: executeRootCommand("assinar", "--entrada", "e.json", "--saida", "s.json", "--alias", "test")
    //   (--modo not specified, defaults to auto)
    // Expected:
    //   - Should connect to HTTP server (auto detected it)
    //   - Result should succeed
    //   - Exit code == 0
}

func TestAssinarWithStrategyAuto_ServerUnavailable(t *testing.T) {
    // Setup:
    //   - No server on 8080
    //   - Write valid entrada.json
    //   - Mock JAVA_BIN to return success
    // Call: executeRootCommand("assinar", "--entrada", "e.json", "--saida", "s.json", "--alias", "test", "--modo", "auto")
    // Expected:
    //   - Should fallback to direct mode (detected no server)
    //   - Execute locally via Java
    //   - Success
}

func TestAssinarWithStrategyHTTP_ServerRequired(t *testing.T) {
    // Setup:
    //   - No server on 8080
    //   - Write valid entrada.json
    // Call: executeRootCommand("assinar", ..., "--modo", "http")
    // Expected:
    //   - Should attempt HTTP connection
    //   - Should fail with connection error (server not found)
    //   - Exit code != 0
}

func TestAssinarWithStrategyDirect_IgnoresServer(t *testing.T) {
    // Setup:
    //   - Start mock server on 8080
    //   - Write valid entrada.json
    //   - Mock JAVA_BIN to succeed
    // Call: executeRootCommand("assinar", ..., "--modo", "direto")
    // Expected:
    //   - Should NOT attempt HTTP
    //   - Should execute directly via Java
    //   - Success (server ignored)
}

func TestVerificarStrategyBehaviorMirrorsAssinar(t *testing.T) {
    // Verificar should behave identically to assinar regarding mode selection
    // Test: auto, http, direto strategies work same way
}

func TestServerStartStopsCorrectly(t *testing.T) {
    // Setup: existing timeout flag
    // Call: executeRootCommand("servidor", "iniciar", "--porta", "8081", "--timeout", "1")
    // Expected:
    //   - Server starts successfully
    //   - After 1 minute of inactivity, process should terminate
    //   - Manual stop via "servidor parar" should work before timeout
}
```

### 4.3 Testes de Aceitação (Cenários Funcionales)

```
Cenário 1: Assinar em modo automático com servidor disponível
  Dado: servidor já está rodando em porta 8080
  Quando: usuário executa: assinador-cli assinar --entrada entrada.json --saida saida.json --alias demo
  Esperado: ✓ Usa modo HTTP, reutiliza servidor
  Exit code: 0

Cenário 2: Assinar em modo automático SEM servidor disponível
  Dado: nenhum servidor rodando
  Quando: usuário executa: assinador-cli assinar --entrada entrada.json --saida saida.json --alias demo
  Esperado: ✓ Fallback para modo direto, executa localmente
  Exit code: 0

Cenário 3: Forçar modo HTTP sem servidor disponível
  Dado: nenhum servidor rodando
  Quando: usuário executa: assinador-cli assinar --entrada entrada.json --saida saida.json --alias demo --modo http
  Esperado: ✗ Falha com erro de conexão
  Exit code: != 0

Cenário 4: Forçar modo direto apesar de servidor disponível
  Dado: servidor está rodando em porta 8080
  Quando: usuário executa: assinador-cli assinar --entrada entrada.json --saida saida.json --alias demo --modo direto
  Esperado: ✓ Ignora servidor, executa localmente
  Exit code: 0

Cenário 5: Servidor inicia com timeout e para por inatividade
  Dado: nenhum processamento
  Quando: usuário executa: assinador-cli servidor iniciar --porta 8080 --timeout 2
  E aguarda 3 minutos
  Esperado: ✓ Processo do servidor é terminado após 2 min de inatividade
  Verificação: ps aux | grep java (processo não aparece)
```

---

## 5. CRITÉRIOS DE PRONTO PARA US-01

### 5.1 Requisitos Funcionais Implementados

- [ ] **CLI aceita modo padrão**: Comando sem `--modo` funciona com auto-detecção
- [ ] **Detecção de servidor ativa**: Antes de assinar/verificar, verifica disponibilidade de porta
- [ ] **Fallback automático**: HTTP → Direto funciona sem erro do usuário
- [ ] **Modo forçado funciona**: `--modo http` e `--modo direto` funcionam como esperado
- [ ] **Timeout de inatividade testado**: Server inicia com `--timeout`, para após inatividade

### 5.2 Código

- [ ] `common.go`: Novo tipo `ExecutionStrategy` com parsing
- [ ] `runner.go`: Funções `IsServerAvailable()` e `RunWithStrategy()`
- [ ] `assinar.go` e `verificar.go`: Integrados com `RunWithStrategy()`
- [ ] `servidor.go`: Validação de `--timeout` confirmada com JAR
- [ ] Sem quebra de compatibilidade com scripts existentes

### 5.3 Testes

- [ ] Testes unitários: Estratégia de modo + detecção de servidor (5+ testes)
- [ ] Testes de integração: Fluxo completo assinar/verificar com auto-detect (4+ cenários)
- [ ] Testes de aceitação: Cobertura de 5 cenários listados em 4.3
- [ ] Cobertura mínima: 80% em `runner.go` e `common.go` (novas funções)

### 5.4 Documentação

- [ ] README atualizado: Exemplos com `--modo auto` (novo padrão)
- [ ] README atualizado: Comportamento de fallback documentado
- [ ] Comentários de código: Funções de estratégia explicadas
- [ ] Changelog: Mudança de padrão documentada

### 5.5 Validação Manual

```bash
# Teste 1: Auto com servidor
assinador-cli servidor iniciar --porta 8080
assinador-cli assinar --entrada entrada.json --saida saida.json --alias demo
# Esperado: Sucesso, usa HTTP

# Teste 2: Auto sem servidor
assinador-cli servidor parar --porta 8080
assinador-cli assinar --entrada entrada.json --saida saida.json --alias demo
# Esperado: Sucesso, fallback para direto

# Teste 3: Forçar HTTP
assinador-cli assinar --entrada entrada.json --saida saida.json --alias demo --modo http
# Esperado: Erro (sem servidor)

# Teste 4: Forçar direto
assinador-cli servidor iniciar --porta 8080
assinador-cli assinar --entrada entrada.json --saida saida.json --alias demo --modo direto
# Esperado: Sucesso, local (ignora servidor)

# Teste 5: Status OK
assinador-cli servidor status --porta 8080
# Esperado: Servidor rodando na porta 8080
```

---

## 6. ROADMAP DE ENTREGA (ETAPA 1)

| Sprint | Tarefa | Estimativa | Predecessor |
|--------|--------|-----------|------------|
| Sprint 1 | Tarefa 1: Refatorar modo com estratégia | 2-3h | — |
| Sprint 1 | Tarefa 2: Implementar detecção de servidor | 3-4h | Tarefa 1 |
| Sprint 2 | Tarefa 3: Atualizar assinar/verificar | 2-3h | Tarefa 2 |
| Sprint 2 | Tarefa 4: Validar timeout | 1-2h | Paralelo |
| Sprint 3 | Testes: Unitários + Integração | 4-5h | Tarefas 1-4 |
| Sprint 3 | Documentação: README + Exemplos | 2h | Tarefas 1-4 |
| **Total** | **Implementação US-01** | **14-20h** | — |

---

## 7. PRÓXIMAS ETAPAS (Fora do Escopo de US-01)

Estas atividades são necessárias para o "Runner completo" mas **não** impedem aceitação de US-01:

1. **Etapa 2**: CLI do Simulador (`simulador-cli`)
2. **Etapa 3**: Provisionamento automático de Java + JARs
3. **Etapa 4**: Testes expandidos (cobertura 90%+)
4. **Etapa 5**: Binários multiplataforma (Windows/Linux/macOS)
5. **Etapa 6**: Assinatura de artefatos (Cosign)
6. **Etapa 7**: Documentação técnica completa

---

## 8. REFERÊNCIAS

- **Especificação**: `especificacao.md` (Seção 5 - US-01)
- **Design**: `design.md` (Diagramas e sequências)
- **Plano Ação**: `plano_acao_runner.md` (Etapa 1)
- **Código Atual**: 
  - `assinador-cli/cmd/*.go` (Cobra CLI)
  - `assinador-cli/internal/runner/runner.go` (Execução)
  - `assinador-verificador/**` (JAR Java - contrato)

---

**Documento gerado**: 19/05/2026  
**Autor**: Engenheiro de Software Sênior | GitHub Copilot  
**Status**: Pronto para revisão e implementação
