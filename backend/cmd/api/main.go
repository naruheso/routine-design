package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/amical/routine-design/backend/internal/handler"
	"github.com/amical/routine-design/backend/internal/repository"
	"github.com/amical/routine-design/backend/internal/service"
)

func main() {
	ctx := context.Background()

	// DB接続
	pool, err := repository.NewDBPool(ctx)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer pool.Close()
	slog.Info("DB接続成功")

	// Repository
	userRepo := repository.NewUserRepository(pool)
	expenseRepo := repository.NewExpenseRepository(pool)
	approvalRepo := repository.NewApprovalRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)
	reportRepo := repository.NewReportRepository(pool)
	txManager := repository.NewTransactionManager(pool)

	// Service
	authService := service.NewAuthService(userRepo)
	expenseService := service.NewExpenseService(expenseRepo, approvalRepo, txManager)
	approvalService := service.NewApprovalService(expenseRepo, approvalRepo, txManager)
	reportService := service.NewReportService(reportRepo)

	// Handler
	authHandler := handler.NewAuthHandler(authService)
	expenseHandler := handler.NewExpenseHandler(expenseService)
	approvalHandler := handler.NewApprovalHandler(approvalService)
	categoryHandler := handler.NewCategoryHandler(categoryRepo)
	reportHandler := handler.NewReportHandler(reportService)

	// Router
	r := handler.SetupRouter(authHandler, expenseHandler, approvalHandler, categoryHandler, reportHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("サーバー起動", "port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("サーバー起動失敗: %v", err)
	}
}
