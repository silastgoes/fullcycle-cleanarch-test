package usecase

import "github.com/silastgoes/fullcycle-cleanarch-test/internal/entity"

type ListOrderOutputDTO []*OrderOutputDTO

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

func (l *ListOrderUseCase) Execute() (ListOrderOutputDTO, error) {
	list := make(ListOrderOutputDTO, 0)

	result, err := l.OrderRepository.List()
	if err != nil {
		return nil, err
	}

	for _, order := range result {
		dto := &OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		list = append(list, dto)
	}

	return list, nil
}
