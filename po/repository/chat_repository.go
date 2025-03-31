package repository

import (
	"github.com/trancecho/ai-proxy/po"
	"gorm.io/gorm"
)

// ChatRepository 结构体
type ChatRepository struct {
	db *gorm.DB
}

// NewChatRepository 返回 ChatRepository 实例
func NewChatRepository(db *gorm.DB) ChatRepository {
	if db == nil {
		panic("数据库连接不能为空")
	}
	return ChatRepository{db: db}
}

// 实现 ChatRepository 接口的方法
func (r *ChatRepository) SaveRequestLog(log po.RequestLog) error {
	return r.db.Create(&log).Error
}

func (r *ChatRepository) GetUserChatHistory(userID uint) ([]po.RequestLog, error) {
	var logs []po.RequestLog
	err := r.db.Where("user_id = ?", userID).Order("request_time DESC").Find(&logs).Error
	return logs, err
}
