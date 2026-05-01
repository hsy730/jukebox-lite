package router

import (
	"jukebox-lite/controller"
	"jukebox-lite/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	songCtrl *controller.SongController,
	orderCtrl *controller.OrderController,
	singerCtrl *controller.SingerController,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		songs := api.Group("/songs")
		{
			songs.GET("/search", songCtrl.Search)
			songs.GET("/categories", songCtrl.GetCategories)
			songs.GET("/category", songCtrl.GetCategorySongs)
			songs.GET("/info", songCtrl.GetSongInfo)
		}

		orders := api.Group("/orders")
		{
			orders.POST("", orderCtrl.CreateOrder)
			orders.GET("/:id", orderCtrl.GetOrder)
			orders.GET("/singer", orderCtrl.ListOrdersBySinger)
			orders.GET("/user", orderCtrl.ListOrdersByUser)
			orders.PUT("/:id/status", orderCtrl.UpdateOrderStatus)
		}

		singers := api.Group("/singers")
		{
			singers.GET("", singerCtrl.ListSingers)
			singers.GET("/:id", singerCtrl.GetSinger)
			singers.POST("", singerCtrl.CreateSinger)
		}
	}

	return r
}
