package dto

import (
	"github.com/trancecho/ai-proxy/po"
)

// ChatMessage 表示单条聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 发送到 AIProxy 的请求结构体（符合 OpenAI 兼容格式）
type ChatRequest struct {
	Model       string        `json:"model"`
	Stream      bool          `json:"stream,omitempty"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

// ChatResponse AIProxy 返回的响应（符合 OpenAI 兼容格式）
type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
	Error string `json:"error,omitempty"`
}

// RequestLog 记录 API 请求日志，并存储 Token 消耗情况
type RequestLog struct {
	po.BaseModel
	UserID      uint          `gorm:"column:user_id;index"` // 关联用户
	Model       string        `gorm:"column:model"`
	Messages    []ChatMessage `gorm:"type:json"` // 存完整对话记录
	Response    string        `gorm:"column:response"`
	RequestTime int64         `gorm:"column:request_time"`
}

// ChatHistory 记录历史消息
type ChatHistory struct {
	ChatMessage
	CreateTime int64 `json:"create_time"` // 使用时间戳格式
}
