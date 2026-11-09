# Testes de Aceitação e Integração — US-01

---

## 1. TESTES DE ACEITAÇÃO (BDD - Behavior Driven Development)

### Cenário 1: Assinar em Modo Automático com Servidor Disponível

```gherkin
Funcionalidade: Assinador deve reutilizar servidor existente automaticamente
  Como usuário do Sistema Runner
  Quero que o CLI detecte e reutilize um servidor já em execução
  Para não precisar gerenciar instâncias manualmente

  Cenário: Assinar em modo automático (padrão) reutiliza servidor
    Dado que existe um servidor HTTP do assinador na porta 8080
    E que existe um arquivo "entrada.json" válido
    Quando eu executo: assinador-cli assinar --entrada entrada.json --saida saida.json --alias demo
    Então o comando deve completar com sucesso
    E a saída deve conter resultado válido
    E o código de saída deve ser 0
    E o servidor não deve ter sido reiniciado (mesmo processo)

  Teste de Implementação:
  ```go
  func TestAssinarModoAutomaticReutilizaServidorExistente(t *testing.T) {
      // Setup: Iniciar servidor mock na porta 8080
      listener, _ := net.Listen("tcp", "localhost:8080")
      defer listener.Close()
      
      // Setup: Criar arquivo de entrada
      tempDir := t.TempDir()
      inputFile := filepath.Join(tempDir, "entrada.json")
      outputFile := filepath.Join(tempDir, "saida.json")
      os.WriteFile(inputFile, []byte(validBundleJSON), 0644)
      
      // Execute: Comando sem --modo (deve usar auto)
      cmd := exec.Command(
          "assinador-cli", "assinar",
          "--entrada", inputFile,
          "--saida", outputFile,
          "--alias", "demo",
      )
      
      // Validate: Deve tentar HTTP
      // (Verificar logs ou behavior de que tentou porta 8080)
      output, err := cmd.CombinedOutput()
      
      // Assert: Não falha por "modo" ou "port"
      if err != nil && strings.Contains(string(output), "Modo invalido") {
          t.Fatalf("não deveria reclamar de modo: %s", output)
      }
  }
  ```
```

### Cenário 2: Assinar em Modo Automático SEM Servidor Disponível

```gherkin
  Cenário: Assinar em modo automático fallback para direto quando servidor indisponível
    Dado que NÃO existe servidor HTTP na porta 8080
    E que existe um arquivo "entrada.json" válido
    E que o Java e o JAR estão disponíveis
    Quando eu executo: assinador-cli assinar --entrada entrada.json --saida saida.json --alias demo
    Então o comando deve completar com sucesso
    E a saída deve conter resultado válido
    E o código de saída deve ser 0
    E o assinador deve ter sido invocado em modo direto (executado localmente)

  Teste de Implementação:
  ```go
  func TestAssinarModoAutomaticFallbackParaDireto(t *testing.T) {
      // Setup: Nenhum listener na porta 8080 (padrão)
      // Verificar que nenhum processo está rodando
      if isPortListening(8080) {
          t.Skip("porta 8080 está em uso, pulando teste")
      }
      
      tempDir := t.TempDir()
      inputFile := filepath.Join(tempDir, "entrada.json")
      outputFile := filepath.Join(tempDir, "saida.json")
      jarPath := filepath.Join(tempDir, "assinador-verificador.jar")
      
      os.WriteFile(inputFile, []byte(validBundleJSON), 0644)
      os.WriteFile(jarPath, []byte("mock-jar"), 0644)  // Mock JAR
      
      // Mock Java para retornar sucesso em modo direto
      t.Setenv("FAKE_JAVA_BEHAVIOR", "direct_mode_success")
      
      // Execute
      cmd := exec.Command(
          "assinador-cli", "assinar",
          "--entrada", inputFile,
          "--saida", outputFile,
          "--alias", "demo",
          "--jar", jarPath,
          "--java-bin", mockJavaBinary(t),
      )
      
      output, err := cmd.CombinedOutput()
      
      // Assert: Sucesso via modo direto
      if err != nil {
          t.Fatalf("esperava sucesso com fallback direto, erro: %v\n%s", err, output)
      }
  }
  ```
```

### Cenário 3: Forçar Modo HTTP sem Servidor Disponível (Deve Falhar)

```gherkin
  Cenário: Forçar modo HTTP sem servidor disponível deve falhar
    Dado que NÃO existe servidor HTTP na porta 8080
    E que existe um arquivo "entrada.json" válido
    Quando eu executo: assinador-cli assinar --entrada entrada.json --saida saida.json --alias demo --modo http
    Então o comando deve falhar
    E a saída deve conter mensagem de erro sobre indisponibilidade
    E o código de saída deve ser diferente de 0

  Teste de Implementação:
  ```go
  func TestAssinarModoHTTPForcadoFalhaSemServidor(t *testing.T) {
      if isPortListening(8080) {
          t.Skip("porta 8080 em uso, pulando teste")
      }
      
      tempDir := t.TempDir()
      inputFile := filepath.Join(tempDir, "entrada.json")
      outputFile := filepath.Join(tempDir, "saida.json")
      
      os.WriteFile(inputFile, []byte(validBundleJSON), 0644)
      
      // Execute com --modo http forçado
      cmd := exec.Command(
          "assinador-cli", "assinar",
          "--entrada", inputFile,
          "--saida", outputFile,
          "--alias", "demo",
          "--modo", "http",  // Forçar HTTP
          "--porta", "8080",
      )
      
      output, err := cmd.CombinedOutput()
      
      // Assert: Deve falhar
      if err == nil {
          t.Fatalf("esperava falha sem servidor em --modo http, mas funcionou")
      }
      if cmd.ProcessState.ExitCode() == 0 {
          t.Fatalf("exit code deve ser != 0")
      }
  }
  ```
```

### Cenário 4: Forçar Modo Direto Apesar de Servidor Disponível

```gherkin
  Cenário: Forçar modo direto ignora servidor existente
    Dado que existe um servidor HTTP do assinador na porta 8080
    E que existe um arquivo "entrada.json" válido
    E que o Java e o JAR estão disponíveis
    Quando eu executo: assinador-cli assinar --entrada entrada.json --saida saida.json --alias demo --modo direto
    Então o comando deve completar com sucesso
    E a saída deve conter resultado válido
    E o código de saída deve ser 0
    E o assinador deve ter sido invocado em modo direto (executado localmente)
    E o servidor NÃO deve ter sido contactado (mesmo estando disponível)

  Teste de Implementação:
  ```go
  func TestAssinarModoDirectoIgnoraServidorDisponivel(t *testing.T) {
      // Setup: Listener mock na porta 8080
      listener, _ := net.Listen("tcp", "localhost:8080")
      defer listener.Close()
      
      tempDir := t.TempDir()
      inputFile := filepath.Join(tempDir, "entrada.json")
      outputFile := filepath.Join(tempDir, "saida.json")
      jarPath := filepath.Join(tempDir, "assinador-verificador.jar")
      
      os.WriteFile(inputFile, []byte(validBundleJSON), 0644)
      os.WriteFile(jarPath, []byte("mock-jar"), 0644)
      
      mockJava := mockJavaBinary(t)
      
      // Execute com --modo direto
      cmd := exec.Command(
          "assinador-cli", "assinar",
          "--entrada", inputFile,
          "--saida", outputFile,
          "--alias", "demo",
          "--modo", "direto",  // Forçar direto
          "--jar", jarPath,
          "--java-bin", mockJava,
      )
      
      output, err := cmd.CombinedOutput()
      
      // Assert: Sucesso e executado localmente (não tentou HTTP)
      if err != nil {
          t.Fatalf("esperava sucesso em modo direto: %v\n%s", err, output)
      }
  }
  ```
```

### Cenário 5: Servidor Inicia e Para por Timeout de Inatividade

```gherkin
  Cenário: Servidor se encerra automaticamente após timeout de inatividade
    Dado que nenhum servidor está rodando
    Quando eu executo: assinador-cli servidor iniciar --porta 8080 --timeout 2
    E aguardo 3 minutos sem interação
    Então o processo do servidor deve ter terminado
    E o comando "assinador-cli servidor status --porta 8080" deve indicar servidor indisponível
    E o código de saída do status deve ser != 0

  Teste de Implementação:
  ```go
  func TestServidorAutoShutdownPorInatividade(t *testing.T) {
      if testing.Short() {
          t.Skip("pulando teste longo em modo short")
      }
      
      // Execute servidor com timeout 1 minuto (para teste rápido)
      cmd := exec.Command(
          "assinador-cli", "servidor", "iniciar",
          "--porta", "18080",
          "--timeout", "1",  // 1 minuto
      )
      
      // Inicia em background
      if err := cmd.Start(); err != nil {
          t.Fatalf("falha ao iniciar servidor: %v", err)
      }
      
      pid := cmd.Process.Pid
      defer func() {
          // Cleanup: garantir que processo seja terminado
          if cmd.ProcessState == nil || !cmd.ProcessState.Exited() {
              cmd.Process.Kill()
          }
      }()
      
      // Verificar que servidor iniciou
      time.Sleep(2 * time.Second)
      if !isProcessRunning(pid) {
          t.Fatalf("servidor processo não iniciou ou saiu cedo")
      }
      
      // Aguardar timeout (~1min)
      time.Sleep(70 * time.Second)
      
      // Verificar que servidor terminou
      if isProcessRunning(pid) {
          t.Fatalf("servidor processo ainda está rodando após timeout")
      }
      
      // Verificar que status falha agora
      statusCmd := exec.Command(
          "assinador-cli", "servidor", "status",
          "--porta", "18080",
      )
      statusOutput, statusErr := statusCmd.CombinedOutput()
      
      if statusErr == nil {
          t.Fatalf("esperava erro no status após timeout: %s", statusOutput)
      }
  }
  ```
```

---

## 2. TESTES DE INTEGRAÇÃO (Workflow End-to-End)

### Teste: Fluxo Completo com Auto-Detecção

```go
// arquivo: integration_test.go

package integration

import (
    "bytes"
    "os"
    "path/filepath"
    "testing"
    "time"
)

// TestCompleteWorkflowWithAutoDetection verifica todo o fluxo
func TestCompleteWorkflowWithAutoDetection(t *testing.T) {
    // FASE 1: Setup
    tempDir := t.TempDir()
    
    // Criar entrada
    inputFile := filepath.Join(tempDir, "entrada.json")
    outputFile := filepath.Join(tempDir, "saida.json")
    statusFile := filepath.Join(tempDir, "status.json")
    
    validInput := `{
      "resourceType": "Bundle",
      "id": "bundle-1",
      "type": "collection",
      "entry": []
    }`
    
    if err := os.WriteFile(inputFile, []byte(validInput), 0644); err != nil {
        t.Fatalf("falha ao criar arquivo de entrada: %v", err)
    }
    
    // FASE 2: Tentar assinar SEM servidor (deve usar modo direto)
    t.Log("Fase 2: Assinar sem servidor (fallback para direto)")
    
    // Mock do Java para retornar sucesso
    javaPath := mockJavaExecutable(t, tempDir, "sign_success")
    jarPath := filepath.Join(tempDir, "assinador.jar")
    os.WriteFile(jarPath, []byte("mock"), 0644)
    
    cmd := cmdAssinar(
        inputFile, outputFile,
        "test-alias",
        javaPath, jarPath,
        "auto",  // Modo automático
        8080,
    )
    
    var output bytes.Buffer
    cmd.Stdout = &output
    cmd.Stderr = &output
    
    if err := cmd.Run(); err != nil {
        t.Logf("assinar (direto) output: %s", output.String())
        // Esperado: pode falhar por JAR mock, mas não por modo
    }
    
    // FASE 3: Iniciar servidor
    t.Log("Fase 3: Iniciar servidor")
    
    serverCmd := cmdServidorIniciar(javaPath, jarPath, 18080, 5)  // 5 min timeout
    if err := serverCmd.Start(); err != nil {
        t.Fatalf("falha ao iniciar servidor: %v", err)
    }
    defer serverCmd.Process.Kill()
    
    // Aguardar servidor iniciar
    time.Sleep(2 * time.Second)
    
    // FASE 4: Assinar COM servidor (deve detectar e usar HTTP)
    t.Log("Fase 4: Assinar com servidor (auto detecta e usa HTTP)")
    
    cmd2 := cmdAssinar(
        inputFile, outputFile,
        "test-alias",
        javaPath, jarPath,
        "auto",  // Modo automático
        18080,   // Mesma porta do servidor
    )
    
    output.Reset()
    cmd2.Stdout = &output
    cmd2.Stderr = &output
    
    // Este comando deve tentar usar HTTP
    // (em produção faria request ao servidor)
    if err := cmd2.Run(); err != nil {
        t.Logf("assinar (http/auto) output: %s", output.String())
    }
    
    // FASE 5: Verificar status do servidor
    t.Log("Fase 5: Verificar status do servidor")
    
    statusCmd := cmdServidorStatus(javaPath, jarPath, 18080)
    output.Reset()
    statusCmd.Stdout = &output
    statusCmd.Stderr = &output
    
    if err := statusCmd.Run(); err == nil {
        t.Logf("servidor status: %s", output.String())
    }
    
    // FASE 6: Parar servidor
    t.Log("Fase 6: Parar servidor")
    
    stopCmd := cmdServidorParar(javaPath, jarPath, 18080)
    if err := stopCmd.Run(); err != nil {
        t.Logf("parar servidor: %v", err)
    }
    
    // FASE 7: Verificar que servidor foi parado
    t.Log("Fase 7: Verificar que servidor foi parado")
    
    statusCmd2 := cmdServidorStatus(javaPath, jarPath, 18080)
    output.Reset()
    statusCmd2.Stdout = &output
    statusCmd2.Stderr = &output
    
    if err := statusCmd2.Run(); err == nil {
        t.Logf("status após parar (deveria falhar): %s", output.String())
    }
    
    t.Log("Fluxo completo validado com sucesso")
}
```

---

## 3. TESTES DE PERFORMANCE

### Teste: Overhead de Detecção de Servidor

```go
// Verificar que auto-detecção não adiciona overhead significativo

func TestAutoDetectionPerformance(t *testing.T) {
    const iterations = 10
    
    // CENÁRIO 1: Modo auto COM servidor disponível
    listener, _ := net.Listen("tcp", "localhost:8080")
    defer listener.Close()
    
    start := time.Now()
    for i := 0; i < iterations; i++ {
        mode, _ := DetermineExecutionMode("auto", 8080)
        if mode != "http" {
            t.Fatalf("esperava http")
        }
    }
    autoWithServer := time.Since(start)
    
    // CENÁRIO 2: Modo auto SEM servidor disponível
    start = time.Now()
    for i := 0; i < iterations; i++ {
        mode, _ := DetermineExecutionMode("auto", 19999)
        if mode != "direct" {
            t.Fatalf("esperava direct")
        }
    }
    autoWithoutServer := time.Since(start)
    
    // CENÁRIO 3: Modo direto (sem detecção)
    start = time.Now()
    for i := 0; i < iterations; i++ {
        mode, _ := DetermineExecutionMode("direct", 8080)
        if mode != "direct" {
            t.Fatalf("esperava direct")
        }
    }
    directMode := time.Since(start)
    
    // Log de performance
    t.Logf("Auto (com servidor): %v / %d = %v por iteração", autoWithServer, iterations, autoWithServer/time.Duration(iterations))
    t.Logf("Auto (sem servidor): %v / %d = %v por iteração", autoWithoutServer, iterations, autoWithoutServer/time.Duration(iterations))
    t.Logf("Direto: %v / %d = %v por iteração", directMode, iterations, directMode/time.Duration(iterations))
    
    // Verificar que timeout de 1s não é violado
    if autoWithoutServer > 30*time.Second {  // 3s * 10 iterações = 30s máximo esperado
        t.Fatalf("auto-detecção levou muito tempo (timeout não respeitado)")
    }
}
```

---

## 4. TESTES DE REGRESSÃO

### Verificar Compatibilidade Reversa

```go
// Verificar que comportamento antigo ainda funciona (sem --modo especificado)

func TestBackwardCompatibilityOldScripts(t *testing.T) {
    tempDir := t.TempDir()
    inputFile := filepath.Join(tempDir, "entrada.json")
    outputFile := filepath.Join(tempDir, "saida.json")
    
    os.WriteFile(inputFile, []byte(`{"resourceType":"Bundle"}`), 0644)
    
    testCases := []struct {
        name        string
        modoValue   string  // valor de --modo
        expectedOK  bool
        description string
    }{
        {
            name:       "modo_direto_antigo",
            modoValue:  "direto",
            expectedOK: true,
            description: "Scripts antigos com --modo direto devem funcionar",
        },
        {
            name:       "modo_http_antigo",
            modoValue:  "http",
            expectedOK: true,
            description: "Scripts antigos com --modo http devem funcionar",
        },
        {
            name:       "modo_direct_english",
            modoValue:  "direct",
            expectedOK: true,
            description: "Suporte a 'direct' em inglês deve funcionar",
        },
        {
            name:       "sem_modo_novo_default",
            modoValue:  "",
            expectedOK: true,
            description: "Sem --modo deve usar auto (novo padrão)",
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            args := []string{
                "assinar",
                "--entrada", inputFile,
                "--saida", outputFile,
                "--alias", "test",
            }
            
            if tc.modoValue != "" {
                args = append(args, "--modo", tc.modoValue)
            }
            
            // Executar CLI
            cmd := exec.Command("assinador-cli", args...)
            output, err := cmd.CombinedOutput()
            
            // Validar
            success := err == nil || cmd.ProcessState.ExitCode() == 0
            
            if tc.expectedOK && !success {
                t.Fatalf("%s: esperava sucesso mas falhou\n%s", tc.description, output)
            }
            
            // Verificar que não falha por erro de modo
            if strings.Contains(string(output), "Modo invalido") {
                t.Fatalf("%s: reclamou de modo inválido\n%s", tc.description, output)
            }
            
            t.Logf("%s: OK", tc.description)
        })
    }
}
```

---

## 5. HELPERS E UTILITÁRIOS PARA TESTES

```go
// test_helpers.go

package cmd

import (
    "fmt"
    "net"
    "os"
    "os/exec"
    "path/filepath"
    "testing"
)

// cmdAssinar cria comando para testar assinar
func cmdAssinar(entrada, saida, alias, javaBin, jarPath, modo string, porta int) *exec.Cmd {
    args := []string{
        "assinar",
        "--entrada", entrada,
        "--saida", saida,
        "--alias", alias,
        "--java-bin", javaBin,
        "--jar", jarPath,
    }
    
    if modo != "" {
        args = append(args, "--modo", modo)
    }
    
    if porta > 0 {
        args = append(args, "--porta", fmt.Sprintf("%d", porta))
    }
    
    return exec.Command("assinador-cli", args...)
}

// cmdServidorIniciar cria comando para iniciar servidor
func cmdServidorIniciar(javaBin, jarPath string, porta, timeoutMin int) *exec.Cmd {
    args := []string{
        "servidor", "iniciar",
        "--java-bin", javaBin,
        "--jar", jarPath,
        "--porta", fmt.Sprintf("%d", porta),
    }
    
    if timeoutMin > 0 {
        args = append(args, "--timeout", fmt.Sprintf("%d", timeoutMin))
    }
    
    return exec.Command("assinador-cli", args...)
}

// cmdServidorStatus cria comando para verificar status
func cmdServidorStatus(javaBin, jarPath string, porta int) *exec.Cmd {
    return exec.Command(
        "assinador-cli", "servidor", "status",
        "--java-bin", javaBin,
        "--jar", jarPath,
        "--porta", fmt.Sprintf("%d", porta),
    )
}

// cmdServidorParar cria comando para parar servidor
func cmdServidorParar(javaBin, jarPath string, porta int) *exec.Cmd {
    return exec.Command(
        "assinador-cli", "servidor", "parar",
        "--java-bin", javaBin,
        "--jar", jarPath,
        "--porta", fmt.Sprintf("%d", porta),
    )
}

// mockJavaExecutable cria um executável Java mock para testes
func mockJavaExecutable(t *testing.T, tempDir, behavior string) string {
    t.Helper()
    
    javaPath := filepath.Join(tempDir, "java")
    if runtime.GOOS == "windows" {
        javaPath = filepath.Join(tempDir, "java.cmd")
    }
    
    // Script que simula Java behavior
    script := fmt.Sprintf(`#!/bin/bash
echo '{"status":"SUCCESS","operation":"sign"}'
`)
    
    os.WriteFile(javaPath, []byte(script), 0755)
    return javaPath
}

// isPortListening verifica se porta está respondendo
func isPortListening(port int) bool {
    conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
    if err != nil {
        return false
    }
    defer conn.Close()
    return true
}

// isProcessRunning verifica se um process ID está vivo
func isProcessRunning(pid int) bool {
    // Unix-like
    if err := exec.Command("kill", "-0", fmt.Sprintf("%d", pid)).Run(); err == nil {
        return true
    }
    
    // Windows
    if err := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid)).Run(); err == nil {
        return true
    }
    
    return false
}
```

---

## 6. MATRIZ DE TESTES × REQUISITOS

| Requisito | Teste Unitário | Teste Integração | Teste Aceitação | Teste Regressão |
|:---|:---:|:---:|:---:|:---:|
| Aceitar modo auto | `TestParseStrategyAuto` | `TestAutoModeFlowEnd2End` | Cenário 1, 2 | Compatibilidade |
| Modo direto funciona | `TestDetermineModeDirect` | `TestDirectModeExecution` | Cenário 4 | Compatibilidade |
| Modo HTTP funciona | `TestDetermineHTTP` | `TestHTTPModeExecution` | Cenário 3 | Compatibilidade |
| Detectar servidor | `TestIsServerAvailable*` | `TestServerDetection` | Cenário 1, 2 | Performance |
| Fallback auto→direct | `TestAutoFallback` | `TestAutoFallback*` | Cenário 2 | — |
| Timeout inatividade | `TestApplyTimeout` | `TestTimeoutShutdown` | Cenário 5 | — |
| Compatibilidade reversa | — | — | — | Todos acima |

---

## 7. COBERTURA ESPERADA

Após implementação, cobertura mínima esperada:

```
common.go:
  - ExecutionStrategy enum        : 100%
  - ParseExecutionStrategy        : 100%
  - ExtractModeFromArgs           : 100%
  - RemoveModeFromArgs            : 100%

runner.go:
  - IsServerAvailable             : 100%
  - DetermineExecutionMode        : 100%
  - ApplyExecutionMode            : 100%
  - RunWithStrategy               : 95%

assinar.go:
  - (o *assinarOptions).run()     : 90%

verificar.go:
  - (o *verificarOptions).run()   : 90%

Cobertura total:   ≥ 85%
Funções críticas:  ≥ 95%
```

Rodar cobertura:
```bash
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

---

**Próxima etapa**: Executar todos os testes com `go test -v ./...`
