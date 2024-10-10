package model

import "github.com/carp-cobain/referrals/domain"

// Signup represents a blockchain address that signed up using a referral campaign.
type Signup struct {
	ID         uint64   `gorm:"primarykey"`
	CampaignID uint64   `gorm:"index;not null"`
	Campaign   Campaign `gorm:"foreignKey:CampaignID"`
	Address    string   `gorm:"uniqueIndex;not null"`
	Status     string
	CreatedAt  Time
	UpdatedAt  Time
}

// ToDomain converts a model to a domain object representation.
func (self Signup) ToDomain() domain.Signup {
	return domain.Signup{
		ID:         self.ID,
		CampaignID: self.CampaignID,
		Address:    self.Address,
		Status:     self.Status,
		CreatedAt:  self.CreatedAt.FromUnix(),
		UpdatedAt:  self.UpdatedAt.FromUnix(),
	}
}
