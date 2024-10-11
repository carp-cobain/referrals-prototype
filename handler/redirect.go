package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/carp-cobain/referrals/keeper"
	"github.com/gin-gonic/gin"
)

// CookieName is the name for referral campaign cookies
var CookieName string = "_referral_campaign"

// MaxAge is the max age for referral campaign cookies
var MaxAge int = 30 * 24 * 60 * 60

// RedirectHandler is the http/json api for managing referral campaigns
type RedirectHandler struct {
	campaignReader keeper.CampaignReader
}

// NewRedirectHandler creates a new referral campaign handler
func NewRedirectHandler(campaignReader keeper.CampaignReader) RedirectHandler {
	return RedirectHandler{campaignReader}
}

// GET /referrals/:id/signup
// CampaignSignupRedirect drops a cookie and redirects the requestor to a signup URL.
func (self RedirectHandler) CampaignSignupRedirect(c *gin.Context) {
	location, path, domain := lookupSignupEnv()
	id, err := uintParam(c, "id")
	if err != nil {
		c.Redirect(http.StatusFound, location)
		return
	}
	if _, err := c.Cookie(CookieName); err == nil {
		c.Redirect(http.StatusFound, location)
		return
	}
	if campaign, err := self.campaignReader.GetCampaign(id); err == nil {
		value := fmt.Sprintf("%d,%s", id, campaign.Address)
		c.SetCookie(CookieName, value, MaxAge, path, domain, false, false)
	}
	c.Redirect(http.StatusFound, location)
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
