package controller

import (
	"jukebox-lite/model"
	"jukebox-lite/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService *service.OrderService
}

func NewOrderController(orderService *service.OrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

func (ctrl *OrderController) CreateOrder(c *gin.Context) {
	var params model.CreateOrderParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 400, Message: "参数错误"})
		return
	}
	if params.SongID == "" || params.SongName == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 400, Message: "歌曲信息不完整"})
		return
	}
	if params.SingerID == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 400, Message: "请选择歌手"})
		return
	}

	userOpenID := c.GetHeader("X-User-OpenID")
	userNick := c.GetHeader("X-User-Nick")

	order, err := ctrl.orderService.CreateOrder(&params, userOpenID, userNick)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "下单成功", Data: order})
}

func (ctrl *OrderController) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := ctrl.orderService.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{Code: 404, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "success", Data: order})
}

func (ctrl *OrderController) ListOrdersBySinger(c *gin.Context) {
	singerID := c.Query("singer_id")
	if singerID == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 400, Message: "singer_id不能为空"})
		return
	}
	page := getIntParam(c, "page", 1)
	limit := getIntParam(c, "limit", 20)

	orders, total, err := ctrl.orderService.ListOrdersBySinger(singerID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 500, Message: err.Error()})
		return
	}
	if orders == nil {
		orders = []*model.Order{}
	}
	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "success", Data: gin.H{
		"list":  orders,
		"total": total,
		"page":  page,
		"limit": limit,
	}})
}

func (ctrl *OrderController) ListOrdersByUser(c *gin.Context) {
	openID := c.GetHeader("X-User-OpenID")
	if openID == "" {
		openID = c.Query("open_id")
	}
	page := getIntParam(c, "page", 1)
	limit := getIntParam(c, "limit", 20)

	orders, total, err := ctrl.orderService.ListOrdersByUser(openID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 500, Message: err.Error()})
		return
	}
	if orders == nil {
		orders = []*model.Order{}
	}
	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "success", Data: gin.H{
		"list":  orders,
		"total": total,
		"page":  page,
		"limit": limit,
	}})
}

func (ctrl *OrderController) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{Code: 400, Message: "参数错误"})
		return
	}
	if err := ctrl.orderService.UpdateOrderStatus(id, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{Code: 500, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "更新成功"})
}

func getIntParam(c *gin.Context, key string, defaultVal int) int {
	val := defaultVal
	if str := c.Query(key); str != "" {
		if v, err := strconv.Atoi(str); err == nil {
			val = v
		}
	}
	return val
}
