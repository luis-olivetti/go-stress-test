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
