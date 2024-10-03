package service

import (
	"context"

	"github.com/silastgoes/fullcycle-cleanarch-test/internal/infra/grpc/pb"
	"github.com/silastgoes/fullcycle-cleanarch-test/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrderUseCase   usecase.ListOrderUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listOrderUseCase usecase.ListOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (r *OrderService) ListOrder(ctx context.Context, _ *pb.ListOrderRequest) (*pb.ListOrderResponse, error) {
	list := make([]*pb.Order, 0)

	output, err := r.ListOrderUseCase.Execute()
	if err != nil {
		return nil, err
	}

	for _, order := range output {
		dto := &pb.Order{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		}
		list = append(list, dto)
	}

	return &pb.ListOrderResponse{Order: list}, nil
}
