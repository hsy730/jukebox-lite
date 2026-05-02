package main

import (
	"fmt"
	"jukebox-lite/config"
	"jukebox-lite/controller"
	"jukebox-lite/repository"
	"jukebox-lite/router"
	"jukebox-lite/service"
)

func main() {
	orderRepo := repository.NewMemoryOrderRepository()
	singerRepo := repository.NewMemorySingerRepository()
	userRepo := repository.NewMemoryUserRepository()

	musicAPIService := service.NewMusicAPIService()
	authService := service.NewAuthService(userRepo)
	songService := service.NewSongService(musicAPIService)
	orderService := service.NewOrderService(orderRepo, singerRepo)
	singerService := service.NewSingerService(singerRepo)

	songCtrl := controller.NewSongController(songService)
	orderCtrl := controller.NewOrderController(orderService)
	singerCtrl := controller.NewSingerController(singerService)
	authCtrl := controller.NewAuthController(authService)

	r := router.SetupRouter(songCtrl, orderCtrl, singerCtrl, authCtrl, authService)

	addr := fmt.Sprintf(":%s", config.C.Port)
	fmt.Printf("🚀 Jukebox Lite Server starting on %s\n", addr)
	if err := r.Run(addr); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
