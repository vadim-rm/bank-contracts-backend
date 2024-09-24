package entity

import "github.com/vadim-rm/bmstu-web-backend/internal/domain"

type Contract struct {
	ID          uint
	Name        string  `gorm:"size:60;not null;unique"`
	Fee         *int32  `gorm:"not null"`
	Description *string `gorm:"size:80"`
	ImageUrl    *string `gorm:"size:80"`
	Type        *string `gorm:"size:20;not null"`

	Deleted  bool      `gorm:"not null"`
	Accounts []Account `gorm:"many2many:account_contracts"`
}

func (c Contract) ToDomain() domain.Contract {
	return domain.Contract{
		Id:          domain.ContractId(c.ID),
		Name:        c.Name,
		Fee:         c.Fee,
		Description: c.Description,
		ImageUrl:    c.ImageUrl,
		Type:        (*domain.ContractType)(c.Type),
	}
}

func (c Contract) ToAccountDomain() domain.AccountContract {
	return domain.AccountContract{
		Id:          domain.ContractId(c.ID),
		Name:        c.Name,
		Fee:         c.Fee,
		Description: c.Description,
		ImageUrl:    c.ImageUrl,
		Type:        (*domain.ContractType)(c.Type),
	}
}
