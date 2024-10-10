package repo_test

import (
	"testing"

	"github.com/carp-cobain/referrals/database"
	"github.com/carp-cobain/referrals/database/repo"
	"gorm.io/gorm"
)

func createTestDB(t *testing.T) *gorm.DB {
	db, err := database.Connect("file::memory:?cache=shared", 1)
	if err != nil {
		t.Fatalf("unable to connect to database: %+v", err)
	}
	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("unable to auto migrate: %+v", err)
	}
	return db
}

func TestCampaignRepo(t *testing.T) {
	db := createTestDB(t)
	referer := "tpabc123"
	campaignRepo := repo.NewCampaignRepo(db, db)
	campaign, err := campaignRepo.CreateCampaign(referer, "UnitTesting")
	if err != nil {
		t.Fatalf("failed to create referral campaign: %+v", err)
	}
	if _, err := campaignRepo.GetCampaign(campaign.ID); err != nil {
		t.Fatalf("failed to get campaign: %+v", err)
	}
	if _, campaigns := campaignRepo.GetCampaigns(referer, 0, 10); len(campaigns) != 1 {
		t.Fatalf("got unexpected number of campaigns")
	}
}

func TestSignupRepo(t *testing.T) {
	db := createTestDB(t)
	referer := "tpabc124"
	campgaignRepo := repo.NewCampaignRepo(db, db)
	campaign, err := campgaignRepo.CreateCampaign(referer, "UnitTesting")
	if err != nil {
		t.Fatalf("failed to create referral campaign: %+v", err)
	}
	referee := "tpabc125"
	signupRepo := repo.NewSignupRepo(db, db)
	if _, err := signupRepo.CreateSignup(campaign.ID, referee); err != nil {
		t.Fatalf("failed to create signup: %+v", err)
	}
	if _, signups := signupRepo.GetSignups(campaign.ID, 0, 10); len(signups) != 1 {
		t.Fatalf("got unexpected number of signups for campaign")
	}
	// Ensure people can't signup for thier own campaigns.
	if _, err := signupRepo.CreateSignup(campaign.ID, referer); err == nil {
		t.Fatalf("expected self referral error")
	}
}
