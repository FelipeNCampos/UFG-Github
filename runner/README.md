# Runner

Sistema Runner para orquestrar aplicacoes Java por CLIs multiplataforma, escondendo do usuario os detalhes de Java, JARs, modos de execucao e ciclo de vida de servidores.

## Status da entrega

O repositorio contem os tres componentes principais do trabalho:

- `assinador-verificador/`: JAR Java que simula assinatura, validacao e modo servidor.
- `assinador-cli/`: CLI Go para assinar, verificar e gerenciar o servidor do assinador.
- `simulador-cli/`: CLI Go para iniciar, consultar e parar o `simulador.jar`, com provisionamento automatico do JAR e do JRE quando necessario.

Tambem foram incluidos:

- workflow de CI em `.github/workflows/ci.yml`;
- workflow de release em `.github/workflows/release.yml`;
- script de empacotamento em `scripts/build-release.ps1`;
- checksums SHA256 e assinatura Cosign automatizados no fluxo de release.

## Estrutura

```text
.
├── assinador-cli/          # CLI do assinador em Go
├── assinador-verificador/  # aplicacao Java empacotada como JAR
├── simulador-cli/          # CLI do simulador em Go
├── docs/                   # especificacao, historias e documentos de apoio
├── diagramas/              # diagramas C4/PlantUML e imagens geradas
├── scripts/                # automacao de build/release
└── .github/workflows/      # CI e publicacao de releases
```

## Validacao local

Requisitos para validar tudo localmente:

- Go 1.22+
- JDK 17+ com `javac`, `jar` e `java` no PATH
- PowerShell 7+ recomendado para execucao multiplataforma dos scripts

```powershell
cd assinador-cli
go test ./...

cd ../simulador-cli
go test ./...

cd ../assinador-verificador
./test.ps1
```

Para gerar os artefatos de release localmente:

```powershell
./scripts/build-release.ps1 -Version 1.0.0
```

Os arquivos serao gravados em `dist/`, incluindo `SHA256SUMS.txt`.

## Uso do assinador-cli

O `assinador-cli` executa o `assinador-verificador.jar` diretamente ou via servidor HTTP.

Exemplos:

```powershell
assinador-cli assinar `
  --entrada .\entrada.json `
  --saida .\assinado.json `
  --modo auto `
  --alias demo-signer

assinador-cli verificar `
  --entrada .\assinado.json `
  --saida .\resultado.json `
  --modo http `
  --porta 8080

assinador-cli servidor iniciar --porta 8080
assinador-cli servidor status --porta 8080
assinador-cli servidor parar --porta 8080
```

Modos aceitos:

- `auto`: usa HTTP quando encontra servidor ativo na porta informada; caso contrario usa execucao direta.
- `http`: usa apenas servidor HTTP.
- `direto`: executa o JAR diretamente.

Quando `--java-bin` nao e informado, o CLI provisiona um Temurin JRE 17 em diretorio local do usuario.

## Uso do simulador-cli

O `simulador-cli` gerencia o ciclo de vida do `simulador.jar`.

```powershell
simulador-cli iniciar --porta 8443
simulador-cli status --porta 8443
simulador-cli parar --porta 8443
```

O JAR do simulador e baixado dinamicamente a partir do `release.json` configurado no codigo ou pelo ambiente:

```powershell
$env:SIMULADOR_RELEASE_URL = "https://exemplo/release.json"
```

Quando Java nao esta disponivel, o CLI provisiona um Temurin JRE 21 em `~/.hubsaude/jre`.

## Release

Uma release e publicada automaticamente ao enviar uma tag SemVer:

```powershell
git tag v1.0.0
git push origin v1.0.0
```

O workflow gera binarios para Windows, Linux e macOS amd64 dos dois CLIs, gera o JAR do assinador, calcula checksums SHA256, assina cada artefato com Cosign keyless/OIDC e publica tudo no GitHub Releases.

Para verificar um artefato assinado:

```bash
cosign verify-blob \
  --certificate assinatura-1.0.0-linux-amd64.pem \
  --signature assinatura-1.0.0-linux-amd64.sig \
  assinatura-1.0.0-linux-amd64
```

## Documentacao

- [Especificacao](docs/especificacao.md)
- [Historias de usuario](docs/HistoriasUser.md)
- [Documento final de entrega](docs/ENTREGA_FINAL.md)
- [Design tecnico](docs/design.md)
