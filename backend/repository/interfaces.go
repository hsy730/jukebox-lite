package repository

import "jukebox-lite/model"

type OrderRepository interface {
	Create(order *model.Order) error
	GetByID(id string) (*model.Order, error)
	ListBySingerID(singerID string, page, limit int) ([]*model.Order, int64, error)
	ListByUserOpenID(openID string, page, limit int) ([]*model.Order, int64, error)
	UpdateStatus(id string, status string) error
}

type SingerRepository interface {
	Create(singer *model.Singer) error
	GetByID(id string) (*model.Singer, error)
	List(page, limit int) ([]*model.Singer, int64, error)
	UpdateStatus(id string, status string) error
}
