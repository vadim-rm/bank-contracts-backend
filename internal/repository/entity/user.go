package entity

import "github.com/vadim-rm/bmstu-web-backend/internal/domain"

type User struct {
	ID           uint
	Login        string `gorm:"size:80;not null;unique"`
	PasswordHash string `gorm:"size:60;not null"`
	IsModerator  bool   `gorm:"not null"`

	CreatorIn   []Account `gorm:"foreignKey:Creator"`
	ModeratorIn []Account `gorm:"foreignKey:Moderator"`
}

func (u User) ToDomain() domain.User {
	return domain.User{
		ID:           domain.UserId(u.ID),
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
		IsModerator:  u.IsModerator,
	}
}
