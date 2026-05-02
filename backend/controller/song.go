package controller

import (
	"jukebox-lite/model"
	"jukebox-lite/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SongController struct {
	songService *service.SongService
}

func NewSongController(songService *service.SongService) *SongController {
	return &SongController{songService: songService}
}

func (ctrl *SongController) Search(c *gin.Context) {
	var params model.SearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 400, Message: "参数错误"})
		return
	}
	if params.Keyword == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 400, Message: "搜索关键词不能为空"})
		return
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}

	songs, err := ctrl.songService.Search(params.Keyword, params.Source, params.Page, params.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 500, Message: err.Error()})
		return
	}
	if songs == nil {
		songs = []*model.Song{}
	}
	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "success", Data: songs})
}

func (ctrl *SongController) GetCategories(c *gin.Context) {
	source := c.Query("source")
	categories, err := ctrl.songService.GetCategories(source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 500, Message: err.Error()})
		return
	}
	if categories == nil {
		categories = []*model.Category{}
	}
	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "success", Data: categories})
}

func (ctrl *SongController) GetCategorySongs(c *gin.Context) {
	var params model.CategorySongsParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 400, Message: "参数错误"})
		return
	}
	// 如果 toplist_id 为空，则使用默认推荐（空字符串表示获取推荐歌曲）
	if params.ToplistID == "" {
		// 返回空列表或调用专门的推荐接口
		songs, err := ctrl.songService.GetRecommendSongs(params.Source, params.Page, params.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 500, Message: err.Error()})
			return
		}
		if songs == nil {
			songs = []*model.Song{}
		}
		c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "success", Data: songs})
		return
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}

	songs, err := ctrl.songService.GetCategorySongs(params.ToplistID, params.Source, params.Page, params.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 500, Message: err.Error()})
		return
	}
	if songs == nil {
		songs = []*model.Song{}
	}
	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "success", Data: songs})
}

func (ctrl *SongController) GetSongInfo(c *gin.Context) {
	source := c.Query("source")
	id := c.Query("id")
	if source == "" || id == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 400, Message: "source和id不能为空"})
		return
	}

	song, err := ctrl.songService.GetSongInfo(source, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 500, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "success", Data: song})
}
