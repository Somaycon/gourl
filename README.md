# GoURL

API de encurtamento de URLs escrita em Go, com Gin, PostgreSQL e Redis.

## Como executar

Execute os comandos abaixo na raiz do projeto, onde está o arquivo `docker-compose.yml`.

### Opção 1: executar tudo com Docker

#### Pré-requisitos

- Docker e Docker Compose v2 (no Windows, você pode usar o Docker Desktop em modo de containers Linux).
- Portas `8080`, `5432` e `6379` disponíveis.
- Acesso à internet para baixar imagens, dependências e, se necessário, o toolchain Go.

#### 1. Configure as variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto com:

```dotenv
POSTGRES_USER=gourl
POSTGRES_PASSWORD=gourl_dev
POSTGRES_DB=gourl

DSN="host=postgres user=gourl password=gourl_dev dbname=gourl port=5432 sslmode=disable"
REDIS_ADDR=redis:6379
BASE_URL=http://localhost:8080/
```

Essas credenciais são apenas para desenvolvimento local. Não publique o `.env` nem use essa senha em produção.

| Variável            | Descrição                                                       |
| ------------------- | --------------------------------------------------------------- |
| `POSTGRES_USER`     | Usuário criado pelo container PostgreSQL.                       |
| `POSTGRES_PASSWORD` | Senha desse usuário.                                            |
| `POSTGRES_DB`       | Banco criado pelo container PostgreSQL.                         |
| `DSN`               | Conexão usada pela API; deve corresponder às credenciais acima. |
| `REDIS_ADDR`        | Endereço do Redis acessível pela API.                           |
| `BASE_URL`          | Endereço público das URLs curtas. Mantenha a `/` no final.      |

Dentro do Docker Compose, os hosts são os nomes dos serviços: `postgres` e `redis`, não `localhost`.

#### 2. Inicie os serviços

```sh
docker compose up -d --build
```

O Compose inicia PostgreSQL e Redis e aguarda os healthchecks antes de iniciar a API. A aplicação cria ou atualiza as tabelas automaticamente pelo `AutoMigrate` do GORM.

A API estará disponível em `http://localhost:8080`.

Para verificar o estado e os logs:

```sh
docker compose ps
docker compose logs --tail=100 app
```

> O `go.mod` exige Go `1.27.0`. O Dockerfile usa `golang:alpine`; a imagem ou o download automático do toolchain precisa fornecer uma versão compatível. Se o build falhar por indisponibilidade dessa versão, revise o requisito do `go.mod` e a versão da imagem antes de continuar.

#### 3. Pare os serviços

```sh
docker compose down
```

Os volumes do banco e do Redis são preservados. Para remover também os dados locais, use `docker compose down -v` **somente se quiser apagar esses dados**.

### Opção 2: executar a API localmente

Além do Docker para PostgreSQL e Redis, instale uma versão do Go compatível com o `go.mod` (atualmente `1.27.0` ou superior).

1. Crie o `.env` conforme o exemplo anterior, mas altere os hosts para `localhost`:

   ```dotenv
   DSN="host=localhost user=gourl password=gourl_dev dbname=gourl port=5432 sslmode=disable"
   REDIS_ADDR=localhost:6379
   ```

   Mantenha as demais variáveis do exemplo.

2. Se a API estiver rodando no Docker, pare-a para liberar a porta `8080`:

   ```sh
   docker compose stop app
   ```

3. Inicie apenas as dependências e aguarde que estejam saudáveis:

   ```sh
   docker compose up -d --wait postgres redis
   ```

4. Baixe as dependências e execute a aplicação:

   ```sh
   go mod download
   go run ./cmd/api
   ```

A aplicação carrega o `.env` da raiz e atende em `http://localhost:8080`. Use `Ctrl+C` para encerrá-la.

Para voltar a executar a API pelo Docker, restaure os hosts `postgres` e `redis` no `.env`.

## Testando a API

### Criar uma URL curta

Envie uma requisição `POST` para `http://localhost:8080/url`, com o header `Content-Type: application/json` e este corpo:

```json
{
  "url": "https://example.com"
}
```

Exemplo com curl em Bash:

```sh
curl -X POST http://localhost:8080/url -H 'Content-Type: application/json' -d '{"url":"https://example.com"}'
```

No PowerShell:

```powershell
$resposta = Invoke-RestMethod -Method Post -Uri 'http://localhost:8080/url' -ContentType 'application/json' -Body '{"url":"https://example.com"}'
$resposta
```

A resposta contém os dados da URL criada, incluindo `code` e `short_url`.

### Acessar a URL original

Abra no navegador o valor de `short_url` retornado pela API. A rota `GET /:short` responde com um redirecionamento `302` para o endereço original, ou `404` se o código não for encontrado.

## Problemas comuns

- **Conexão recusada:** confira `docker compose ps`, os logs e os hosts definidos no `.env`. Use `postgres`/`redis` com a API no Docker e `localhost` com a API local.
- **Falha de autenticação no PostgreSQL:** confirme que o `DSN` usa as mesmas credenciais de `POSTGRES_USER`, `POSTGRES_PASSWORD` e `POSTGRES_DB`. Alterar essas variáveis não atualiza um banco já inicializado em um volume existente.
- **Porta ocupada:** encerre o processo que usa a porta ou ajuste o mapeamento no Compose e as configurações correspondentes. Evite executar a API local e a API no Docker na mesma porta.
- **Erro de versão do Go:** verifique `go version` e o requisito de toolchain no `go.mod`.
