package handler

import (
	"fmt"
	"strings"

	"github.com/carp-cobain/referrals/keeper"
	"github.com/gin-gonic/gin"
)

// CampaignHandler is the http/json api for managing referral campaigns
type CampaignHandler struct {
	campaignKeeper keeper.CampaignKeeper
}

// NewCampaignHandler creates a new referral campaign handler
func NewCampaignHandler(campaignKeeper keeper.CampaignKeeper) CampaignHandler {
	return CampaignHandler{campaignKeeper}
}

// GET /campaigns
// GetCampaigns gets a page of campaigns for an address
func (self CampaignHandler) GetCampaigns(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		badRequestJson(c, fmt.Errorf("address query param is required"))
		return
	}
	cursor, limit := getPageParams(c)
	next, campaigns := self.campaignKeeper.GetCampaigns(address, cursor, limit)
	okJson(c, gin.H{"cursor": next, "campaigns": campaigns})
}

// GET /campaigns/:id
// GetCampaign gets campaigns by ID
func (self CampaignHandler) GetCampaign(c *gin.Context) {
	id, err := uintParam(c, "id")
	if err != nil {
		badRequestJson(c, err)
		return
	}
	campaign, err := self.campaignKeeper.GetCampaign(id)
	if err != nil {
		notFoundJson(c, err)
		return
	}
	okJson(c, gin.H{"campaign": campaign})
}

// POST /campaigns
// CreateCampaign creates new named campaigns
func (self CampaignHandler) CreateCampaign(c *gin.Context) {
	var request CampaignRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		badRequestJson(c, err)
		return
	}
	address, name, err := request.Validate()
	if err != nil {
		badRequestJson(c, err)
		return
	}
	campaign, err := self.campaignKeeper.CreateCampaign(address, name)
	if err != nil {
		badRequestJson(c, err)
		return
	}
	okJson(c, gin.H{"campaign": campaign})
}

// CampaignRequest is the request type for creating referral campaigns.
type CampaignRequest struct {
	Address string `json:"address" binding:"required,min=41,max=61"`
	Name    string `json:"name"`
}

// Validate campaign request address
func (self CampaignRequest) Validate() (string, string, error) {
	address := strings.TrimSpace(self.Address)
	if address == "" {
		return "", "", fmt.Errorf("address cannot be blank")
	}
	if strings.ToLower(address) != address {
		return "", "", fmt.Errorf("address must be all lower case")
	}
	if !strings.HasPrefix(address, "tp") {
		return "", "", fmt.Errorf("address must have prefix: tp")
	}
	return address, strings.TrimSpace(self.Name), nil
}
