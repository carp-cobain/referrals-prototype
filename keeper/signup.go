package keeper

import "github.com/carp-cobain/referrals/domain"

// SignupKeeper manages referral campaign signups
type SignupKeeper interface {
	GetSignups(campaignID, cursor uint64, limit int) (uint64, []domain.Signup)
	CreateSignup(campaignID uint64, address string) (domain.Signup, error)
	UpdateSignup(campaignID, signupID uint64, status string) (domain.Signup, error)
}
