package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/ai-proxy/config"
	"github.com/trancecho/ai-proxy/domain"
	"github.com/trancecho/ai-proxy/initialize/database"
	"github.com/trancecho/ai-proxy/pkg/utils"
	"github.com/trancecho/ai-proxy/po/repository"
	"log"
)

func main() {

	config.InitViper()

	db := database.InitDB() // 初始化数据库并获取数据库实例
	if db == nil {
		log.Fatal("数据库连接失败")
	}

	utils.InitSecret()

	// 2️⃣ 创建 ChatRepository 实例
	chatRepo := repository.NewChatRepository(db)
	//
	//// 3️⃣ 创建 AIService 实例并传递 ChatRepository
	aiService := domain.NewAIService(chatRepo)

	// 4️⃣ 创建 ChatHandler 实例
	// chatHandler 不需要传递 aiService，因为 NewChatHandler 已经创建了 aiService
	chatHandler := domain.NewChatHandler(aiService)
	// 7️⃣ 启动 Gin 服务器
	r := gin.Default()
	MakeRoutes(r, chatHandler)

	log.Println("✅ 服务器启动成功，监听端口：12349")
	log.Fatal(r.Run(":12349"))
}

// 注册 API 路由
func MakeRoutes(g *gin.Engine, chatHandler *domain.ChatHandler) {
	// 根路由，返回欢迎信息或健康检查
	g.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the AI Proxy API",
		})
	})

	// 发送 AI 询问
	g.POST("/api/v1/chat", utils.JWTAuthMiddleware(), chatHandler.Chat)

	// 查询历史对话
	g.GET("/api/v1/chat/history", utils.JWTAuthMiddleware(), chatHandler.GetChatHistory)
}
