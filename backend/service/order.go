package service

import (
	"jukebox-lite/model"
	"jukebox-lite/repository"
)

type OrderService struct {
	orderRepo  repository.OrderRepository
	singerRepo repository.SingerRepository
}

func NewOrderService(orderRepo repository.OrderRepository, singerRepo repository.SingerRepository) *OrderService {
	return &OrderService{
		orderRepo:  orderRepo,
		singerRepo: singerRepo,
	}
}

func (s *OrderService) CreateOrder(params *model.CreateOrderParams, userOpenID, userNick string) (*model.Order, error) {
	singer, err := s.singerRepo.GetByID(params.SingerID)
	if err != nil {
		return nil, err
	}

	order := &model.Order{
		SongID:     params.SongID,
		SongName:   params.SongName,
		Artist:     params.Artist,
		Cover:      params.Cover,
		Source:     params.Source,
		Message:    params.Message,
		SingerID:   params.SingerID,
		SingerName: singer.Name,
		UserOpenID: userOpenID,
		UserNick:   userNick,
		Price:      params.Price,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrder(id string) (*model.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) ListOrdersBySinger(singerID string, page, limit int) ([]*model.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return s.orderRepo.ListBySingerID(singerID, page, limit)
}

func (s *OrderService) ListOrdersByUser(openID string, page, limit int) ([]*model.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return s.orderRepo.ListByUserOpenID(openID, page, limit)
}

func (s *OrderService) UpdateOrderStatus(id, status string) error {
	return s.orderRepo.UpdateStatus(id, status)
}
