package main

import (
	"github.com/dghubble/sessions"
	"github.com/gin-gonic/gin"
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/static"
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/strongswan"
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/web_service"
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/wireguard"
	"net/http"
)

func routerMain(cookieStore *sessions.CookieStore, multiClient *MultiClient) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	/*
	authService := web_service.NewAuthService(cookieStore)
	authService.Install(r)
	*/

	if multiClient.StrongswanClient != nil {
		strongSwanService, err := strongswan.NewService(multiClient.StrongswanClient)
		if err != nil {
			panic(err)
		}
		strongSwanService.Install(r)
	}
	if multiClient.WireguardClient != nil {
		wireguardService, err := wireguard.NewService(multiClient.WireguardClient)
		if err != nil {
			panic(err)
		}
		wireguardService.Install(r)
	}

	pingService, err := web_service.NewPingService(multiClient.GetAddresses)
	if err != nil {
		panic(err)
	}
	pingService.Install(r)

	r.GET("/", indexHandler)
	r.GET("/index.html", indexHandler)

	r.StaticFS("/public", static.HTTP)

	return r
}

func indexHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/public/")
}
