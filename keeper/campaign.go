package keeper

import "github.com/carp-cobain/referrals/domain"

// CampaignKeeper manages referral campaigns
type CampaignKeeper interface {
	CampaignReader
	CampaignWriter
}

// CampaignReader reads referral campaigns
type CampaignReader interface {
	GetCampaign(id uint64) (campaign domain.Campaign, err error)
	GetCampaigns(address string, cursor uint64, limit int) (uint64, []domain.Campaign)
}

// CampaignWriter writes referral campaigns
type CampaignWriter interface {
	CreateCampaign(address, name string) (campaign domain.Campaign, err error)
}
