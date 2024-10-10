package model

import "github.com/carp-cobain/referrals/domain"

// Campaign represents a named referral campaign for a blockchain address.
type Campaign struct {
	ID        uint64 `gorm:"primarykey"`
	Address   string `gorm:"index;not null"`
	Name      string
	CreatedAt Time
	UpdatedAt Time
}

// ToDomain converts a model to a domain object representation.
func (self Campaign) ToDomain() domain.Campaign {
	return domain.Campaign{
		ID:        self.ID,
		Address:   self.Address,
		Name:      self.Name,
		CreatedAt: self.CreatedAt.FromUnix(),
		UpdatedAt: self.UpdatedAt.FromUnix(),
	}
}
