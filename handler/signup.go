package handler

import (
	"fmt"
	"strings"

	"github.com/carp-cobain/referrals/keeper"
	"github.com/gin-gonic/gin"
)

// SignupHandler is the http/json api for managing referral campaign signups
type SignupHandler struct {
	campaignReader keeper.CampaignReader
	signupKeeper   keeper.SignupKeeper
}

// NewSignupHandler creates a new referral campaign handler
func NewSignupHandler(
	campaignReader keeper.CampaignReader, signupKeeper keeper.SignupKeeper) SignupHandler {

	return SignupHandler{campaignReader, signupKeeper}
}

// GET /campaigns/:id/signups
// GetSignups gets a page of signups for a referral campaign
func (self SignupHandler) GetSignups(c *gin.Context) {
	campaignID, err := uintParam(c, "id")
	if err != nil {
		badRequestJson(c, err)
		return
	}
	if _, err := self.campaignReader.GetCampaign(campaignID); err != nil {
		notFoundJson(c, err)
		return
	}
	cursor, limit := getPageParams(c)
	next, signups := self.signupKeeper.GetSignups(campaignID, cursor, limit)
	okJson(c, gin.H{"cursor": next, "signups": signups})
}

// POST /campaigns/:id/signups
// CreateSignup creates a signup for a referral campaign
func (self SignupHandler) CreateSignup(c *gin.Context) {
	campaignID, err := uintParam(c, "id")
	if err != nil {
		badRequestJson(c, err)
		return
	}
	var request SignupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		badRequestJson(c, err)
		return
	}
	address, err := request.Validate()
	if err != nil {
		badRequestJson(c, err)
		return
	}
	if _, err := self.campaignReader.GetCampaign(campaignID); err != nil {
		notFoundJson(c, err)
		return
	}
	signup, err := self.signupKeeper.CreateSignup(campaignID, address)
	if err != nil {
		badRequestJson(c, err)
		return
	}
	okJson(c, gin.H{"signup": signup})
}

// PATCH /campaigns/:id/signups/:sid
// UpdateSignup updates the status of a signup referral.
func (self SignupHandler) UpdateSignup(c *gin.Context) {
	campaignID, err := uintParam(c, "id")
	if err != nil {
		badRequestJson(c, err)
		return
	}
	signupID, err := uintParam(c, "sid")
	if err != nil {
		badRequestJson(c, err)
		return
	}
	var request UpdateSignupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		badRequestJson(c, err)
		return
	}
	status, err := request.Validate()
	if err != nil {
		badRequestJson(c, err)
		return
	}
	if _, err := self.campaignReader.GetCampaign(campaignID); err != nil {
		notFoundJson(c, err)
		return
	}
	signup, err := self.signupKeeper.UpdateSignup(campaignID, signupID, status)
	if err != nil {
		badRequestJson(c, err)
		return
	}
	okJson(c, gin.H{"signup": signup})
}

// SignupRequest is the request type for consuming referral campaigns.
type SignupRequest struct {
	Address string `json:"address" binding:"required,min=41,max=61"`
}

// Validate signup request fields
func (self SignupRequest) Validate() (string, error) {
	address := strings.TrimSpace(self.Address)
	if address == "" {
		return "", fmt.Errorf("address cannot be blank")
	}
	if strings.ToLower(address) != address {
		return "", fmt.Errorf("address must be all lower case")
	}
	if !strings.HasPrefix(address, "tp") {
		return "", fmt.Errorf("address must have prefix: tp")
	}
	return address, nil
}

// UpdateSignupRequest is the request type for updating signup status.
type UpdateSignupRequest struct {
	Status string `json:"status" binding:"required"`
}

// Validate ensures a signup status is a valid variant.
func (self UpdateSignupRequest) Validate() (string, error) {
	status := strings.ToLower(strings.TrimSpace(self.Status))
	if status == "" {
		return "", fmt.Errorf("invalid status: empty string")
	}
	if status != "verified" && status != "pending" {
		return "", fmt.Errorf("invalid status variant: %s", status)
	}
	return status, nil
}
