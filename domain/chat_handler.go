package domain

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/ai-proxy/domain/dto"
	"net/http"
)

type ChatHandler struct {
	service *AIServicePo
}

// 参数：AIService实例，包含所有必要的依赖
func NewChatHandler(aiService *AIServicePo) *ChatHandler {
	return &ChatHandler{
		service: aiService,
	}
}

// 处理 AI 询问,交通警察
func (h *ChatHandler) Chat(c *gin.Context) {
	var request dto.ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	// 直接从上下文获取 userID,是不是因为这里调用的时候，根本都没有调用中间件函数啊，学到了，原来需要身份验证的路由，先要调用中间件，然后再调用

	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	userID, ok := uid.(int64) // 确保 userID 是 int64
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户ID类型错误"})
		return
	}

	// 调用 AIProxy，调用函数来调用大模型
	response, err := h.service.CallAIProxy(uint(userID), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// 处理历史对话查询
func (h *ChatHandler) GetChatHistory(c *gin.Context) {
	// 直接从上下文获取 userID
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	userID, ok := uid.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户ID类型错误"})
		return
	}

	// 获取历史记录
	history, err := h.service.GetChatHistory(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
