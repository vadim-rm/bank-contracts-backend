package domain

type UserId uint

type User struct {
	ID           UserId
	Email        string
	PasswordHash string
	IsModerator  bool
}
