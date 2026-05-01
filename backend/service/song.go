package service

import "jukebox-lite/model"

type SongService struct {
	musicAPI *MusicAPIService
}

func NewSongService(musicAPI *MusicAPIService) *SongService {
	return &SongService{musicAPI: musicAPI}
}

func (s *SongService) Search(keyword, source string, page, limit int) ([]*model.Song, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return s.musicAPI.Search(keyword, source, page, limit)
}

func (s *SongService) GetCategories(source string) ([]*model.Category, error) {
	return s.musicAPI.GetToplists(source)
}

func (s *SongService) GetCategorySongs(toplistID, source string, page, limit int) ([]*model.Song, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return s.musicAPI.GetToplistSongs(toplistID, source, page, limit)
}

func (s *SongService) GetSongInfo(source, id string) (*model.Song, error) {
	return s.musicAPI.GetSongInfo(source, id)
}
