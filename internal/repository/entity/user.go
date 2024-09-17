package entity

type User struct {
	ID           uint
	Email        string `gorm:"size:80;not null;unique"`
	PasswordHash string `gorm:"size:60;not null"`
	IsModerator  bool   `gorm:"not null"`

	CreatorIn   []Account `gorm:"foreignKey:Creator"`
	ModeratorIn []Account `gorm:"foreignKey:Moderator"`
}
