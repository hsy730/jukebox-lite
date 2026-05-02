package controller

import (
	"jukebox-lite/model"
	"jukebox-lite/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var params model.LoginParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 400, Message: "参数错误"})
		return
	}
	if params.Code == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 400, Message: "code不能为空"})
		return
	}

	user, token, err := ctrl.authService.Login(params.Code, params.NickName, params.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "登录成功", Data: gin.H{
		"user":  user,
		"token": token,
	}})
}

func (ctrl *AuthController) SwitchRole(c *gin.Context) {
	var params model.SwitchRoleParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 400, Message: "参数错误"})
		return
	}

	openID, _ := c.Get("open_id")
	user, err := ctrl.authService.SwitchRole(openID.(string), params.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "切换成功", Data: user})
}

func (ctrl *AuthController) GetProfile(c *gin.Context) {
	openID, _ := c.Get("open_id")
	user, err := ctrl.authService.GetUserByOpenID(openID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{Code: 404, Message: "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "success", Data: user})
}
