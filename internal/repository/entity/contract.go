package entity

type Contract struct {
	ID          uint
	Name        string `gorm:"size:30;not null;unique"`
	Fee         *int32
	Description *string `gorm:"size:80"`
	ImageUrl    *string `gorm:"size:80"`
	Type        *string `gorm:"size:20"`

	Deleted  bool      `gorm:"not null"`
	Accounts []Account `gorm:"many2many:account_contracts"`
}
