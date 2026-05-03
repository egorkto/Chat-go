package storage_postgres_gorm

import "github.com/egorkto/Chat-go/internal/domain"

type UserModel struct {
	ID       uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Version  int    `gorm:"column:version;default:1"`
	FullName string `gorm:"column:full_name"`
	Login    string `gorm:"column:login"`
	Password string `gorm:"column:password"`
}

func (m UserModel) TableName() string {
	return "chat.users"
}

func (m UserModel) ToDomain() domain.User {
	return domain.NewUser(
		int(m.ID),
		m.Version,
		m.FullName,
		m.Login,
	)
}
