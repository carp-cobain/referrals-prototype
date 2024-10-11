package main

import (
	"log"
	"os"

	"github.com/carp-cobain/referrals/database"
	"github.com/carp-cobain/referrals/database/repo"
	"github.com/carp-cobain/referrals/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	if _, ok := os.LookupEnv("DISABLE_COLOR"); ok {
		gin.DisableConsoleColor()
	}

	readDB, writeDB, err := database.ConnectAndMigrate()
	if err != nil {
		log.Panicf("unable to connnect to db: %+v", err)
	}

	campaignRepo := repo.NewCampaignRepo(readDB, writeDB)
	signupRepo := repo.NewSignupRepo(readDB, writeDB)

	campaignHandler := handler.NewCampaignHandler(campaignRepo)
	redirectHandler := handler.NewRedirectHandler(campaignRepo)
	signupHandler := handler.NewSignupHandler(campaignRepo, signupRepo)

	// Router
	r := gin.Default()

	// Signup redirects
	r.GET("/referrals/:id/signup", redirectHandler.CampaignSignupRedirect)

	// API
	v1 := r.Group("/referrals/api/v1")
	{
		v1.GET("/campaigns", campaignHandler.GetCampaigns)
		v1.POST("/campaigns", campaignHandler.CreateCampaign)
		v1.GET("/campaigns/:id", campaignHandler.GetCampaign)
		v1.GET("/campaigns/:id/signups", signupHandler.GetSignups)
		v1.POST("/campaigns/:id/signups", signupHandler.CreateSignup)
		v1.PATCH("/campaigns/:id/signups/:sid", signupHandler.UpdateSignup)
	}

	if err := r.Run(); err != nil {
		log.Panicf("unable to start referral server:  %+v", err)
	}
}
