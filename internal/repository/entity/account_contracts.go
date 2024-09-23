package entity

type AccountContracts struct {
	AccountID  uint `gorm:"primaryKey"`
	ContractID uint `gorm:"primaryKey"`
	IsMain     bool `gorm:"not null"` // todo. add fields according to requirements
}
