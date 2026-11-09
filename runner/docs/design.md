# Sistema Runner - Design Arquitetural (Mermaid)

Este documento consolida os diagramas arquiteturais em **Markdown + Mermaid**, com foco nas histórias de usuário atuais:
- o usuário interage com o **Sistema Runner** (que expõe os dois CLIs);
- os dois JARs são tratados como **sistemas externos**;
- o Runner valida dependências e orquestra a comunicação com os JARs.

---

## 1) Diagrama de Contexto

```mermaid
flowchart LR
    U[Ator: Usuário Comum]
    R["Sistema Runner<br/>(CLI Assinador + CLI Simulador)"]
    AJ["Sistema Externo:<br/>assinador-verificador.jar"]
    SJ["Sistema Externo:<br/>simulador.jar"]

        U -->|Comandos CLI| R
    R -->|Assinar / Verificar<br/>one-time ou server| AJ
    R -->|Start / Status / Stop<br/>simulador| SJ
    AJ -->|Resultado da operação| R
    SJ -->|Status e confirmação| R
    R -->|Feedback e códigos de saída| U
```

### Tabela descritiva (Contexto)

| Elemento | Tipo | Responsabilidade |
|---|---|---|
| Usuário Comum | Ator | Executa comandos nos CLIs sem precisar conhecer Java. |
| Sistema Runner | Sistema principal | Interface única de execução (CLI Assinador + CLI Simulador), valida dependências e orquestra chamadas. |
| assinador-verificador.jar | Sistema externo | Processa assinatura/verificação em modo one-time ou server. |
| simulador.jar | Sistema externo | Fornece serviço de simulação controlado por start/status/stop. |

---

## 2) Diagrama de Contêineres (visão interna do Runner)

```mermaid
flowchart LR
    U[Usuário Comum]

    subgraph R[Sistema Runner]
      AC[Container: CLI Assinador]
      SC[Container: CLI Simulador]
      DEP["Container: Gerenciador de Dependências<br/>(Java + download de JAR)"]
    end

    AJ["Sistema Externo:<br/>assinador-verificador.jar"]
    SJ["Sistema Externo:<br/>simulador.jar"]

    U -->|Comandos assinar/verificar/server| AC
    U -->|Comandos start/status/stop| SC

    AC -->|verifica/provisiona Java<br/>verifica/baixa jar| DEP
    SC -->|verifica/provisiona Java<br/>verifica/baixa jar| DEP

    AC -->|invoca operações| AJ
    SC -->|controla ciclo de vida| SJ
```

### Tabela descritiva (Contêineres)

| Origem | Destino | Interface | Objetivo |
|---|---|---|---|
| Usuário Comum | CLI Assinador | CLI | Disparar assinatura/verificação e comandos server. |
| Usuário Comum | CLI Simulador | CLI | Controlar simulador com start/status/stop. |
| CLI Assinador | Gerenciador de Dependências | Chamada interna | Garantir Java e JAR antes de executar operações. |
| CLI Simulador | Gerenciador de Dependências | Chamada interna | Garantir Java e JAR antes de controlar o simulador. |
| CLI Assinador | assinador-verificador.jar | Processo local/HTTP | Executar one-time ou gerir server. |
| CLI Simulador | simulador.jar | Processo local/HTTP | Iniciar, consultar status e parar simulador. |

---

## 3) Diagrama de Sequência — Assinador (one-time e server)

```mermaid
sequenceDiagram
    actor U as Usuário Comum
    participant AC as CLI Assinador (Runner)
    participant DEP as Gerenciador de Dependências
    participant AJ as assinador-verificador.jar (Externo)

    U->>AC: comando assinar/verificar/server
    AC->>DEP: verificar jar do assinador
    DEP-->>AC: jar pronto (baixado ou já existente)
    AC->>DEP: verificar Java compatível
    DEP-->>AC: Java pronto (provisionado ou já existente)

    alt modo one-time
        AC->>AC: validar parâmetros (pathin/pathout)
        alt parâmetros inválidos
            AC-->>U: erro + código != 0
        else parâmetros válidos
            AC->>AJ: executar operação one-time
            AJ-->>AC: resultado assinatura/verificação
            AC-->>U: sucesso + código 0
        end
    else modo server
        AC->>AJ: start/status/stop
        AJ-->>AC: endpoint/estado/confirmação
        AC-->>U: retorno formatado
    end
```

### Tabela descritiva (Sequência Assinador)

| Etapa | Descrição |
|---|---|
| Pré-checagem | Runner garante presença de `assinador-verificador.jar` e Java antes da execução. |
| One-time | Runner valida parâmetros e só invoca o JAR se entrada estiver correta. |
| Server | Runner controla ciclo de vida do serviço (`start`, `status`, `stop`). |

---

## 4) Diagrama de Sequência — Simulador (start/status/stop)

```mermaid
sequenceDiagram
    actor U as Usuário Comum
    participant SC as CLI Simulador (Runner)
    participant DEP as Gerenciador de Dependências
    participant SJ as simulador.jar (Externo)

    U->>SC: comando start/status/stop
    SC->>DEP: verificar simulador.jar
    DEP-->>SC: jar pronto (baixado ou já existente)
    SC->>DEP: verificar Java compatível
    DEP-->>SC: Java pronto (provisionado ou já existente)

    alt start
        SC->>SJ: iniciar serviço
        SJ-->>SC: confirmação + porta
        SC-->>U: simulador ativo
    else status
        SC->>SJ: consultar status
        SJ-->>SC: ativo/inativo
        SC-->>U: status formatado
    else stop
        SC->>SJ: encerrar serviço
        SJ-->>SC: confirmação
        SC-->>U: simulador encerrado
    end
```

### Tabela descritiva (Sequência Simulador)

| Etapa | Descrição |
|---|---|
| Pré-checagem | Runner valida presença de `simulador.jar` e Java antes dos comandos. |
| Start | Inicia o `simulador.jar` e devolve confirmação ao usuário. |
| Status | Consulta estado atual do serviço de simulação. |
| Stop | Finaliza o serviço de simulação com retorno explícito ao usuário. |
