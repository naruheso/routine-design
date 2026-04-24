package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/amical/routine-design/backend/internal/handler"
	"github.com/amical/routine-design/backend/internal/middleware"
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

	// Service
	authService := service.NewAuthService(userRepo)
	expenseService := service.NewExpenseService(expenseRepo, approvalRepo)
	approvalService := service.NewApprovalService(expenseRepo, approvalRepo)
	reportService := service.NewReportService(reportRepo)

	// Handler
	authHandler := handler.NewAuthHandler(authService)
	expenseHandler := handler.NewExpenseHandler(expenseService)
	approvalHandler := handler.NewApprovalHandler(approvalService)
	categoryHandler := handler.NewCategoryHandler(categoryRepo)
	reportHandler := handler.NewReportHandler(reportService)

	// Router
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		// 認証（公開）
		api.POST("/auth/login", authHandler.Login)

		// 認証必須エンドポイント
		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 経費
			expenses := authorized.Group("/expenses")
			{
				expenses.GET("", expenseHandler.List)
				expenses.POST("", expenseHandler.Create)
				expenses.GET("/:id", expenseHandler.GetByID)
				expenses.PUT("/:id", expenseHandler.Update)
				expenses.DELETE("/:id", expenseHandler.Delete)
				expenses.POST("/:id/submit", expenseHandler.Submit)
			}

			// 承認
			approvals := authorized.Group("/approvals")
			{
				approvals.GET("/pending", approvalHandler.GetPending)
				approvals.POST("/:id/approve", approvalHandler.Approve)
				approvals.POST("/:id/return", approvalHandler.Return)
				approvals.POST("/:id/reject", approvalHandler.Reject)
			}

			// 勘定科目
			authorized.GET("/categories", categoryHandler.List)

			// 集計レポート
			reports := authorized.Group("/reports")
			{
				reports.GET("/by-employee", reportHandler.ByEmployee)
				reports.GET("/by-category", reportHandler.ByCategory)
			}
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("サーバー起動", "port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("サーバー起動失敗: %v", err)
	}
}
