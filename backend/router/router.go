package router

import (
	"jukebox-lite/controller"
	"jukebox-lite/middleware"
	"jukebox-lite/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	songCtrl *controller.SongController,
	orderCtrl *controller.OrderController,
	singerCtrl *controller.SingerController,
	authCtrl *controller.AuthController,
	authService *service.AuthService,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authCtrl.Login)
		}

		// 公开接口：无需登录即可浏览歌曲、歌手等信息
		publicSongs := api.Group("/songs")
		{
			publicSongs.GET("/search", songCtrl.Search)
			publicSongs.GET("/categories", songCtrl.GetCategories)
			publicSongs.GET("/category", songCtrl.GetCategorySongs)
			publicSongs.GET("/info", songCtrl.GetSongInfo)
		}

		publicSingers := api.Group("/singers")
		{
			publicSingers.GET("", singerCtrl.ListSingers)
			publicSingers.GET("/:id", singerCtrl.GetSinger)
		}

		// 需要登录的接口：点歌、查看订单、切换角色、入驻歌手等操作
		authRequired := api.Group("")
		authRequired.Use(middleware.Auth(authService))
		{
			authGroup := authRequired.Group("/auth")
			{
				authGroup.GET("/profile", authCtrl.GetProfile)
				authGroup.PUT("/role", authCtrl.SwitchRole)
			}

			orders := authRequired.Group("/orders")
			{
				orders.POST("", orderCtrl.CreateOrder)
				orders.GET("/:id", orderCtrl.GetOrder)
				orders.GET("/singer", orderCtrl.ListOrdersBySinger)
				orders.GET("/user", orderCtrl.ListOrdersByUser)
				orders.PUT("/:id/status", orderCtrl.UpdateOrderStatus)
			}

			authSingers := authRequired.Group("/singers")
			{
				authSingers.POST("", singerCtrl.CreateSinger)
			}
		}
	}

	return r
}
