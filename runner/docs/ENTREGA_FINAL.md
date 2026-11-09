# Entrega final do Sistema Runner

Este documento resume o estado final do projeto e o procedimento recomendado para avaliacao, build e publicacao.

## Componentes entregues

| Componente | Local | Estado |
|---|---|---|
| CLI do assinador | `assinador-cli/` | Implementado em Go com comandos `assinar`, `verificar` e `servidor` |
| Assinador Java | `assinador-verificador/` | Implementado em Java, com build para JAR e testes automatizados |
| CLI do simulador | `simulador-cli/` | Implementado em Go com comandos `iniciar`, `status` e `parar` |
| Provisionamento de Java | `assinador-cli/internal/javaruntime/`, `simulador-cli/internal/provisioner/` | Implementado para runtimes Temurin |
| Provisionamento do simulador.jar | `simulador-cli/internal/provisioner/jar.go` | Implementado via `release.json` remoto |
| CI | `.github/workflows/ci.yml` | Testa Go e Java em Windows, Linux e macOS |
| Release | `.github/workflows/release.yml` | Gera binarios, checksums, assinaturas Cosign e GitHub Release |

## Como validar

Ambiente necessario:

- Go 1.22+
- JDK 17+
- PowerShell 7+

Comandos:

```powershell
cd assinador-cli
go test ./...

cd ../simulador-cli
go test ./...

cd ../assinador-verificador
./test.ps1
```

## Como gerar artefatos

```powershell
./scripts/build-release.ps1 -Version 1.0.0
```

Saida esperada em `dist/`:

- `assinatura-1.0.0-windows-amd64.exe`
- `assinatura-1.0.0-linux-amd64`
- `assinatura-1.0.0-macos-amd64`
- `simulador-1.0.0-windows-amd64.exe`
- `simulador-1.0.0-linux-amd64`
- `simulador-1.0.0-macos-amd64`
- `assinador-verificador-1.0.0.jar`
- `SHA256SUMS.txt`

## Como publicar release

Crie e envie uma tag SemVer:

```powershell
git tag v1.0.0
git push origin v1.0.0
```

O workflow `Release` executa automaticamente:

- build do JAR Java;
- cross-build dos CLIs Go para Windows, Linux e macOS amd64;
- geracao de `SHA256SUMS.txt`;
- assinatura de cada artefato com Cosign keyless/OIDC;
- publicacao dos arquivos no GitHub Releases.

Para cada artefato publicado, o release tambem inclui:

- `<artefato>.sig`
- `<artefato>.pem`

## Observacoes de avaliacao

- O ambiente desta sessao nao possuia `go` nem `javac` no PATH, entao a validacao completa deve ser feita no CI ou em uma maquina com esses requisitos instalados.
- O script `scripts/build-release.ps1` foi criado para ser a fonte unica de empacotamento local e remoto.
- Os binarios Linux e macOS sao executaveis Go nativos. Caso a avaliacao exija especificamente `.AppImage` ou `.dmg`, estes formatos podem ser adicionados como uma etapa posterior de empacotamento visual, sem alterar os CLIs.
