package po

import (
	"context"
	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	GetUserByID(ctx context.Context, userID string) (*UserInfo, error)
}

// ChatRepository 对话历史仓库接口
type ChatRepository interface {
	NewChatRepository(db *gorm.DB) ChatRepository
	SaveRequestLog(log RequestLog) error
	GetUserChatHistory(userID uint) ([]RequestLog, error)
}
