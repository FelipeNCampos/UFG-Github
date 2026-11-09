# Planejamento de Desenvolvimento — Sistema Runner

## 0. Fase de Entendimento e Definições Iniciais:

### 0.1 Definição da Stack Tecnológica
Definir as tecnologias que serão utilizadas no projeto:
- Linguagem do CLI: Go
- Biblioteca para CLI: Cobra
- Aplicação Java: definir entre CLI simples ou framework web (ex: Spring Boot)
- Persistência local: arquivos JSON no diretório do usuário (`~/.hubsaude`)

### 0.2 Definição dos Comandos e Parâmetros
- Listar comandos disponíveis:
  - `sign`
  - `validate`
  - `simulador start`
  - `simulador stop`
  - `simulador status`
- Definir flags e parâmetros de cada comando
- Definir formato de entrada (flags, arquivos ou JSON)

### 0.3 Definição do Modelo de Dados
- Definir entradas para assinatura e validação
- Definir formato das respostas
- Garantir consistência entre CLI e Java

---

## 1. Fase de Implementação do CLI Básico (Go)

### 1.1 Inicialização do Projeto
- Criar módulo Go
- Instalar e configurar Cobra

### 1.2 Estruturação dos Comandos
```text
assinatura
├── sign
├── validate
└── simulador
    ├── start
    ├── stop
    └── status
```

### 1.3 Parsing de Argumentos
- Implementar flags
- Validar entradas básicas
- Exibir valores recebidos (modo inicial)

---

## 2. Fase de Execução de Processos Java (Modo Local)

### 2.1 Aplicação Java Inicial
- Criar projeto Java
- Implementar leitura de argumentos

### 2.2 Integração Go → Java
- Executar processo via `exec.Command`
- Passar parâmetros

### 2.3 Captura de Saída
- Capturar `stdout` e `stderr`
- Exibir resultado no CLI

---

## 3. Fase de Implementação da Lógica no Java

### 3.1 Estrutura do Projeto
- Configurar Maven ou Gradle
- Definir classe principal

### 3.2 Interface de Serviço
```java
interface SignatureService {
    String sign(String message, String alias);
    boolean validate(String message, String signature, String publicKey);
}
```

### 3.3 Implementação Simulada
- Criar `FakeSignatureService`
- Implementar lógica fake

### 3.4 Fluxos de Negócio
- Implementar `sign`
- Implementar `validate`

### 3.5 Validação de Parâmetros
- Validar obrigatoriedade
- Validar formato
- Retornar erros claros

---

## 4. Fase de Integração CLI ↔ Java (Modo Local)

### 4.1 Mapeamento de Parâmetros
- Converter flags do CLI em argumentos Java

### 4.2 Tratamento de Erros
- Interpretar erros do Java
- Exibir mensagens amigáveis

### 4.3 Formatação de Saída
- Padronizar saída
- Garantir clareza

---

## 5. Fase de Implementação do Modo Servidor (HTTP)

### 5.1 API no Java
- Criar endpoints:
  - `/sign`
  - `/validate`

### 5.2 Cliente HTTP no Go
- Enviar requisições
- Trabalhar com JSON

### 5.3 Seleção de Modo
- Usar HTTP se servidor ativo
- Caso contrário, usar execução local

---

## 6. Fase de Gerenciamento de Processos

### 6.1 Inicialização do Servidor
- Iniciar Java via CLI

### 6.2 Detecção de Instâncias
- Verificar se já está em execução

### 6.3 Encerramento
- Parar processo com segurança

### 6.4 Gerenciamento de Portas
- Definir porta padrão
- Detectar conflitos
- Selecionar porta disponível

---

## 7. Fase de Persistência Local

### 7.1 Estrutura de Diretórios

~/.hubsaude/


### 7.2 Arquivos
- `config.json`
- `processes.json`

### 7.3 Dados Persistidos
- PID
- Porta
- Versões

---

## 8. Fase de Provisionamento de JDK

### 8.1 Detecção
- Verificar se JDK está disponível

### 8.2 Download
- Baixar versão compatível

### 8.3 Configuração
- Descompactar
- Configurar uso local

---

## 9. Fase de Gerenciamento do Simulador

### 9.1 Download
- Baixar via GitHub Releases

### 9.2 Controle de Versão
- Evitar downloads desnecessários

### 9.3 Execução
- `start`
- `stop`
- `status`

---

## 10. Fase de Tratamento de Erros e Refinamento

### 10.1 Padronização
- Definir mensagens de erro

### 10.2 Tratamento de Exceções
- Garantir robustez

### 10.3 Experiência do Usuário
- Melhorar clareza
- Padronizar comportamento

---

## 11. Fase de Testes

### 11.1 Testes Unitários
- Go e Java

### 11.2 Testes de Integração
- CLI ↔ Java

### 11.3 Testes de Erro
- Cenários inválidos

### 11.4 Testes de Aceitação
- Baseados nas user stories

---

## 12. Fase de Build e Distribuição

### 12.1 CLI
- Windows
- Linux
- macOS

### 12.2 Java
- Gerar `.jar`

### 12.3 Publicação
- GitHub Releases
- Checksums SHA256

---

## 13. Fase de Segurança

### 13.1 Assinatura com Cosign
- Assinar artefatos

### 13.2 Arquivos de Verificação
- `.sig`
- `.pem`

### 13.3 Automação
- CI/CD

### 13.4 Verificação
- Garantir integridade dos binários

## 14. Visão Geral e Direcionamento

Este planejamento define uma abordagem incremental e modular para o desenvolvimento do Sistema Runner, começando por um CLI funcional em Go e evoluindo para uma arquitetura híbrida com execução local e modo servidor em Java. A separação clara entre camadas (CLI, execução de processos, API e persistência) favorece manutenibilidade, testabilidade e futura escalabilidade. Além disso, a inclusão antecipada de preocupações como gerenciamento de processos, provisionamento de JDK, versionamento e segurança garante que o sistema não seja apenas funcional, mas também robusto e pronto para distribuição em múltiplos ambientes.
