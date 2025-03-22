package dtos

// MqDTO 消息body
type MqDTO struct {
	Topic   string `json:"topic" binding:"required"`   // 账户 ID
	Message string `json:"message" binding:"required"` // 账户号（必填）
}
