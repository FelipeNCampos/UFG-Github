# Plano de ação para aderência do repositório `ufgfelipencampos/runner` à especificação do Sistema Runner

Este documento consolida os requisitos que ainda não estão plenamente atendidos no repositório analisado e propõe um plano de ação em etapas para evolução do projeto. A comparação foi feita entre o repositório do aluno e a documentação externa do `kyriosdata/runner` fornecida em anexo [file:106][file:107][file:108], além dos arquivos já analisados do repositório `ufgfelipencampos/runner` [cite:1][cite:2][cite:3][cite:5][cite:6][cite:9][cite:10][cite:11][cite:12].

## Situação atual

O repositório atual já implementa com boa maturidade o núcleo do **CLI do Assinador**, incluindo comandos para assinar, verificar e operar o modo servidor, com validações de parâmetros, integração com `assinador-verificador.jar` e tratamento de saída/erros [cite:2][cite:6][cite:9][cite:10][cite:11][cite:12]. Ainda assim, a especificação oficial exige um escopo mais amplo, incluindo CLI do simulador, provisionamento automático de Java, download dinâmico de artefatos, distribuição multiplataforma e assinatura de releases [file:108][file:107].

## Requisitos não atendidos ou atendidos parcialmente

| Área | Requisito | Situação | Evidência | Ação necessária |
|---|---|---|---|---|
| US-01 | Detectar instância já em execução do `assinador.jar` em modo servidor e reutilizá-la | Parcial | O código possui `server start/status/stop`, mas não evidencia descoberta e reuso automático de instância existente [cite:9][cite:12][file:108] | Implementar verificação prévia de status antes de iniciar novo processo |
| US-01 | Usar modo servidor por padrão quando não orientado para modo local | Parcial | O CLI exige `--modo` e não indica fallback automático para HTTP como comportamento padrão [cite:10][cite:11][file:108] | Definir estratégia padrão de execução e fallback configurável |
| US-01 | Requisição de interrupção programada após minutos sem interação | Parcial | Existe flag `--timeout` em `servidor iniciar`, mas a aderência completa ao critério funcional não ficou comprovada no comportamento final [cite:12][file:108] | Confirmar no JAR e testar semanticamente a inatividade |
| US-02 | Simulação completa do comportamento do `assinador.jar` | Parcial | O CLI está pronto para integração, mas a inspeção atual não comprovou toda a lógica simulada do módulo Java [cite:2][cite:9][file:108] | Auditar e completar o módulo Java com testes de aceitação |
| US-02 | Suporte PKCS#11 efetivo | Parcial | O CLI aceita flags PKCS#11, mas não ficou demonstrado o uso funcional completo no JAR [cite:10][file:107][file:108] | Validar integração no módulo Java e documentar cenários suportados |
| US-03 | CLI do simulador | Não atendido | A especificação exige CLI próprio para gerenciar o Simulador; o design externo também o prevê [file:107][file:108] | Criar módulo `simulador-cli` com `iniciar`, `status` e `parar` |
| US-03 | Verificação da porta 8443 antes de iniciar simulador | Não atendido | Não foi encontrada implementação correspondente [file:108] | Implementar checagem de porta ocupada/livre antes do start |
| US-03 | Obtenção dinâmica do `simulador.jar` via GitHub Releases | Não atendido | Não foi encontrada lógica de download/versionamento dinâmico [file:108] | Implementar `release.json`, resolução de versão e cache local |
| US-03 | Não baixar novamente artefato já atualizado | Não atendido | Sem mecanismo encontrado de cache/versionamento local [file:108] | Persistir metadados de versão instalada e comparar com release |
| US-04 | Provisionamento automático de JDK/JRE | Não atendido | O projeto atual depende de `java` previamente instalado ou `--java-bin` informado manualmente [cite:2][cite:6][cite:9][file:108] | Implementar download e uso local do runtime em `.hubsaude` |
| US-04 | Download multiplataforma do runtime | Não atendido | Não foi encontrada lógica específica para Windows/Linux/macOS [file:108] | Abstrair URLs por SO/arquitetura e automatizar instalação |
| US-05 | Binários pré-compilados para Windows/Linux/macOS | Não atendido | Não houve evidência de releases multiplataforma publicadas [file:108] | Criar pipeline de build e empacotamento |
| US-05 | Checksums SHA256 nos artefatos | Não atendido | Não houve evidência de geração/publicação de checksums [file:108] | Gerar checksum por artefato e anexar à release |
| US-05 | Versionamento SemVer e GitHub Releases | Parcial | Não ficou comprovado pelo material inspecionado [file:108] | Definir estratégia de versionamento e processo de release |
| Entregáveis | Assinatura dos artefatos com Cosign | Não atendido | A especificação exige `.sig` e `.pem` para cada artefato [file:108] | Adicionar assinatura automática no pipeline de release |
| Entregáveis | Testes de integração e aceitação amplos | Parcial | Há testes no CLI, mas a cobertura completa exigida não foi comprovada [cite:5][file:108] | Expandir suíte automatizada por requisito funcional |
| Entregáveis | Documentação de instalação e operação completa | Parcial | Há README e exemplos, mas a especificação pede manual, guia de instalação e documentação técnica de integração [file:106][file:108][cite:2] | Estruturar documentação por público: usuário, dev e operação |

## Etapas de desenvolvimento

## Etapa 1 — Fechar lacunas do Assinador

Objetivo: concluir a aderência funcional do fluxo já iniciado no `assinador-cli` [cite:10][cite:11][cite:12].

Atividades:
- Implementar política padrão de execução: tentar modo servidor primeiro e cair para local quando configurado ou necessário [file:108].
- Adicionar detecção de instância ativa do `assinador.jar` antes de subir novo processo [file:108][cite:12].
- Confirmar e testar o comportamento do `--timeout` como desligamento por inatividade [cite:12][file:108].
- Revisar o módulo Java para garantir que a simulação de assinatura, validação e erros atende aos cenários exigidos [file:108].
- Validar ponta a ponta o suporte PKCS#11 no JAR, não apenas no CLI [cite:10][file:107][file:108].

Entregáveis da etapa:
- CLI do assinador aderente à US-01.
- Testes automatizados dos modos direto e HTTP.
- Documento de comportamento esperado dos modos de execução.

## Etapa 2 — Implementar o CLI do Simulador

Objetivo: criar o segundo CLI previsto pela arquitetura do Runner [file:107][file:108].

Atividades:
- Criar módulo `simulador-cli` com comandos `iniciar`, `status` e `parar` [file:108].
- Implementar checagem da porta 8443 antes da inicialização [file:108].
- Integrar consumo dos endpoints esperados, incluindo `/api/info` para status e `/shutdown` para parada [file:108].
- Definir padrão de saída legível e códigos de retorno consistentes, alinhados ao CLI do assinador [cite:6][file:108].

Entregáveis da etapa:
- CLI do simulador funcional.
- Testes de start/status/stop.
- Manual rápido de operação do simulador.

## Etapa 3 — Implementar gerenciador de dependências

Objetivo: remover a dependência de instalação manual de Java e da disponibilidade manual dos JARs, conforme a proposta central do Runner [file:108].

Atividades:
- Criar camada de provisionamento em diretório local, por exemplo `.hubsaude/` [file:108].
- Implementar download do JRE/JDK por sistema operacional e arquitetura, usando as URLs previstas na especificação [file:108].
- Implementar download dinâmico de `simulador.jar` e, idealmente, também de `assinador.jar`, baseado em `release.json` [file:108].
- Persistir metadados locais de versão instalada para evitar downloads desnecessários [file:108].
- Integrar esse gerenciador aos dois CLIs antes da execução de qualquer operação [file:107][cite:3].

Entregáveis da etapa:
- Módulo de dependências reutilizável.
- Cache/versionamento local de runtime e JARs.
- Testes de provisionamento em Windows, Linux e macOS.

## Etapa 4 — Fortalecer testes e qualidade

Objetivo: comprovar aderência dos requisitos por automação [file:108].

Atividades:
- Mapear cada user story para testes unitários, de integração e de aceitação [file:108].
- Criar cenários de erro explícitos: arquivo ausente, JSON inválido, porta ocupada, JAR ausente, Java ausente, timeout, parâmetros indevidos [cite:6][cite:9][file:108].
- Medir cobertura mínima para componentes críticos: parser de argumentos, execução de processos, provisionamento e status de servidor.
- Padronizar códigos de saída e mensagens de erro entre todos os CLIs [cite:6][file:108].

Entregáveis da etapa:
- Suíte de testes ampliada.
- Matriz requisito x teste.
- Relatório simples de cobertura e cenários críticos.

## Etapa 5 — Empacotamento e releases

Objetivo: atender os requisitos de distribuição multiplataforma [file:108].

Atividades:
- Gerar binários para Windows amd64, Linux amd64 e macOS amd64 [file:108].
- Empacotar nomes e formatos esperados conforme estratégia do projeto [file:108].
- Gerar checksums SHA256 para cada artefato [file:108].
- Adotar SemVer no versionamento do projeto e publicar via GitHub Releases [file:108].

Entregáveis da etapa:
- Primeira release multiplataforma.
- Arquivos de checksum anexados.
- Procedimento de release documentado.

## Etapa 6 — Segurança da cadeia de suprimentos

Objetivo: cumprir a exigência de assinatura de artefatos com Cosign [file:108].

Atividades:
- Configurar pipeline CI/CD para assinatura automática de artefatos via Cosign e OIDC [file:108].
- Publicar para cada artefato os arquivos correspondentes `.sig` e `.pem` [file:108].
- Documentar o processo de verificação para usuários finais com exemplos de comando `cosign verify-blob` [file:108].

Entregáveis da etapa:
- Releases assinadas.
- Evidência de verificação de integridade.
- Documentação de verificação para usuários.

## Etapa 7 — Documentação final

Objetivo: fechar os entregáveis documentais exigidos pela disciplina [file:106][file:108].

Atividades:
- Produzir manual do usuário do assinador [file:108].
- Produzir manual do usuário do simulador [file:108].
- Produzir guia de instalação sem pré-requisitos manuais, contemplando provisionamento automático [file:108].
- Produzir documentação técnica da integração entre CLI, JARs e gerenciador de dependências [file:108][file:107].
- Revisar os diagramas arquiteturais para refletir a implementação final [file:107][cite:3].

Entregáveis da etapa:
- Conjunto documental completo.
- Exemplos de uso por cenário.
- Arquitetura atualizada e coerente com o código.

## Ordem recomendada de execução

A sequência mais segura é: **Etapa 1 -> Etapa 2 -> Etapa 3 -> Etapa 4 -> Etapa 5 -> Etapa 6 -> Etapa 7**. Essa ordem prioriza primeiro o fechamento funcional do núcleo da aplicação, depois a expansão do escopo, em seguida a automação operacional e, por fim, a distribuição segura e a documentação final [file:108][file:107].

## Priorização prática

Se houver pouco tempo, a priorização mínima recomendada é:
- Prioridade alta: Etapas 1, 2 e 3.
- Prioridade média: Etapa 4.
- Prioridade alta para entrega final: Etapas 5 e 7.
- Prioridade obrigatória para conformidade completa da especificação: Etapa 6 [file:108].

## Observação final

O repositório atual já oferece uma base sólida para o **Assinador**, especialmente no CLI e na orquestração do `assinador-verificador.jar` [cite:2][cite:6][cite:9][cite:10][cite:11][cite:12]. O maior esforço restante está na transformação do projeto em um **Runner completo**, com simulador, provisionamento automático, distribuição multiplataforma e pipeline de release seguro, exatamente como exige a especificação oficial [file:108][file:107].
