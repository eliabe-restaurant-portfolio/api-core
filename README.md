# API para restaurantes

## A documentação da API

Está disponibilizada em 

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