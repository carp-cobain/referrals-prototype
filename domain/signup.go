package domain

import "time"

// Signup represents a blockchain address that signed up using a referral campaign.
type Signup struct {
	ID         uint64    `json:"id"`
	CampaignID uint64    `json:"campaignId"`
	Address    string    `json:"address"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
