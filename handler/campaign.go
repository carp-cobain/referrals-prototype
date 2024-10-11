package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/carp-cobain/referrals/keeper"
	"github.com/gin-gonic/gin"
)

// CookieName is the name for referral campaign cookies
var CookieName string = "c-referrer"

// MaxAge is the max age for referral campaign cookies
var MaxAge int = 30 * 24 * 60 * 60

// CampaignHandler is the http/json api for managing referral campaigns
type CampaignHandler struct {
	campaignKeeper keeper.CampaignKeeper
}

// NewCampaignHandler creates a new referral campaign handler
func NewCampaignHandler(campaignKeeper keeper.CampaignKeeper) CampaignHandler {
	return CampaignHandler{campaignKeeper}
}

// GET /referrals/:id/signup
// CampaignSignupRedirect drops a cookie and redirects the requestor to a signup URL.
func (self CampaignHandler) CampaignSignupRedirect(c *gin.Context) {
	location, path, domain := lookupSignupEnv()
	id, err := uintParam(c, "id")
	if err != nil {
		c.Redirect(http.StatusFound, location)
		return
	}
	if campaign, err := self.campaignKeeper.GetCampaign(id); err == nil {
		c.SetCookie(CookieName, fmt.Sprintf("%d", campaign.ID), MaxAge, path, domain, false, false)
	}
	c.Redirect(http.StatusFound, location)
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

// Lookup signup redirect info from env var
func lookupSignupEnv() (string, string, string) {
	url, ok := os.LookupEnv("SIGNUP_URL")
	if !ok {
		log.Panicf("SIGNUP_URL not defined")
	}
	path, ok := os.LookupEnv("SIGNUP_COOKIE_PATH")
	if !ok {
		log.Panicf("SIGNUP_COOKIE_PATH not defined")
	}
	domain, ok := os.LookupEnv("SIGNUP_COOKIE_DOMAIN")
	if !ok {
		log.Panicf("SIGNUP_COOKIE_DOMAIN not defined")
	}
	return url, path, domain
}
