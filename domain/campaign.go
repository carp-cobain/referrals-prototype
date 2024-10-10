package domain

import "time"

// Campaign represents a referral campaign for a blockchain address.
type Campaign struct {
	ID        uint64    `json:"id"`
	Address   string    `json:"address"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
