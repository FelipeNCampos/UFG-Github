# assinador-verificador

Modulo Java do `assinador-verificador.jar`.

## Escopo atual

- `sign` em modo `direct`
- `validate` em modo `direct`
- `sign` em modo `http`
- `validate` em modo `http`
- `server start`, `server status` e `server stop`
- validacao forte de flags e arquivos
- resposta JSON padronizada em `stdout` e no arquivo de saida
- simulacao deterministica de assinatura via digest SHA-256

## Contrato real de modos

O contrato atual do JAR usa estes valores:

- `--mode direct`
- `--mode http`

O `assinador-cli` em Go expoe `--modo direto|http` e faz o mapeamento automaticamente.

## Exemplos de uso direto do JAR

### Sign em modo direct

```powershell
java -jar .\build\dist\assinador-verificador.jar `
  sign `
  --pathin .\examples\fhir-bundle.json `
  --pathout .\build\saida\assinado.json `
  --mode direct `
  --alias demo-signer `
  --pkcs11-lib token.dll `
  --pkcs11-slot 0
```

### Validate em modo direct

```powershell
java -jar .\build\dist\assinador-verificador.jar `
  validate `
  --pathin .\build\saida\assinado.json `
  --pathout .\build\saida\resultado-validacao.json `
  --mode direct
```

### Sign e validate via HTTP

```powershell
java -jar .\build\dist\assinador-verificador.jar `
  server start `
  --port 8080

java -jar .\build\dist\assinador-verificador.jar `
  sign `
  --pathin .\examples\fhir-bundle.json `
  --pathout .\build\saida\assinado-http.json `
  --mode http `
  --port 8080 `
  --alias demo-signer

java -jar .\build\dist\assinador-verificador.jar `
  validate `
  --pathin .\build\saida\assinado-http.json `
  --pathout .\build\saida\resultado-validacao-http.json `
  --mode http `
  --port 8080

java -jar .\build\dist\assinador-verificador.jar `
  server stop `
  --port 8080
```

## Exit codes

| Codigo | Constante | Significado |
| --- | --- | --- |
| `0` | `SUCCESS` | Operacao concluida com sucesso. |
| `1` | `RUNTIME_ERROR` | Erro inesperado em tempo de execucao. |
| `2` | `VALIDATION_ERROR` | Argumento invalido, arquivo ausente ou conteudo incorreto. |
| `3` | `SERVER_RUNNING` | Servidor iniciado e mantido ativo intencionalmente. |

## Premissas desta versao

- `sign` exige um JSON de entrada nao vazio contendo o campo `"resourceType"`.
- `validate` espera como entrada um JSON gerado pelo comando `sign`.
- `--alias` e obrigatorio somente em `sign`.
- `--pkcs11-lib` e `--pkcs11-slot` sao opcionais em `sign`.
- `validate` nao aceita `--alias`, `--pkcs11-lib` nem `--pkcs11-slot`.
- `--port` so se aplica a `http` e aos subcomandos de `server`.
- `--timeout` so se aplica a `server start`.

## Scripts locais

- `build.ps1`: compila as classes e empacota o JAR.
- `test.ps1`: compila o modulo e executa os testes sem dependencias externas.
