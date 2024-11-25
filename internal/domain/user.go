package domain

type UserId uint

type User struct {
	ID           UserId
	Login        string
	PasswordHash string
	IsModerator  bool
}
