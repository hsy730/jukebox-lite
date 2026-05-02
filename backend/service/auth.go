package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"jukebox-lite/config"
	"jukebox-lite/model"
	"jukebox-lite/repository"
	"net/http"
	"time"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

type WxCode2SessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func (s *AuthService) Login(code, nickName, avatar string) (*model.User, string, error) {
	openID, err := s.code2Session(code)
	if err != nil {
		return nil, "", fmt.Errorf("微信登录失败: %v", err)
	}

	user, err := s.userRepo.GetByOpenID(openID)
	if err != nil {
		user = &model.User{
			OpenID:    openID,
			NickName:  nickName,
			Avatar:    avatar,
			Role:      "user",
			CreatedAt: time.Now().Unix(),
		}
		if createErr := s.userRepo.Create(user); createErr != nil {
			return nil, "", createErr
		}
	} else {
		if nickName != "" {
			user.NickName = nickName
		}
		if avatar != "" {
			user.Avatar = avatar
		}
	}

	token := s.generateToken(openID)

	return user, token, nil
}

func (s *AuthService) SwitchRole(openID, role string) (*model.User, error) {
	if role != "user" && role != "singer" {
		return nil, fmt.Errorf("无效的角色: %s", role)
	}
	if err := s.userRepo.UpdateRole(openID, role); err != nil {
		return nil, err
	}
	return s.userRepo.GetByOpenID(openID)
}

func (s *AuthService) GetUserByOpenID(openID string) (*model.User, error) {
	return s.userRepo.GetByOpenID(openID)
}

func (s *AuthService) ValidateToken(token string) (string, error) {
	if len(token) < 64 {
		return "", fmt.Errorf("invalid token")
	}
	hexPart := token[:64]
	openID := token[64:]
	expectedHMAC := s.computeHMAC(openID)
	if !hmac.Equal([]byte(hexPart), []byte(expectedHMAC)) {
		return "", fmt.Errorf("invalid token")
	}
	return openID, nil
}

func (s *AuthService) code2Session(code string) (string, error) {
	if config.C.Wechat.AppID == "" || config.C.Wechat.AppSecret == "" {
		return s.devCode2Session(code)
	}

	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		config.C.Wechat.AppID,
		config.C.Wechat.AppSecret,
		code,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求微信API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var result WxCode2SessionResp
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信API错误: %d %s", result.ErrCode, result.ErrMsg)
	}

	return result.OpenID, nil
}

func (s *AuthService) devCode2Session(code string) (string, error) {
	if code == "" {
		return "", fmt.Errorf("code不能为空")
	}
	return "dev_" + code, nil
}

func (s *AuthService) generateToken(openID string) string {
	hmacStr := s.computeHMAC(openID)
	return hmacStr + openID
}

func (s *AuthService) computeHMAC(openID string) string {
	key := "jukebox-lite-secret-key"
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(openID))
	return hex.EncodeToString(mac.Sum(nil))
}
