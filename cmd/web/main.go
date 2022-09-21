package main

import (
	"flag"
	"log"
	"net/http"
	"terraform-plan-parser/cmd/web/provider"
	"terraform-plan-parser/internals"

	"github.com/gin-gonic/gin"
)

func main() {
	planPath := flag.String("planjson", "", "Path to the plan JSON file")
	flag.Parse()

	if *planPath == "" {
		log.Fatalf("planjson is a required argument")
	}

	parser := internals.NewParser()
	result, err := parser.ParseFile(*planPath)
	if err != nil {
		log.Fatalf("can't process file %s: %v", *planPath, err)
	}

	siteService := provider.NewSiteService(result)

	router := gin.Default()
	router.Static("/assets", "./site")
	router.LoadHTMLGlob("./templates/*.html")

	// Redirect root
	router.GET("/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/")
	})

	router.GET("/", siteService.Home)
	router.GET("/state", siteService.State)
	router.GET("/changes", siteService.Changes)
	router.GET("/config", siteService.Config)

	api := router.Group("api/v1")

	api.GET("state/", siteService.GetStateList)
	api.GET("state/:address", siteService.GetStateResource)
	api.GET("change/", siteService.GetChangeList)
	api.GET("change/:address", siteService.GetChangeResource)
	api.GET("config/", siteService.GetConfigAsync)
	api.GET("config/:address", siteService.GetItem)

	if err := router.Run(":7034"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
