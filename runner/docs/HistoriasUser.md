# Histórias de Usuário Estruturadas — Plataforma Runner

## 1) Visão do produto (refatorada)

A solução final contempla **4 entregáveis**, organizados em **2 CLIs + 2 JARs**:

1. **CLI Simulador (executável)**
   - abstrai o uso do **simulador.jar**;
   - verifica/baixa o JAR correspondente quando necessário;
   - verifica/provisiona Java automaticamente quando ausente;
   - inicia, para e consulta status via comandos simples.

2. **simulador.jar**
   - serviço de simulação gerenciado pelo CLI Simulador.

3. **CLI Assinador (executável)**
   - abstrai o uso do **assinador-verificador.jar**;
   - verifica/baixa o JAR correspondente quando necessário;
   - verifica/provisiona Java automaticamente quando ausente;
   - executa assinatura/verificação em modo one-time e em modo server.

4. **assinador-verificador.jar**
   - valida parâmetros;
   - simula assinatura e verificação;
   - suporta execução por comando único (one-time) ou como servidor.

---

## 2) Diretrizes de experiência

- Público inclui usuários não técnicos.
- O usuário deve conseguir operar **mesmo sem Java previamente instalado**.
- Cada CLI deve garantir pré-checagem de dependências (JAR + Java) antes de executar.
- Comandos devem ser objetivos e com mensagens de erro acionáveis.
- `--help` deve conter exemplos prontos.
- Os dois CLIs devem ser distribuídos em **3 versões**: Windows, macOS e Linux.

---

## 3) Personas

### Persona A — Usuário comum (não técnico)
- Quer rodar comandos prontos sem conhecer Java.
- Espera mensagens claras e orientadas à ação.

---

## 4) Histórias por produto

## 4.1 CLI Simulador + simulador.jar

### US-SIM-01 — Baixar simulador automaticamente
**Como** usuário comum  
**Quero** que o CLI Simulador baixe o `simulador.jar` quando necessário  
**Para que** eu consiga executar o simulador sem setup manual.

**Critérios de aceitação**
1. O CLI verifica se o `simulador.jar` existe localmente.
2. Se ausente ou desatualizado, realiza download automático.
3. Se atualizado localmente, reutiliza o arquivo sem novo download.
4. Em falha de rede, mostra erro com instrução de tentativa.

---

### US-SIM-01B — Garantir Java para execução do simulador
**Como** usuário comum  
**Quero** que o CLI Simulador valide e provisione Java automaticamente  
**Para que** o `simulador.jar` execute mesmo em ambiente sem Java pré-instalado.

**Critérios de aceitação**
1. O CLI detecta se existe Java compatível no host antes de executar o JAR.
2. Se Java estiver ausente/incompatível, o CLI realiza download/configuração automática.
3. O processo continua no mesmo comando após provisionamento de Java.
4. Em falha de provisionamento, a mensagem orienta próximo passo.

---

### US-SIM-02 — Controlar ciclo de vida do simulador
**Como** usuário comum  
**Quero** comandos de `start`, `stop` e `status`  
**Para que** eu controle o `simulador.jar` sem comandos Java manuais.

**Critérios de aceitação**
1. `start` inicia o serviço e retorna confirmação.
2. `stop` encerra o serviço com segurança.
3. `status` informa claramente se está ativo/inativo.
4. Erros de porta em uso retornam mensagem acionável.

---

## 4.2 CLI Assinador + assinador-verificador.jar

### US-AS-01 — Executar sem Java pré-instalado
**Como** usuário comum  
**Quero** executar o CLI Assinador em ambiente sem Java  
**Para que** eu use assinatura/verificação sem instalar runtime manualmente.

**Critérios de aceitação**
1. O CLI detecta Java compatível no host.
2. Se Java ausente/incompatível, baixa e configura automaticamente.
3. Após provisionamento, a execução continua no mesmo comando.
4. Mensagens de falha orientam como corrigir/repetir a operação.

---

### US-AS-01B — Garantir presença do assinador-verificador.jar antes da execução
**Como** usuário comum  
**Quero** que o CLI Assinador valide se o JAR está disponível localmente  
**Para que** eu execute assinatura/verificação sem gerenciamento manual de artefatos.

**Critérios de aceitação**
1. O CLI verifica a presença/versão do `assinador-verificador.jar` antes de qualquer operação.
2. Se ausente ou desatualizado, faz download automático e valida integridade básica do arquivo.
3. Se já atualizado, reutiliza o artefato local sem novo download.
4. Em falha de rede/download, retorna erro acionável.

---

### US-AS-02 — Executar em modo one-time com validação de parâmetros
**Como** usuário comum  
**Quero** executar assinatura/verificação por comando único  
**Para que** eu obtenha resultado rápido sem iniciar servidor.

**Critérios de aceitação**
1. O modo one-time valida todos os parâmetros antes do processamento.
2. Parâmetros inválidos retornam erro claro e não processam operação.
3. Parâmetros válidos retornam resultado simulado de assinatura/verificação.
4. Códigos de saída refletem sucesso/falha para automação.

**Exemplo (ilustrativo)**
```bash
assinador-cli assinar --pathin ./entrada.json --pathout ./assinado.json --mode one-time
assinador-cli verificar --pathin ./assinado.json --pathout ./resultado.json --mode one-time
```

---

### US-AS-03 — Executar em modo server
**Como** usuário comum  
**Quero** subir o assinador como serviço  
**Para que** múltiplas requisições sejam atendidas com menor overhead.

**Critérios de aceitação**
1. CLI inicia o `assinador-verificador.jar` em modo server.
2. CLI informa endpoint e status do serviço.
3. Há comando para parar o serviço.
4. Falha de bind/porta retorna erro com ação sugerida.

**Exemplo (ilustrativo)**
```bash
assinador-cli server start --port 8080
assinador-cli server status
assinador-cli server stop
```

---

### US-AS-04 — Invocar operação local ou via servidor
**Como** usuário comum  
**Quero** escolher entre execução local one-time e chamada ao servidor  
**Para que** eu adapte desempenho e simplicidade ao meu cenário.

**Critérios de aceitação**
1. CLI aceita seleção explícita de modo (`one-time` ou `server`).
2. Resultado funcional mantém contrato consistente entre modos.
3. Timeout/conexão indisponível retornam erro padronizado.

---

## 5) Histórias transversais

### US-TR-01 — Abrir e operar ambos serviços no mesmo ambiente sem Java
**Como** usuário comum  
**Quero** executar CLI Simulador e CLI Assinador na mesma máquina sem Java prévio  
**Para que** eu consiga abrir e operar os dois serviços com fricção mínima.

**Critérios de aceitação**
1. Ambos CLIs funcionam sem pré-instalação manual de Java.
2. Cada CLI mantém seu fluxo de download/provisionamento sem conflito.
3. Mensagens de sucesso/erro deixam claro qual serviço foi afetado.

---

### US-TR-03 — Distribuição multiplataforma dos CLIs
**Como** usuário comum  
**Quero** versões dos CLIs para Windows, macOS e Linux  
**Para que** eu use a solução no meu sistema operacional sem adaptações complexas.

**Critérios de aceitação**
1. Existe build/publicação de `simulador-cli` para Windows, macOS e Linux.
2. Existe build/publicação de `assinador-cli` para Windows, macOS e Linux.
3. A documentação de instalação/uso diferencia comandos por sistema quando necessário.
4. O fluxo de verificação de JAR + Java é consistente entre os 3 sistemas.

---

### US-TR-02 — Ajuda amigável para usuário não técnico
**Como** usuário comum  
**Quero** `--help` com exemplos reais  
**Para que** eu execute tarefas sem suporte especializado.

**Critérios de aceitação**
1. Cada comando crítico contém exemplo completo.
2. Exemplos cobrem `pathin`, `pathout`, `one-time` e `server`.
3. Erros de uso indicam comando correto equivalente.

---

## 6) Priorização sugerida

### MVP
1. US-AS-01 (rodar sem Java pré-instalado)
2. US-SIM-01B (Java automático no CLI Simulador)
3. US-AS-01B (garantia de download do assinador-verificador.jar)
4. US-SIM-01 (download automático do simulador.jar)
5. US-AS-02 (one-time com validação obrigatória)
6. US-SIM-02 (start/stop/status)
7. US-TR-01 (operar ambos serviços no mesmo ambiente)

### Incremento 2
8. US-AS-03 (modo server)
9. US-AS-04 (escolha explícita de modo)
10. US-TR-02 (help avançado e UX)
11. US-TR-03 (distribuição Windows/macOS/Linux)

---

## 7) Definição de pronto (DoD)

- Critérios de aceitação implementados e testados.
- Mensagens de erro revisadas para linguagem simples.
- `--help` atualizado com exemplos válidos.
- Evidência de execução em ambiente sem Java pré-instalado.
- Evidência de pré-checagem de JAR e Java em ambos CLIs.
- Evidência de pacote/binário publicado para Windows, macOS e Linux.