package po

// UserInfo 用户信息模型（主项目已有，不需要重新创建）
// 这里假设主项目的 User 表已经存在，我们只是在代码中映射它
type UserInfo struct {
	BaseModel
	UserID     int64  `gorm:"column:user_id;primaryKey"` // 这里不建议使用 `int64` 作为 ID，建议 UUID 或字符串
	Username   string `gorm:"column:username;not null"`
	Points     int64  `gorm:"column:points;not null;default:0"`     // 账户积分
	Experience int64  `gorm:"column:experience;not null;default:0"` // 经验值
	Level      int    `gorm:"column:level;not null;default:1"`      // 用户等级
}

// RequestLog 记录 API 请求日志，并存储 Token 消耗情况
type RequestLog struct {
	BaseModel
	UserID      uint   `gorm:"column:user_id;index"` // 关联用户
	Model       string `gorm:"column:model"`
	Prompt      string `gorm:"column:prompt"`
	Response    string `gorm:"column:response"`
	RequestTime int64  `gorm:"column:request_time"`
}

// ChatHistory 记录历史消息
//type ChatHistory struct {
//	ChatMessage
//	CreateTime time.Time `json:"create_time"`
//}

// ChatMessage 单条聊天消息
//type ChatMessage struct {
//	Role    RoleType `json:"role"`
//	Content string   `json:"content"`
//}

// RoleType 角色类型
//type RoleType string

// Model类型，为用户提供更多选择。
//type Model struct {
//	Model string `json:"model"`
//}
