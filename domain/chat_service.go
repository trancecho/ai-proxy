package domain

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/trancecho/ai-proxy/config"
	"github.com/trancecho/ai-proxy/domain/dto"
	"github.com/trancecho/ai-proxy/pkg/utils"
	"github.com/trancecho/ai-proxy/po"
	"github.com/trancecho/ai-proxy/po/repository"
	"io"
	"net/http"
	"strings"
	"time"
)

type AIServicePo struct {
	repo repository.ChatRepository
}

func NewAIService(repo repository.ChatRepository) *AIServicePo {
	return &AIServicePo{
		repo: repo,
	}
}

// 获取聊天历史
func (s *AIServicePo) GetChatHistory(userID uint) ([]po.RequestLog, error) {
	return s.repo.GetUserChatHistory(userID)
}

// 调用大模型API
func (s *AIServicePo) CallAIProxy(userID uint, req dto.ChatRequest) (string, error) {
	apiURL := config.GetAPIURL()
	apiKey := config.GetAPIKey()

	// 组装请求
	reqData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request failed: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqData))
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	resp, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// 如果是流式响应，使用逐行解析
	if req.Stream {
		return s.handleStreamResponse(resp.Body)
	}

	// 普通 JSON 响应，直接解析
	var aiResp dto.ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return "", fmt.Errorf("decode response failed: %w", err)
	}

	// 只返回 Choices 中的第一个内容
	if len(aiResp.Choices) > 0 {
		// 处理可能含有 Markdown 的内容
		processor := utils.MarkdownProcessor{}
		return processor.Do(aiResp.Choices[0].Message.Content), nil
	}

	return "", fmt.Errorf("no content found in response")
}

// 处理流式响应
func (s *AIServicePo) handleStreamResponse(body io.ReadCloser) (string, error) {
	defer body.Close()

	reader := bufio.NewReader(body)
	var responseContent string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("error reading stream: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" || line == "data: [DONE]" {
			continue
		}

		// 移除 "data: " 前缀
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		}

		// 解析 JSON
		var chunk dto.ChatResponse
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			fmt.Printf("Failed to parse JSON chunk: %s\n", line) // 调试日志
			continue
		}

		// 拼接 Message.Content
		for _, choice := range chunk.Choices {
			responseContent += choice.Message.Content // 修正这里，直接使用 `choice.Message.Content`
		}
	}

	// 返回拼接的内容
	return responseContent, nil
}

//bug 1 没有保存到历史对话中去
