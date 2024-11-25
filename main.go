package main

import (
	"PaymentAPI/entity"
	"PaymentAPI/handler"
	"PaymentAPI/repository"
	"PaymentAPI/service"
	"PaymentAPI/storage"
	"github.com/gin-gonic/gin"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	customerRepository := repository.NewCustomerRepository(storage.NewJsonFileHandler[entity.Customer]())
	walletRepository := repository.NewWalletRepository(storage.NewJsonFileHandler[entity.Wallet]())
	refreshTokenRepository := repository.NewRefreshTokenRepository(storage.NewJsonFileHandler[entity.RefreshToken]())
	blacklistRepository := repository.NewBlacklistRepository(storage.NewJsonFileHandler[entity.Blacklist]())

	walletService := service.NewWalletService(walletRepository)
	refreshTokenService := service.NewRefreshTokenService(refreshTokenRepository)
	blacklistService := service.NewBlacklistService(blacklistRepository)
	customerService := service.NewCustomerService(customerRepository, walletService)
	authService := service.NewAuthService(customerService, refreshTokenService, blacklistService)

	authHandler := handler.NewAuthHandler(authService, customerService)

	r := gin.Default()

	public := r.Group("/api/public")
	{
		public.POST("/auth/register", authHandler.HandleRegister)
		public.POST("/auth/login", authHandler.HandleLogin)
		public.POST("/auth/logout", authHandler.HandleLogout)
		public.POST("/auth/refresh-token", authHandler.HandleRefreshToken)
	}

	r.Run(":8081")
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
