package domain

type UserId uint

type User struct {
	ID           UserId
	Name         string
	Email        string
	PasswordHash string
	IsModerator  bool
}
