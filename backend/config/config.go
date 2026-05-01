package config

type Config struct {
	Port     string
	MusicAPI MusicAPIConfig
	Wechat   WechatConfig
}

type MusicAPIConfig struct {
	BaseURL string
	Timeout int
}

type WechatConfig struct {
	AppID     string
	AppSecret string
}

var C = Config{
	Port: "8080",
	MusicAPI: MusicAPIConfig{
		BaseURL: "https://music-dl.sayqz.com/api",
		Timeout: 10,
	},
	Wechat: WechatConfig{
		AppID:     "",
		AppSecret: "",
	},
}
