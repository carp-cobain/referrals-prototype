package query

import (
	"github.com/carp-cobain/referrals/database/model"
	"gorm.io/gorm"
)

// SelectCampaign selects a referral campaign by id
func SelectCampaign(db *gorm.DB, id uint64) (campaign model.Campaign, err error) {
	err = db.Where("id = ?", id).First(&campaign).Error
	return
}

// SelectCampaigns selects a page of referral campaigns for an address
func SelectCampaigns(
	db *gorm.DB, address string, cursor uint64, limit int) (campaigns []model.Campaign) {

	db.Where("address = ?", address).
		Where("id > ?", cursor).
		Order("id").
		Limit(limit).
		Find(&campaigns)

	return
}

// InsertCampaign inserts a new named campaign for an address
func InsertCampaign(db *gorm.DB, address, name string) (campaign model.Campaign, err error) {
	campaign = model.Campaign{Address: address, Name: name}
	err = db.Create(&campaign).Error
	return
}
