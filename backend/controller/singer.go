package controller

import (
	"jukebox-lite/model"
	"jukebox-lite/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SingerController struct {
	singerService *service.SingerService
}

func NewSingerController(singerService *service.SingerService) *SingerController {
	return &SingerController{singerService: singerService}
}

func (ctrl *SingerController) ListSingers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	singers, total, err := ctrl.singerService.ListSingers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 500, Message: err.Error()})
		return
	}
	if singers == nil {
		singers = []*model.Singer{}
	}
	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "success", Data: gin.H{
		"list":  singers,
		"total": total,
		"page":  page,
		"limit": limit,
	}})
}

func (ctrl *SingerController) GetSinger(c *gin.Context) {
	id := c.Param("id")
	singer, err := ctrl.singerService.GetSinger(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{Code: 404, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "success", Data: singer})
}

func (ctrl *SingerController) CreateSinger(c *gin.Context) {
	var singer model.Singer
	if err := c.ShouldBindJSON(&singer); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 400, Message: "参数错误"})
		return
	}
	if err := ctrl.singerService.CreateSinger(&singer); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 500, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "创建成功", Data: singer})
}
