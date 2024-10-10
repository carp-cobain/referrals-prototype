package repo

import (
	"fmt"

	"github.com/carp-cobain/referrals/database/model"
	"github.com/carp-cobain/referrals/database/query"
	"github.com/carp-cobain/referrals/domain"
	"gorm.io/gorm"
)

// CampaignRepo manages referral campaigns in a database.
type CampaignRepo struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

// NewCampaignRepo creates a new repository for managing referral campaigns.
func NewCampaignRepo(readDB, writeDB *gorm.DB) CampaignRepo {
	return CampaignRepo{readDB, writeDB}
}

// GetCampaign gets a campaign by ID
func (self CampaignRepo) GetCampaign(id uint64) (campaign domain.Campaign, err error) {
	var model model.Campaign
	if model, err = query.SelectCampaign(self.readDB, id); err == nil {
		campaign = model.ToDomain()
	}
	if err != nil {
		err = fmt.Errorf("GetCampaign %d: %s", id, err.Error())
	}
	return
}

// GetCampaigns gets a page of referral campaigns for a blockchain address
func (self CampaignRepo) GetCampaigns(
	address string, cursor uint64, limit int) (next uint64, campaigns []domain.Campaign) {

	models := query.SelectCampaigns(self.readDB, address, cursor, limit)
	campaigns = make([]domain.Campaign, len(models))
	for i, model := range models {
		campaigns[i] = model.ToDomain()
		next = max(next, model.ID)
	}
	return
}

// CreateCampaign creates a new named campaign
func (self CampaignRepo) CreateCampaign(
	address, name string) (campaign domain.Campaign, err error) {

	if model, err := query.InsertCampaign(self.writeDB, address, name); err == nil {
		campaign = model.ToDomain()
	}
	if err != nil {
		err = fmt.Errorf("CreateCampaign: %s", err.Error())
	}
	return
}
