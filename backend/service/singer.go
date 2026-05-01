package service

import (
	"jukebox-lite/model"
	"jukebox-lite/repository"
)

type SingerService struct {
	singerRepo repository.SingerRepository
}

func NewSingerService(singerRepo repository.SingerRepository) *SingerService {
	return &SingerService{singerRepo: singerRepo}
}

func (s *SingerService) GetSinger(id string) (*model.Singer, error) {
	return s.singerRepo.GetByID(id)
}

func (s *SingerService) ListSingers(page, limit int) ([]*model.Singer, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return s.singerRepo.List(page, limit)
}

func (s *SingerService) CreateSinger(singer *model.Singer) error {
	return s.singerRepo.Create(singer)
}

func (s *SingerService) UpdateSingerStatus(id, status string) error {
	return s.singerRepo.UpdateStatus(id, status)
}
