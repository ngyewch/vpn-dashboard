//go:generate fileb0x b0x.yaml
package main

import (
	"github.com/dghubble/sessions"
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/prometheus"
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/strongswan"
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/wireguard"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	sessionSecret = os.Getenv("SESSION_SECRET")
	cookieStore   = sessions.NewCookieStore([]byte(sessionSecret), nil)
	g             errgroup.Group
)

func main() {
	strongswanClient, err := strongswan.NewClient()
	if err != nil {
		panic(err)
	}
	wireguardClient, err := wireguard.NewClient()
	if err != nil {
		panic(err)
	}
	multiClient, err := NewMultiClient(strongswanClient, wireguardClient)
	if err != nil {
		panic(err)
	}

	serverMain := &http.Server{
		Addr:         ":8080",
		Handler:      routerMain(cookieStore, multiClient),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	serverPrometheus := &http.Server{
		Addr:         ":2112",
		Handler:      prometheus.RouterPrometheus(multiClient.GetAddresses),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		err := serverMain.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	g.Go(func() error {
		err := serverPrometheus.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
