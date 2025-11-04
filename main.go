package main

import (
	"log"
	"os"

	"github.com/Arterning/go-rag/database"
	"github.com/Arterning/go-rag/handlers"
	"github.com/Arterning/go-rag/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 0. 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用默认配置")
	}

	// 1. 初始化数据库
	log.Println("初始化数据库...")
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 2. 初始化 LLM 客户端
	log.Println("初始化 LLM 客户端...")
	if err := services.InitLLM(); err != nil {
		log.Fatalf("LLM 初始化失败: %v\n提示：请设置 ANTHROPIC_API_KEY 环境变量", err)
	}

	// 3. 创建 Gin 路由
	r := gin.Default()

	// 设置最大上传文件大小（默认 32MB）
	r.MaxMultipartMemory = 32 << 20

	// 4. 配置路由
	api := r.Group("/api")
	{
		// 文件上传接口
		api.POST("/upload", handlers.UploadDocxHandler)

		// AI 问答接口
		api.POST("/qa", handlers.QAHandler)

		// 健康检查接口
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"message": "RAG 服务运行正常",
			})
		})
	}

	// 5. 获取端口配置
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 6. 启动服务器
	log.Printf("服务器启动在端口 %s...", port)
	log.Println("API 端点:")
	log.Println("  POST /api/upload  - 上传 docx 文件")
	log.Println("  POST /api/qa      - AI 问答")
	log.Println("  GET  /api/health  - 健康检查")

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
