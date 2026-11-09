# IMPLEMENTAÇÃO US-01 — Guia Prático e Direto

**Data**: 19 de maio de 2026  
**Foco**: Mudanças mínimas e reaproveitar código existente

---

## 📋 ARQUIVOS A ALTERAR

```
assinador-cli/
├── cmd/
│   ├── common.go              [MODIFICAR] Adicionar ExecutionStrategy
│   ├── assinar.go             [MODIFICAR] Usar RunWithStrategy()
│   ├── verificar.go           [MODIFICAR] Usar RunWithStrategy()
│   └── strategy_test.go       [NOVO] Testes de estratégia
│
├── internal/runner/
│   ├── runner.go              [MODIFICAR] Adicionar detecção e RunWithStrategy()
│   └── runner_test.go         [MODIFICAR] Adicionar testes de detecção
│
└── go.mod                      [SEM MUDANÇA]
```

---

## ✂️ MUDANÇA 1: common.go — Adicionar ExecutionStrategy

### Local da Mudança
**Após** `type runtimeFlags struct { ... }` (linha ~30)

### Código a Adicionar

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

// RemoveModeFromArgs remove qualquer flag --mode e seu valor dos argumentos
func RemoveModeFromArgs(args []string) []string {
	result := []string{}
	i := 0
	for i < len(args) {
		if args[i] == "--mode" {
			i += 2 // Pula flag e valor
			continue
		}
		result = append(result, args[i])
		i++
	}
	return result
}
```

### Imports Necessários em common.go

Verificar se `"net"` e `"time"` estão importados. Se não, adicionar junto com outros imports:

```go
import (
	// ... existing imports ...
	"net"
	"time"
)
```

---

## ✂️ MUDANÇA 2: runner.go — Adicionar Detecção e RunWithStrategy

### Imports Necessários em runner.go

Verificar se os seguintes imports estão presentes:

```go
import (
	// ... existing ...
	"net"
	"time"
)
```

### Código a Adicionar (antes de `func (c Config) Run(...)`)

```go
// IsServerAvailable verifica se o servidor responde na porta dada
// Usa teste de conexão TCP (funciona independente de implementação HTTP)
// Retorna true se conexão sucede, false se timeout ou recusada
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
// Para "auto": retorna "http" se servidor disponível, senão "direct"
// Para "http": retorna "http"
// Para "direct": retorna "direct"
func DetermineExecutionMode(strategy string, port int) (string, error) {
	switch strings.ToLower(strategy) {
	case "direct":
		return "direct", nil

	case "http":
		return "http", nil

	case "auto":
		// Tenta HTTP com 1 segundo de timeout
		if IsServerAvailable(port, 1) {
			return "http", nil
		}
		// Fallback para direto
		return "direct", nil

	default:
		return "", fmt.Errorf("estrategia desconhecida: %s", strategy)
	}
}

// ApplyExecutionMode modifica os argumentos do comando para setar --mode e --port
// Remove qualquer flag --mode existente, adiciona nova --mode, e --port se necessário
func ApplyExecutionMode(args []string, mode string, port int) []string {
	// Remove flag --mode existente
	filtered := RemoveModeFromArgs(args)

	// Adiciona nova --mode
	filtered = append(filtered, "--mode", mode)

	// Adiciona --port se modo é http e não está presente
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

// RunWithStrategy executa com seleção automática de modo baseada na estratégia
// Se strategy é "auto": tenta HTTP primeiro, fallback para direto se necessário
// Se strategy é "http": força modo HTTP
// Se strategy é "direct": força modo direto
func (c Config) RunWithStrategy(args []string, strategy string, port int) (Result, error) {
	if err := c.Validate(); err != nil {
		return Result{}, err
	}

	// Determina modo final baseado em estratégia e disponibilidade
	mode, err := DetermineExecutionMode(strategy, port)
	if err != nil {
		return Result{}, err
	}

	// Aplica modo aos argumentos
	finalArgs := ApplyExecutionMode(args, mode, port)

	// Executa com modo determinado
	return c.Run(finalArgs)
}
```

---

## ✂️ MUDANÇA 3: assinar.go — Integrar RunWithStrategy

### Modificar função `(o *assinarOptions) run()`

**Local**: Encontre a linha com `mode, err := modeToJarValue(o.Modo)`

**Substituir TODO O BLOCO** que começa com essa linha (até `result, err := newRunnerConfig...`) por:

```go
	// Novo: Converte --modo para estratégia
	strategy, err := ParseExecutionStrategy(o.Modo)
	if err != nil {
		return err
	}

	// Se HTTP foi forçado, valida porta
	if strategy == StrategyHTTP {
		if err := ensurePort(command, "http", o.Porta); err != nil {
			return err
		}
	}

	// Constrói argumentos para JAR (sem --mode, será adicionado por RunWithStrategy)
	args := []string{
		"sign",
		"--pathin", o.Entrada,
		"--pathout", o.Saida,
		"--alias", o.Alias,
	}

	// Adiciona PKCS#11 se fornecido
	if o.BibliotecaPKCS11 != "" {
		args = append(args, "--pkcs11-lib", o.BibliotecaPKCS11)
	}
	if o.SlotPKCS11 != "" {
		args = append(args, "--pkcs11-slot", o.SlotPKCS11)
	}

	// Novo: Usa RunWithStrategy para seleção automática de modo
	config := newRunnerConfig(o.runtimeFlags)
	result, err := config.RunWithStrategy(args, strategy.String(), o.Porta)
```

### Atualizar flag `--modo` em `newAssinarCmd()`

**Encontre** a linha com `flags.StringVar(&options.Modo...`

**Substitua** por:

```go
	flags.StringVar(&options.Modo, "modo", "auto",
		"Estrategia de execucao: auto, http ou direto.\n"+
			"  auto   (padrao): Tenta servidor HTTP, fallback para direto se indisponivel.\n"+
			"  http   : Usa apenas servidor HTTP, falha se indisponivel.\n"+
			"  direto : Usa apenas execucao direta (sem servidor).")
```

---

## ✂️ MUDANÇA 4: verificar.go — Integrar RunWithStrategy

### Modificar função `(o *verificarOptions) run()`

**Mesmo processo que assinar.go**, mas sem PKCS#11:

```go
	// Novo: Converte --modo para estratégia
	strategy, err := ParseExecutionStrategy(o.Modo)
	if err != nil {
		return err
	}

	// Validação de modo inválido
	if command.Flags().Changed("alias") || command.Flags().Changed("biblioteca-pkcs11") || command.Flags().Changed("slot-pkcs11") {
		return validationError("O comando verificar nao aceita --alias, --biblioteca-pkcs11 nem --slot-pkcs11.")
	}

	// Se HTTP foi forçado, valida porta
	if strategy == StrategyHTTP {
		if err := ensurePort(command, "http", o.Porta); err != nil {
			return err
		}
	}

	// Constrói argumentos (sem --mode)
	args := []string{
		"validate",
		"--pathin", o.Entrada,
		"--pathout", o.Saida,
	}

	// Novo: Usa RunWithStrategy
	config := newRunnerConfig(o.runtimeFlags)
	result, err := config.RunWithStrategy(args, strategy.String(), o.Porta)
```

### Atualizar flag `--modo` em `newVerificarCmd()`

Mesmo que assinar.go:

```go
	flags.StringVar(&options.Modo, "modo", "auto",
		"Estrategia de execucao: auto, http ou direto.\n"+
			"  auto   (padrao): Tenta servidor HTTP, fallback para direto se indisponivel.\n"+
			"  http   : Usa apenas servidor HTTP, falha se indisponivel.\n"+
			"  direto : Usa apenas execucao direta (sem servidor).")
```

---

## ✂️ MUDANÇA 5: Adicionar Testes

### Novo arquivo: `cmd/strategy_test.go`

```go
package cmd

import (
	"errors"
	"testing"
)

func TestParseExecutionStrategyDefault(t *testing.T) {
	strategy, err := ParseExecutionStrategy("")
	if err != nil {
		t.Fatalf("erro inesperado para string vazia: %v", err)
	}
	if strategy != StrategyAuto {
		t.Fatalf("esperava StrategyAuto, obteve %v", strategy)
	}
}

func TestParseExecutionStrategyAuto(t *testing.T) {
	strategy, err := ParseExecutionStrategy("auto")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if strategy != StrategyAuto {
		t.Fatalf("esperava StrategyAuto, obteve %v", strategy)
	}
}

func TestParseExecutionStrategyHTTP(t *testing.T) {
	strategy, err := ParseExecutionStrategy("http")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if strategy != StrategyHTTP {
		t.Fatalf("esperava StrategyHTTP, obteve %v", strategy)
	}
}

func TestParseExecutionStrategyDireto(t *testing.T) {
	strategy, err := ParseExecutionStrategy("direto")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if strategy != StrategyDirect {
		t.Fatalf("esperava StrategyDirect, obteve %v", strategy)
	}
}

func TestParseExecutionStrategyInvalid(t *testing.T) {
	_, err := ParseExecutionStrategy("invalido")
	if err == nil {
		t.Fatalf("esperava erro de validacao para estrategia invalida")
	}

	var exitErr *ExitError
	if !errors.As(err, &exitErr) || exitErr.Code != 2 {
		t.Fatalf("esperava erro de validacao com codigo 2")
	}
}

func TestParseExecutionStrategyCaseInsensitive(t *testing.T) {
	testCases := []struct {
		input    string
		expected ExecutionStrategy
	}{
		{"AUTO", StrategyAuto},
		{"HTTP", StrategyHTTP},
		{"DIRETO", StrategyDirect},
		{"Auto", StrategyAuto},
	}

	for _, tc := range testCases {
		strategy, err := ParseExecutionStrategy(tc.input)
		if err != nil {
			t.Fatalf("erro inesperado para %q: %v", tc.input, err)
		}
		if strategy != tc.expected {
			t.Fatalf("para entrada %q: esperava %v, obteve %v", tc.input, tc.expected, strategy)
		}
	}
}
```

### Adicionar a `internal/runner/runner_test.go` (ao final do arquivo)

```go
func TestIsServerAvailableWhenListening(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("falha ao criar listener: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port

	if !IsServerAvailable(port, 1) {
		t.Fatalf("esperava que servidor estivesse disponivel na porta %d", port)
	}
}

func TestIsServerAvailableWhenNotListening(t *testing.T) {
	// Usa uma porta que dificilmente estará em uso
	if IsServerAvailable(59999, 1) {
		t.Fatalf("esperava que servidor nao estivesse disponivel na porta 59999")
	}
}

func TestDetermineExecutionModeDirect(t *testing.T) {
	mode, err := DetermineExecutionMode("direct", 8080)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if mode != "direct" {
		t.Fatalf("esperava 'direct', obteve %q", mode)
	}
}

func TestDetermineExecutionModeHTTP(t *testing.T) {
	mode, err := DetermineExecutionMode("http", 8080)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if mode != "http" {
		t.Fatalf("esperava 'http', obteve %q", mode)
	}
}

func TestDetermineExecutionModeAutoWithServer(t *testing.T) {
	listener, _ := net.Listen("tcp", "localhost:0")
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port

	mode, err := DetermineExecutionMode("auto", port)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if mode != "http" {
		t.Fatalf("esperava 'http' com servidor disponivel, obteve %q", mode)
	}
}

func TestDetermineExecutionModeAutoWithoutServer(t *testing.T) {
	mode, err := DetermineExecutionMode("auto", 59999)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if mode != "direct" {
		t.Fatalf("esperava 'direct' sem servidor, obteve %q", mode)
	}
}

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
		t.Fatalf("flag --mode nao foi adicionada aos argumentos: %v", result)
	}
	if !hasPort {
		t.Fatalf("flag --port nao foi adicionada para modo http: %v", result)
	}
}

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
		t.Fatalf("esperava exatamente uma flag --mode, encontrou %d em: %v", count, result)
	}
}

func TestRemoveModeFromArgs(t *testing.T) {
	args := []string{"sign", "--mode", "direct", "--pathin", "e.json", "--mode", "http"}
	result := RemoveModeFromArgs(args)

	for i := 0; i < len(result); i++ {
		if result[i] == "--mode" {
			t.Fatalf("flag --mode nao foi removida: %v", result)
		}
	}
}
```

### Adicionar imports em runner_test.go

Garantir que `"net"` está importado:

```go
import (
	// ... existing ...
	"net"
)
```

---

## ⚠️ RISCOS E MITIGAÇÃO

| Risco | Probabilidade | Impacto | Mitigação |
|:---|:---:|:---:|:---|
| Timeout de 1s em detecção muito longo | Baixa | Médio | Testar em rede lenta, ajustar se necessário |
| Compatibilidade reversa quebrada | Baixa | Alto | Testes com `--modo direto` e `--modo http` |
| Falha na detecção causa fallback inesperado | Baixa | Médio | Logs e testes cobrem cenários |
| JAR não implementa timeout | Média | Médio | Validar com `servidor iniciar --timeout` |

---

## 🔄 SEQUÊNCIA DE COMMITS IDEAL

### Commit 1: Adicionar tipos e funções base
```bash
git add assinador-cli/cmd/common.go assinador-cli/internal/runner/runner.go
git commit -m "US-01: Adicionar ExecutionStrategy e funções de detecção

- Adicionar tipo ExecutionStrategy (auto/http/direct)
- Implementar ParseExecutionStrategy() com suporte a português/inglês
- Adicionar IsServerAvailable() para detecção de servidor
- Adicionar DetermineExecutionMode() com fallback automático
- Adicionar ApplyExecutionMode() para aplicar modo aos argumentos
- Implementar RunWithStrategy() para execução com seleção automática"
```

### Commit 2: Integrar em assinar
```bash
git add assinador-cli/cmd/assinar.go
git commit -m "US-01: Integrar RunWithStrategy em comando assinar

- Modificar função run() para usar ParseExecutionStrategy()
- Chamar RunWithStrategy() em vez de Run()
- Atualizar flag --modo com novo padrão 'auto'
- Adicionar suporte a fallback automático HTTP → direto"
```

### Commit 3: Integrar em verificar
```bash
git add assinador-cli/cmd/verificar.go
git commit -m "US-01: Integrar RunWithStrategy em comando verificar

- Modificar função run() para usar ParseExecutionStrategy()
- Chamar RunWithStrategy() em vez de Run()
- Atualizar flag --modo com novo padrão 'auto'
- Adicionar suporte a fallback automático HTTP → direto"
```

### Commit 4: Adicionar testes
```bash
git add assinador-cli/cmd/strategy_test.go assinador-cli/internal/runner/runner_test.go
git commit -m "US-01: Adicionar testes de estratégia e detecção

- Novo arquivo strategy_test.go com 6 testes de parsing
- Adicionar 8 testes de detecção e modo em runner_test.go
- Cobertura total: 14+ testes para US-01
- Testes cobrem: parsing, detecção, fallback, compatibilidade"
```

---

## ✅ CHECKLIST DE VALIDAÇÃO

Depois de implementar, executar:

```bash
# 1. Compilar
cd assinador-cli
go build .
# Esperado: Sem erros

# 2. Testes
go test ./...
# Esperado: Todos os testes passam

# 3. Cobertura
go test -cover ./...
# Esperado: common.go e runner.go com ≥80%

# 4. Testes de regressão
go test -run "TestAssinarRejectsMissingAlias|TestServidorIniciar" ./...
# Esperado: Testes antigos ainda passam

# 5. Validação manual
./assinador-cli assinar --help | grep -A3 "modo"
# Esperado: mostra "auto" como padrão

./assinador-cli assinar --entrada entrada.json --saida s.json --alias demo
# Esperado: Sem reclamação de "modo obrigatório"

./assinador-cli assinar --entrada entrada.json --saida s.json --alias demo --modo direto
# Esperado: Funciona (compatibilidade reversa)
```

---

## 🎯 ESTIMATIVA FINAL

| Etapa | Tempo |
|:---|---:|
| Ler este documento | 20 min |
| Implementar mudanças | 30-45 min |
| Executar testes | 15 min |
| Validar manualmente | 20 min |
| **TOTAL** | **1.5-2 horas** |

---

## 📝 NOTAS IMPORTANTES

1. **Compatibilidade**: Todos os comandos antigos funcionam, `--modo` é opcional agora
2. **Padrão**: `--modo` padrão agora é `auto`, não mais obrigatório
3. **Detecção**: Timeout de 1 segundo, sem erro se falhar
4. **Fallback**: HTTP → direto automático, sem intervenção do usuário
5. **Testes**: Roda com `go test ./...` como sempre

---

**Pronto para implementar? Siga os commits na sequência proposta! 🚀**
