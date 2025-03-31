package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/trancecho/ai-proxy/config"
	"github.com/trancecho/ai-proxy/domain/dto"
	"github.com/trancecho/ai-proxy/po"
	"github.com/trancecho/ai-proxy/po/repository"
	"net/http"
	"time"
)

type AIServicePo struct {
	repo repository.ChatRepository // 改为接口类型，便于测试和扩展
}

// NewAIService 接受 ChatRepository 接口作为参数
func NewAIService(repo repository.ChatRepository) *AIServicePo {
	return &AIServicePo{
		repo: repo,
	}
}

func (s *AIServicePo) CallAIProxy(userID uint, req dto.ChatRequest) (*dto.ChatResponse, error) {
	apiURL := config.GetAPIURL()
	fmt.Println("apiURL:", apiURL)
	apiKey := config.GetAPIKey() // 获取 API Key

	// 确保 Messages 不是 nil
	var messages []dto.ChatMessage
	if len(req.Messages) > 0 {
		messages = req.Messages
	}

	// 组装 OpenAI 兼容的请求结构
	reqData := dto.ChatRequest{
		Model:       req.Model,
		Stream:      req.Stream,
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}

	// JSON 序列化请求体
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	// 解析响应
	var aiResp dto.ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	// 解析返回的 AI 响应内容
	if len(aiResp.Choices) == 0 || aiResp.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("empty response from AI")
	}
	aiMessage := aiResp.Choices[0].Message.Content

	// 记录日志
	messagesJSON, _ := json.Marshal(messages) // 转换 messages 为 JSON 字符串存入数据库
	log := po.RequestLog{
		UserID:      userID,
		Model:       req.Model,
		Prompt:      string(messagesJSON),
		Response:    aiMessage,
		RequestTime: time.Now().Unix(),
	}

	// 保存请求日志
	if err := s.repo.SaveRequestLog(log); err != nil {
		fmt.Printf("failed to save log: %v\n", err)
	}

	return &aiResp, nil
}

// GetChatHistory 获取指定用户的聊天历史
func (s *AIServicePo) GetChatHistory(userID uint) ([]po.RequestLog, error) {
	return s.repo.GetUserChatHistory(userID)
}
