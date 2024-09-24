package entity

type AccountContracts struct {
	ID         uint
	AccountID  uint `gorm:"index:idx_composite_key,unique"`
	ContractID uint `gorm:"index:idx_composite_key,unique"`
	IsMain     bool `gorm:"not null"`
}
