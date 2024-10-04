# fullcycle-cleanarch-test
Desafio full cycle Clear Architecture

## Funcionalidades

- **REST API**: Endpoints para criar e listar pedidos.
- **gRPC API**: Serviços gRPC para criar e listar pedidos.
- **GraphQL API**: Interface para criar e listar pedidos via GraphQL.
- **Eventos com RabbitMQ**: Eventos de criação de pedido são despachados através do RabbitMQ.
- **Banco de Dados MySQL**: Persistência dos pedidos com suporte a migrações automáticas.
  
## Tecnologias Utilizadas

- Go
- gRPC
- GraphQL
- RabbitMQ
- MySQL
- gqlgen (para GraphQL)
- `sql-driver/mysql`
- `google.golang.org/grpc`
- `github.com/streadway/amqp`

## Pré-requisitos

- Go 1.18+
- Docker (para RabbitMQ e MySQL)
- [gqlgen](https://gqlgen.com/getting-started/)

## Configuração

1. **Clone o repositório:**

```bash
git clone https://github.com/seu-usuario/seu-repositorio.git
cd seu-repositorio
```

2. **Configuração do Banco de Dados e RabbitMQ:**

Certifique-se de que o MySQL e o RabbitMQ estejam rodando. Você pode usar Docker para isso:

```bash
docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=root -p 3306:3306 -d mysql:latest
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

ou usar o docker compode com comando:

```bash
docker_compose run
```

3. **Configuração do arquivo `.env`:**

Crie o arquivo `.env` com as seguintes configurações:

```
DB_DRIVER=mysql
DB_USER=root
DB_PASSWORD=root
DB_HOST=localhost
DB_PORT=3306
DB_NAME=order_system
WEB_SERVER_PORT=8080
GRPC_SERVER_PORT=50051
GRAPHQL_SERVER_PORT=8081
```
4. **Executar as migrações:**

O serviço de migrações para criar as tabelas no banco de dados é executado na execução da aplicação

## Executar a Aplicação

Após configurar o ambiente, você pode rodar a aplicação com:

```bash
go run main.go
```

A aplicação iniciará três servidores:

* REST API: Rodando na porta 8080 para interagir via HTTP.
* gRPC Server: Rodando na porta 50051 para interações gRPC.
* GraphQL Server: Rodando na porta 8081 para interações GraphQL.

## Endpoints REST

* POST /order: Cria um novo pedido.
* GET /orders: Lista todos os pedidos.

## GraphQL Playground

Acesse o playground GraphQL em http://localhost:8081 para testar as queries.

## Adicionando um Novo Método Compartilhado (REST, gRPC e GraphQL)

Para adicionar um novo método que será compartilhado entre as três APIs (REST, gRPC e GraphQL), siga os passos abaixo:

1. Crie o Caso de Uso (Use Case)
No diretório internal, crie o novo caso de uso. Suponha que você queira adicionar um método para atualizar pedidos:

```go
package usecase

import "database/sql"

type UpdateOrderUseCase struct {
    DB *sql.DB
}

func NewUpdateOrderUseCase(db *sql.DB) *UpdateOrderUseCase {
    return &UpdateOrderUseCase{DB: db}
}

func (u *UpdateOrderUseCase) Execute(orderID string, data OrderData) error {
    // Lógica para atualizar o pedido
}
```

2. Adicione o Método na API REST

No handler da API REST, adicione um novo endpoint que usa o novo caso de uso:

```go
func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
    // Lógica para chamar o caso de uso de atualização
}
```

E no servidor web:

```go
webserver.AddHandler(http.MethodPut, "/order/{id}", webOrderHandler.Update)
```

3. Adicione o Método na API gRPC

No serviço gRPC, adicione o método ao proto e implemente o serviço:

```proto
service OrderService {
    rpc UpdateOrder(UpdateOrderRequest) returns (UpdateOrderResponse);
}
```

Implemente o método:

```go
func (s *OrderService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
    // Lógica para chamar o caso de uso de atualização
}
```

4. Adicione o Método na API GraphQL

No resolver GraphQL, adicione a nova mutation:

```go
func (r *mutationResolver) UpdateOrder(ctx context.Context, input UpdateOrderInput) (*Order, error) {
    // Lógica para chamar o caso de uso de atualização
}
```

Atualize o schema GraphQL:

```go
type Mutation {
    updateOrder(input: UpdateOrderInput!): Order!
}
```
