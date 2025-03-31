package repository

import (
	"context"
	"errors"

	"github.com/trancecho/ai-proxy/po"
	"gorm.io/gorm"
)

// UserRepositoryImpl 实现 UserRepository 接口
type UserRepositoryImpl struct {
	db *gorm.DB // ✅ 这里改成 `*gorm.DB`
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	if db == nil {
		panic("数据库连接不能为空")
	}
	return &UserRepositoryImpl{db: db}
}

// GetUserByID 通过用户ID获取用户信息
func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, userID string) (*po.UserInfo, error) {
	var user po.UserInfo
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&user) // ✅ 这里不会报错
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}
	return &user, nil
}
