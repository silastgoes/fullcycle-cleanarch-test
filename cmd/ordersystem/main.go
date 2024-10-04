package main

import (
	"database/sql"
	"fmt"

	"log"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/silastgoes/fullcycle-cleanarch-test/cmd/ordersystem/wire"
	"github.com/silastgoes/fullcycle-cleanarch-test/configs"
	"github.com/silastgoes/fullcycle-cleanarch-test/internal/event/handler"

	migrate "github.com/silastgoes/fullcycle-cleanarch-test/internal/infra/database/migrate"
	graph "github.com/silastgoes/fullcycle-cleanarch-test/internal/infra/graph/generated"
	graphResolver "github.com/silastgoes/fullcycle-cleanarch-test/internal/infra/graph/generated/resolvers"
	"github.com/silastgoes/fullcycle-cleanarch-test/internal/infra/grpc/pb"
	"github.com/silastgoes/fullcycle-cleanarch-test/internal/infra/grpc/service"
	"github.com/silastgoes/fullcycle-cleanarch-test/internal/infra/web/webserver"
	"github.com/silastgoes/fullcycle-cleanarch-test/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := migrate.NewMigrateService(db).Up(); err != nil {
		log.Fatal(fmt.Errorf("error ou migrar: %s", err.Error()))
	}

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := wire.NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUsecase := wire.NewListOrderUseCase(db)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := wire.NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler(http.MethodPost, "/order", webOrderHandler.Create)
	webserver.AddHandler(http.MethodGet, "/orders", webOrderHandler.List)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUsecase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	resolver := &graphResolver.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUsecase:   *listOrderUsecase,
	}

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
