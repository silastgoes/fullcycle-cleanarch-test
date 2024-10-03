//go:build wireinject
// +build wireinject

package wire

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/silastgoes/fullcycle-cleanarch-test/internal/entity"
	"github.com/silastgoes/fullcycle-cleanarch-test/internal/event"
	"github.com/silastgoes/fullcycle-cleanarch-test/internal/infra/database"
	"github.com/silastgoes/fullcycle-cleanarch-test/internal/infra/web"
	"github.com/silastgoes/fullcycle-cleanarch-test/internal/usecase"
	"github.com/silastgoes/fullcycle-cleanarch-test/pkg/events"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewListOrderUseCase(db *sql.DB) *usecase.ListOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewListOrderUseCase,
	)
	return &usecase.ListOrderUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
