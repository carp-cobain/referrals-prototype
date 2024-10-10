package repo

import (
	"fmt"

	"github.com/carp-cobain/referrals/database/model"
	"github.com/carp-cobain/referrals/database/query"
	"github.com/carp-cobain/referrals/domain"
	"gorm.io/gorm"
)

// SignupRepo manages signups for referral campaigns.
type SignupRepo struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

// NewSignupRepo creates a new repository for managing signups for referral campaigns.
func NewSignupRepo(readDB, writeDB *gorm.DB) SignupRepo {
	return SignupRepo{readDB, writeDB}
}

// GetSignups gets a page of signups for a referral campaign.
func (self SignupRepo) GetSignups(
	campaignID, cursor uint64, limit int) (next uint64, signups []domain.Signup) {

	models := query.SelectSignups(self.readDB, campaignID, cursor, limit)
	signups = make([]domain.Signup, len(models))
	for i, model := range models {
		signups[i] = model.ToDomain()
		next = max(next, model.ID)
	}
	return
}

// CreateSignup creates a signup for a referral campaign
func (self SignupRepo) CreateSignup(
	campaignID uint64, address string) (signup domain.Signup, err error) {

	var campaign model.Campaign
	campaign, err = query.SelectCampaign(self.readDB, campaignID)
	if err != nil {
		err = fmt.Errorf("campaign %d: %s", campaignID, err.Error())
		return
	}
	if campaign.Address == address {
		err = fmt.Errorf("self referral error: %s", address)
		return
	}
	var model model.Signup
	if model, err = query.InsertSignup(self.writeDB, campaignID, address); err == nil {
		signup = model.ToDomain()
	}
	return
}

// UpdateSignup updates the status of a signup for a referral campaign.
func (self SignupRepo) UpdateSignup(
	campaignID, signupID uint64, status string) (signup domain.Signup, err error) {

	var model model.Signup
	if model, err = query.UpdateSignup(self.writeDB, campaignID, signupID, status); err == nil {
		signup = model.ToDomain()
	}
	return
}
