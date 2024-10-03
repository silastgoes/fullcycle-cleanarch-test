package resolvers

import "github.com/silastgoes/fullcycle-cleanarch-test/internal/usecase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrderUsecase   usecase.ListOrderUseCase
}
