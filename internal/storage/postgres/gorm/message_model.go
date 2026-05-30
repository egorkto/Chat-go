package storage_postgres_gorm

import (
	"time"

	"github.com/egorkto/Chat-go/internal/domain"
)

type MessageModel struct {
	ID      uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Version int       `gorm:"column:version;default:1"`
	UserID  int       `gorm:"column:user_id"`
	User    UserModel `gorm:"foreignKey:UserID;references:ID"`
	Text    string    `gorm:"column:text"`
	SentAt  time.Time `gorm:"column:sent_at"`
}

func (m MessageModel) TableName() string {
	return "chat.messages"
}

func (m MessageModel) ToDomain() domain.Message {
	return domain.NewMessage(
		int(m.ID),
		m.Version,
		m.User.ToDomain(),
		m.Text,
		m.SentAt,
	)
}
