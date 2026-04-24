package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/amical/routine-design/backend/internal/middleware"
)

// SetupRouter はアプリケーションの全てのエンドポイントを登録したGinルーターを生成して返します。
func SetupRouter(
	authHandler *AuthHandler,
	expenseHandler *ExpenseHandler,
	approvalHandler *ApprovalHandler,
	categoryHandler *CategoryHandler,
	reportHandler *ReportHandler,
) *gin.Engine {
	r := gin.Default()
	
	// CORS設定
	r.Use(middleware.CORSMiddleware())

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		// 認証不要の公開API
		api.POST("/auth/login", authHandler.Login)

		// 認証必須のAPI群
		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 経費関連
			expenses := authorized.Group("/expenses")
			{
				expenses.GET("", expenseHandler.List)
				expenses.POST("", expenseHandler.Create)
				expenses.GET("/:id", expenseHandler.GetByID)
				expenses.PUT("/:id", expenseHandler.Update)
				expenses.POST("/:id/submit", expenseHandler.Submit)
				expenses.DELETE("/:id", expenseHandler.Delete)
			}

			// 承認関連
			approvals := authorized.Group("/approvals")
			{
				approvals.GET("/pending", approvalHandler.GetPending)
				approvals.POST("/:id/approve", approvalHandler.Approve)
				approvals.POST("/:id/return", approvalHandler.Return)
				approvals.POST("/:id/reject", approvalHandler.Reject)
			}

			// カテゴリ関連
			authorized.GET("/categories", categoryHandler.List)

			// 集計・レポート関連
			reports := authorized.Group("/reports")
			{
				reports.GET("/by-employee", reportHandler.ByEmployee)
				reports.GET("/by-category", reportHandler.ByCategory)
			}
		}
	}

	return r
}
