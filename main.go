package main

import (
	"PaymentAPI/config"
	"PaymentAPI/entity"
	"PaymentAPI/handler"
	"PaymentAPI/middleware"
	"PaymentAPI/repository"
	"PaymentAPI/service"
	"PaymentAPI/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	customerRepository := repository.NewCustomerRepository(storage.NewJsonFileHandler[entity.Customer]())
	walletRepository := repository.NewWalletRepository(storage.NewJsonFileHandler[entity.Wallet]())
	refreshTokenRepository := repository.NewRefreshTokenRepository(storage.NewJsonFileHandler[entity.RefreshToken]())
	blacklistRepository := repository.NewBlacklistRepository(storage.NewJsonFileHandler[entity.Blacklist]())
	transactionRepository := repository.NewTransactionRepository(storage.NewJsonFileHandler[entity.Transaction]())

	walletService := service.NewWalletService(walletRepository)
	refreshTokenService := service.NewRefreshTokenService(refreshTokenRepository)
	blacklistService := service.NewBlacklistService(blacklistRepository)
	customerService := service.NewCustomerService(customerRepository, walletService)
	authService := service.NewAuthService(customerService, refreshTokenService, blacklistService)
	transactionService := service.NewTransactionService(transactionRepository, walletService)

	authHandler := handler.NewAuthHandler(authService, customerService)
	transactionHandler := handler.NewTransactionHandler(transactionService, walletService)
	customerHandler := handler.NewCustomerHandler(customerService)

	r := gin.Default()

	public := r.Group("/api/public")
	{
		public.POST("/auth/register", authHandler.HandleRegister)
		public.POST("/auth/login", authHandler.HandleLogin)
		public.POST("/auth/logout", authHandler.HandleLogout)
		public.POST("/auth/refresh-token", authHandler.HandleRefreshToken)
	}

	r.Use(middleware.AuthMiddleware(blacklistService))

	transaction := r.Group("/api/transactions")
	{
		transaction.POST("", transactionHandler.HandleCreateTransaction)
	}

	customer := r.Group("/api/customers")
	{
		customer.GET("/:id", customerHandler.HandleGetCustomerById)
	}

	err := r.Run(":" + config.ServerPort)
	if err != nil {
		return
	}
}
