package model

type Song struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Cover    string `json:"cover"`
	Source   string `json:"source"`
	Duration int64  `json:"duration,omitempty"`
	URL      string `json:"url,omitempty"`
}

type Category struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Source string `json:"source"`
	Cover  string `json:"cover,omitempty"`
}

type Order struct {
	ID         string `json:"id"`
	SongID     string `json:"song_id"`
	SongName   string `json:"song_name"`
	Artist     string `json:"artist"`
	Cover      string `json:"cover"`
	Source     string `json:"source"`
	Message    string `json:"message"`
	SingerID   string `json:"singer_id"`
	SingerName string `json:"singer_name"`
	UserOpenID string `json:"user_open_id"`
	UserNick   string `json:"user_nick"`
	Status     string `json:"status"`
	Price      int64  `json:"price"`
	CreatedAt  int64  `json:"created_at"`
}

type Singer struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	OpenID string `json:"open_id"`
	Status string `json:"status"`
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type SearchParams struct {
	Keyword string `form:"keyword"`
	Source  string `form:"source"`
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
}

type CategorySongsParams struct {
	ToplistID string `form:"toplist_id"`
	Source    string `form:"source"`
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
}

type User struct {
	OpenID    string `json:"open_id"`
	NickName  string `json:"nick_name"`
	Avatar    string `json:"avatar"`
	Role      string `json:"role"`
	CreatedAt int64  `json:"created_at"`
}

type CreateOrderParams struct {
	SongID   string `json:"song_id"`
	SongName string `json:"song_name"`
	Artist   string `json:"artist"`
	Cover    string `json:"cover"`
	Source   string `json:"source"`
	Message  string `json:"message"`
	SingerID string `json:"singer_id"`
	Price    int64  `json:"price"`
}

type LoginParams struct {
	Code     string `json:"code"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
}

type SwitchRoleParams struct {
	Role string `json:"role"`
}
