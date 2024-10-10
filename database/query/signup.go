package query

import (
	"fmt"

	"github.com/carp-cobain/referrals/database/model"
	"gorm.io/gorm"
)

// SelectSignup selects a signup by id
func SelectSignup(db *gorm.DB, id uint64) (signup model.Signup, err error) {
	err = db.Where("id = ?", id).First(&signup).Error
	return
}

// SelectSignups selects all referrals for a campaign.
func SelectSignups(db *gorm.DB, campaignID, cursor uint64, limit int) (signups []model.Signup) {
	db.Where("campaign_id = ?", campaignID).
		Where("id > ?", cursor).
		Order("id").
		Limit(limit).
		Find(&signups)
	return
}

// InsertSignup inserts a new referral for a campaign.
func InsertSignup(db *gorm.DB, campaignID uint64, address string) (signup model.Signup, err error) {
	signup = model.Signup{CampaignID: campaignID, Address: address, Status: "pending"}
	err = db.Create(&signup).Error
	return
}

// UpdateSignup updates a referral for a campaign.
func UpdateSignup(
	db *gorm.DB, campaignID, signupID uint64, status string) (signup model.Signup, err error) {

	signup, err = SelectSignup(db, signupID)
	if err != nil {
		return
	}
	if signup.CampaignID != campaignID {
		err = fmt.Errorf("invalid campaign: %d", campaignID)
		return
	}
	result := db.Model(&signup).Updates(updates{"status": status})
	err = result.Error
	return
}
