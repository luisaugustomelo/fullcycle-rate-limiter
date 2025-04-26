#FullCycle Rate Limiter

## Descrição

Este projeto implementa um Rate Limiter em Go, utilizando Redis como mecanismo de persistência. O objetivo é limitar o número de requisições que podem ser feitas:

- Por endereço IP

- Por token de acesso (passado no header API_KEY)

Caso o limite de requisições seja ultrapassado:

- é retornado o código HTTP 429 Too Many Requests

- é aplicada uma penalização de tempo (block), impedindo novas requisições até o tempo expirar

Todas as informações de limites são armazenadas no Redis.

Comportamento da Aplicação

- IP-Based Limiting:

Limita as requisições feitas por cada IP.

Exemplo: Se o limite é 2 req/s, a 3ª será bloqueada.

- Token-Based Limiting:

Se o header API_KEY estiver presente, o limite configurado para o token é usado.

Exemplo: Token TESTTOKEN pode ter limite diferente do IP.

- Prioridade:

Se houver token, o limite é baseado no token.

Caso contrário, é baseado no IP.

- Bloqueio Temporário:

Se ultrapassado, um "block" é gravado no Redis para impedir novas requisições temporariamente.

- Expiração do Block:

Após o tempo de penalização, o IP/token volta a poder fazer requisições.

## Estrutura do Projeto

```
rate-limiter/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── limiter/
│   │   ├── limiter.go
│   │   ├── redis_strategy.go
│   │   └── strategy.go
│   └── middleware/
│       └── limiter.go
│       └──rate_limiter_test.go
├── .env
├── docker-compose.yml
├── go.mod
└── go.sum
````

## Variáveis de Ambiente

Configure no arquivo .env:

```
RATE_LIMIT_IP=2
RATE_BLOCK_DURATION_SECONDS=10
RATE_LIMIT_TOKEN_TESTTOKEN=2
```

Onde:

- RATE_LIMIT_IP define o limite de requisições por IP.

- RATE_BLOCK_DURATION_SECONDS define a duração do bloqueio.

- RATE_LIMIT_TOKEN_<TOKEN> define o limite para cada token.

Como Rodar Localmente

1. Subir o Redis

Certifique-se de ter o Docker instalado. Rode o seguinte comando:

```
docker-compose up -d
```

2. Rodar os Testes

Este projeto possui testes automatizados para validar todo o comportamento:

Para rodar os testes:

```
make test
```

O comando acima irá:

- Rodar todos os testes unitários.

- Mostrar os logs -v (verbose) de cada requisição e resposta.

Se tudo estiver correto, o resultado será:

- Todas as validações de rate limit passando

- Logs indicando 200 para requisições válidas e 429 para bloqueios