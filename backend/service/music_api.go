package service

import (
	"encoding/json"
	"fmt"
	"io"
	"jukebox-lite/config"
	"jukebox-lite/model"
	"net/http"
	"net/url"
	"time"
)

type MusicAPIService struct {
	client  *http.Client
	baseURL string
}

func NewMusicAPIService() *MusicAPIService {
	return &MusicAPIService{
		client: &http.Client{
			Timeout: time.Duration(config.C.MusicAPI.Timeout) * time.Second,
		},
		baseURL: config.C.MusicAPI.BaseURL,
	}
}

type musicAPIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (s *MusicAPIService) Search(keyword, source string, page, limit int) ([]*model.Song, error) {
	params := url.Values{}
	params.Set("type", "aggregateSearch")
	params.Set("keyword", keyword)
	if source != "" {
		params.Set("source", source)
	}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("limit", fmt.Sprintf("%d", limit))

	resp, err := s.get(params)
	if err != nil {
		return s.mockSearch(keyword, page, limit), nil
	}

	var songs []*model.Song
	if err := json.Unmarshal(resp, &songs); err != nil {
		return s.parseSearchResponse(resp)
	}
	return songs, nil
}

func (s *MusicAPIService) parseSearchResponse(data json.RawMessage) ([]*model.Song, error) {
	var rawList []map[string]interface{}
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, err
	}
	var songs []*model.Song
	for _, item := range rawList {
		song := &model.Song{
			ID:     fmt.Sprintf("%v", item["id"]),
			Name:   fmt.Sprintf("%v", item["name"]),
			Artist: fmt.Sprintf("%v", item["artist"]),
			Album:  fmt.Sprintf("%v", item["album"]),
			Source: fmt.Sprintf("%v", item["source"]),
		}
		if pic, ok := item["pic"]; ok && pic != nil {
			song.Cover = fmt.Sprintf("%v", pic)
		} else if cover, ok := item["cover"]; ok && cover != nil {
			song.Cover = fmt.Sprintf("%v", cover)
		}
		if urlStr, ok := item["url"]; ok && urlStr != nil {
			song.URL = fmt.Sprintf("%v", urlStr)
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func (s *MusicAPIService) GetToplists(source string) ([]*model.Category, error) {
	params := url.Values{}
	params.Set("type", "toplists")
	if source != "" {
		params.Set("source", source)
	} else {
		params.Set("source", "netease")
	}

	resp, err := s.get(params)
	if err != nil {
		return s.mockCategories(), nil
	}

	var categories []*model.Category
	if err := json.Unmarshal(resp, &categories); err != nil {
		return s.parseToplistsResponse(resp)
	}
	return categories, nil
}

func (s *MusicAPIService) parseToplistsResponse(data json.RawMessage) ([]*model.Category, error) {
	var rawList []map[string]interface{}
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, err
	}
	var categories []*model.Category
	for _, item := range rawList {
		cat := &model.Category{
			ID:     fmt.Sprintf("%v", item["id"]),
			Name:   fmt.Sprintf("%v", item["name"]),
			Source: fmt.Sprintf("%v", item["source"]),
		}
		if cover, ok := item["cover"]; ok && cover != nil {
			cat.Cover = fmt.Sprintf("%v", cover)
		} else if pic, ok := item["pic"]; ok && pic != nil {
			cat.Cover = fmt.Sprintf("%v", pic)
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func (s *MusicAPIService) GetToplistSongs(toplistID, source string, page, limit int) ([]*model.Song, error) {
	params := url.Values{}
	params.Set("type", "toplist")
	params.Set("id", toplistID)
	if source != "" {
		params.Set("source", source)
	} else {
		params.Set("source", "netease")
	}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("limit", fmt.Sprintf("%d", limit))

	resp, err := s.get(params)
	if err != nil {
		return s.mockSearch("", page, limit), nil
	}

	return s.parseSearchResponse(resp)
}

func (s *MusicAPIService) GetSongInfo(source, id string) (*model.Song, error) {
	params := url.Values{}
	params.Set("type", "info")
	params.Set("source", source)
	params.Set("id", id)

	resp, err := s.get(params)
	if err != nil {
		return &model.Song{
			ID:     id,
			Name:   "未知歌曲",
			Artist: "未知歌手",
			Source: source,
		}, nil
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(resp, &raw); err != nil {
		return nil, err
	}

	song := &model.Song{
		ID:     id,
		Name:   fmt.Sprintf("%v", raw["name"]),
		Artist: fmt.Sprintf("%v", raw["artist"]),
		Album:  fmt.Sprintf("%v", raw["album"]),
		Source: source,
	}
	if pic, ok := raw["pic"]; ok {
		song.Cover = fmt.Sprintf("%v", pic)
	}
	if urlStr, ok := raw["url"]; ok {
		song.URL = fmt.Sprintf("%v", urlStr)
	}
	return song, nil
}

func (s *MusicAPIService) get(params url.Values) (json.RawMessage, error) {
	reqURL := fmt.Sprintf("%s/?%s", s.baseURL, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp musicAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("invalid api response: %s", string(body))
	}

	if apiResp.Code != 200 {
		return nil, fmt.Errorf("api error: %s", apiResp.Message)
	}

	return apiResp.Data, nil
}

func (s *MusicAPIService) mockCategories() []*model.Category {
	return []*model.Category{
		{ID: "cat_1", Name: "华语流行", Source: "netease"},
		{ID: "cat_2", Name: "粤语经典", Source: "netease"},
		{ID: "cat_3", Name: "摇滚", Source: "netease"},
		{ID: "cat_4", Name: "民谣", Source: "netease"},
		{ID: "cat_5", Name: "DJ/Remix", Source: "netease"},
		{ID: "cat_6", Name: "欧美流行", Source: "netease"},
		{ID: "cat_7", Name: "日韩精选", Source: "netease"},
		{ID: "cat_8", Name: "影视原声", Source: "netease"},
	}
}

func (s *MusicAPIService) mockSearch(keyword string, page, limit int) []*model.Song {
	allSongs := []*model.Song{
		{ID: "mock_1", Name: "十年", Artist: "陈奕迅", Album: "黑白灰", Cover: "https://p1.music.126.net/1vHhp2JE0RnI1jdH6Ha7hg==/109951163499530495.jpg", Source: "netease"},
		{ID: "mock_2", Name: "晴天", Artist: "周杰伦", Album: "叶惠美", Cover: "https://p1.music.126.net/1vHhp2JE0RnI1jdH6Ha7hg==/109951163499530495.jpg", Source: "netease"},
		{ID: "mock_3", Name: "光年之外", Artist: "邓紫棋", Album: "光年之外", Cover: "https://p1.music.126.net/1vHhp2JE0RnI1jdH6Ha7hg==/109951163499530495.jpg", Source: "netease"},
		{ID: "mock_4", Name: "海阔天空", Artist: "Beyond", Album: "乐与怒", Cover: "https://p1.music.126.net/1vHhp2JE0RnI1jdH6Ha7hg==/109951163499530495.jpg", Source: "netease"},
		{ID: "mock_5", Name: "平凡之路", Artist: "朴树", Album: "猎户星座", Cover: "https://p1.music.126.net/1vHhp2JE0RnI1jdH6Ha7hg==/109951163499530495.jpg", Source: "netease"},
		{ID: "mock_6", Name: "起风了", Artist: "买辣椒也用券", Album: "起风了", Cover: "https://p1.music.126.net/1vHhp2JE0RnI1jdH6Ha7hg==/109951163499530495.jpg", Source: "netease"},
		{ID: "mock_7", Name: "稻香", Artist: "周杰伦", Album: "魔杰座", Cover: "https://p1.music.126.net/1vHhp2JE0RnI1jdH6Ha7hg==/109951163499530495.jpg", Source: "netease"},
		{ID: "mock_8", Name: "后来", Artist: "刘若英", Album: "我等你", Cover: "https://p1.music.126.net/1vHhp2JE0RnI1jdH6Ha7hg==/109951163499530495.jpg", Source: "netease"},
		{ID: "mock_9", Name: "小幸运", Artist: "田馥甄", Album: "我的少女时代", Cover: "https://p1.music.126.net/1vHhp2JE0RnI1jdH6Ha7hg==/109951163499530495.jpg", Source: "netease"},
		{ID: "mock_10", Name: "匆匆那年", Artist: "王菲", Album: "匆匆那年", Cover: "https://p1.music.126.net/1vHhp2JE0RnI1jdH6Ha7hg==/109951163499530495.jpg", Source: "netease"},
		{ID: "mock_11", Name: "夜曲", Artist: "周杰伦", Album: "十一月的萧邦", Cover: "https://p1.music.126.net/1vHhp2JE0RnI1jdH6Ha7hg==/109951163499530495.jpg", Source: "netease"},
		{ID: "mock_12", Name: "浮夸", Artist: "陈奕迅", Album: "U87", Cover: "https://p1.music.126.net/1vHhp2JE0RnI1jdH6Ha7hg==/109951163499530495.jpg", Source: "netease"},
	}

	if keyword != "" {
		var filtered []*model.Song
		for _, song := range allSongs {
			if containsIgnoreCase(song.Name, keyword) || containsIgnoreCase(song.Artist, keyword) {
				filtered = append(filtered, song)
			}
		}
		return paginate(filtered, page, limit)
	}
	return paginate(allSongs, page, limit)
}

func containsIgnoreCase(s, substr string) bool {
	sLower := toLower(s)
	subLower := toLower(substr)
	for i := 0; i <= len(sLower)-len(subLower); i++ {
		if sLower[i:i+len(subLower)] == subLower {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			result[i] = byte(c + 32)
		} else {
			result[i] = byte(c)
		}
	}
	return string(result)
}

func paginate(songs []*model.Song, page, limit int) []*model.Song {
	start := (page - 1) * limit
	if start >= len(songs) {
		return []*model.Song{}
	}
	end := start + limit
	if end > len(songs) {
		end = len(songs)
	}
	return songs[start:end]
}
