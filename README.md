# API para restaurantes

## A documentação da API

Está disponibilizada em 

## Principais comandos

``` shell
make build # constroi a iamge dos containers

make up # docker compose up -d

make down # docker compose down

make logs # verifica os logs do docker_restaurant_api

make migrate-up # faz up em todos os migratios

make migrate-down # faz down em no último migration
```

## Como executar?

1. Usando Linux, instale o Make.

2. Faça o comando abaixo para gerar as chaves necessárias para a autenticação

``` shell
make generate-key
```

3. Adicione o .env seguindo example.env como referência.

4. Faça o build:

``` shell
make build
```

5. Rode os migrations, caso necessaŕio:

``` shell
make migrate-up
```