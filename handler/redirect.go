package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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
	signupKeeper   keeper.SignupKeeper
}

// NewRedirectHandler creates a new referral campaign handler
func NewRedirectHandler(
	campaignReader keeper.CampaignReader, signupKeeper keeper.SignupKeeper) RedirectHandler {
	return RedirectHandler{campaignReader, signupKeeper}
}

// GET /referrals/:id/signup
// Signup drops a cookie and redirects the requestor to a signup URL.
func (self RedirectHandler) Signup(c *gin.Context) {
	signupURL, path, domain := lookupSignupEnv()
	campaignID, err := uintParam(c, "id")
	if err != nil {
		c.Redirect(http.StatusFound, signupURL)
		return
	}
	if campaign, err := self.campaignReader.GetCampaign(campaignID); err == nil {
		value := fmt.Sprintf("%d", campaign.ID)
		c.SetCookie(CookieName, value, MaxAge, path, domain, false, false)
	}
	c.Redirect(http.StatusFound, signupURL)
}

// GET /referrals
// Referrals records referrals from campaign cookie and redirects to a provided URL.
func (self RedirectHandler) Referrals(c *gin.Context) {
	// Check for a redirect URL in query params. Use signup url if not provided.
	signup, _, _ := lookupSignupEnv()
	url := c.Query("url")
	if url == "" {
		url = signup
	}
	// Assume a redirect here from signup complete, so try and pull blockchain address.
	// If not found, redirect
	address := c.GetHeader("x-account-address")
	if address == "" {
		log.Printf("no address header; redirecting to: %s", url)
		c.Redirect(http.StatusFound, url)
		return
	}
	// Check for cookie, redirect if not found.
	cookie, err := c.Cookie(CookieName)
	if err == http.ErrNoCookie {
		log.Printf("no cookie; redirecting to: %s", url)
		c.Redirect(http.StatusFound, url)
		return
	}
	// Get campaign from cookie
	campaignID, err := strconv.ParseUint(cookie, 10, 64)
	if err != nil {
		log.Printf("failed to parse referral campaign cookie ID: %s", err.Error())
		c.Redirect(http.StatusFound, url)
		return
	}
	campaign, err := self.campaignReader.GetCampaign(campaignID)
	if err != nil {
		log.Printf("failed to get referral campaign %d: %s", campaignID, err.Error())
		c.Redirect(http.StatusFound, url)
		return
	}
	// Store referral signup
	if _, err := self.signupKeeper.CreateSignup(campaign.ID, address); err != nil {
		log.Printf("failed to record signup referral: %s", err.Error())
	}
	// Send user on their way
	c.Redirect(http.StatusFound, url)
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
