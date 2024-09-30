package dto

import (
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"time"
)

type AccountsFilter struct {
	From   *time.Time
	Status *domain.AccountStatus
}
