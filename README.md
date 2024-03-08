# go-stresstest

## Descrição do Desafio:

**Objetivo:** Criar um sistema CLI em Go para realizar testes de carga em um serviço web.

O usuário deve fornecer os seguintes parâmetros via CLI:

- `--url`: URL do serviço a ser testado.
- `--requests`: Número total de requisições.
- `--concurrency`: Número de chamadas simultâneas.

### Execução do Teste:

- Realizar requisições HTTP para a URL especificada.
- Distribuir as requisições de acordo com o nível de concorrência definido.
- Garantir que o número total de requisições seja atendido.

### Geração de Relatório:

Ao final dos testes, apresentar um relatório contendo:

- Tempo total gasto na execução.
- Quantidade total de requisições realizadas.
- Quantidade de requisições com status HTTP 200.
- Distribuição de outros códigos de status HTTP (como 404, 500, etc.).

---

## Documentação do Projeto

### Introdução
Este projeto é um CLI para realizar testes de carga em um serviço web.

### Como executar?

#### Ambiente Dev
Altere o arquivo .env com os seguintes valores:

```
DOCKERFILE=Dockerfile.dev
IS_DEV=true
```

Execute o seguinte comando através do Docker Compose:

```shell
$ docker compose up -d
```

Conecte-se no container **stresstest** e execute o serviço:

```shell
$ docker compose exec stresstest sh
$ go run main.go stress -u http://www.google.com -c 5 -r 100
```

Dica: Utilize a extensão **Remote Development** no **VSCode** para realizar um ´Attach to running container´.

### Ambiente Produção
Altere o arquivo .env com os seguintes valores:

```
DOCKERFILE=Dockerfile.prod
IS_DEV=false
```

Execute o seguinte comando através do Docker Compose:

```shell
$ docker compose up --build
```

O contêiner **stresstest** estará pronto para uso.

#### Como testar?

```shell
$ ./go-stresstest stress -u http://www.google.com -c 5 -r 100
```

### Testes unitários

Foi utilizado o pacote [gotestsum](https://github.com/gotestyourself/gotestsum)

```shell
$ gotestsum --format=short -- -coverprofile=coverage.out ./...
$ go tool cover -html=coverage.out -o coverage.html
```

Após a geração, abra o arquivo coverage.html para verificar a cobertura.
