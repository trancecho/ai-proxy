package dto

import (
	"github.com/trancecho/ai-proxy/po"
)

// 定义角色类型的枚举
type RoleType string

const (
	RoleSystem    RoleType = "system"
	RoleUser      RoleType = "user"
	RoleAssistant RoleType = "assistant"
)

// ChatMessage 表示单条聊天消息
type ChatMessage struct {
	Role    RoleType `json:"role"`    // 使用 RoleType 类型确保角色的合法性
	Content string   `json:"content"` // 消息内容
}

// ChatRequest 发送到 AIProxy 的请求结构体（符合 OpenAI 兼容格式）
type ChatRequest struct {
	Model       string        `json:"model"`                 // 请求的模型
	Stream      bool          `json:"stream,omitempty"`      // 是否启用流式响应
	Messages    []ChatMessage `json:"messages"`              // 消息数组
	MaxTokens   int           `json:"max_tokens,omitempty"`  // 最大 tokens 数量
	Temperature float64       `json:"temperature,omitempty"` // 温度，控制生成文本的随机性
}

// ChatResponse AIProxy 返回的完整响应结构（包含 index、message、finish_reason）
type ChatResponse struct {
	Model   string `json:"model"` // 返回的模型
	Choices []struct {
		Index        int         `json:"index"`         // 选择的索引
		Message      ChatMessage `json:"message"`       // 包含消息内容
		FinishReason string      `json:"finish_reason"` // 完成原因
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`     // 提示tokens数量
		CompletionTokens int `json:"completion_tokens"` // 生成tokens数量
		TotalTokens      int `json:"total_tokens"`      // 总tokens数量
	} `json:"usage"`
}

// RequestLog 记录 API 请求日志，并存储 Token 消耗情况
type RequestLog struct {
	po.BaseModel
	UserID      uint          `gorm:"column:user_id;index"` // 关联用户
	Model       string        `gorm:"column:model"`         // 使用的模型
	Messages    []ChatMessage `gorm:"type:json"`            // 存完整对话记录
	Response    string        `gorm:"column:response"`      // 响应内容
	RequestTime int64         `gorm:"column:request_time"`  // 请求时间
}

// ChatHistory 记录历史消息
type ChatHistory struct {
	ChatMessage
	CreateTime int64 `json:"create_time"` // 使用时间戳格式
}
