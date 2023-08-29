package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/resources"
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/strongswan"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"net/http"
)

var (
	g errgroup.Group
)

func main() {
	strongswanClient, err := strongswan.NewClient()
	if err != nil {
		panic(err)
	}

	listenPort := 8080

	g.Go(func() error {
		e := echo.New()
		e.HideBanner = true
		e.Use(middleware.Recover())
		e.Use(middleware.Logger())

		strongswanService, err := strongswan.NewService(strongswanClient)
		if err != nil {
			return err
		}
		strongswanService.Install(e)

		rootFs, err := fs.Sub(resources.UiFS, "ui")
		if err != nil {
			return err
		}
		e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(rootFs))))

		err = e.Start(fmt.Sprintf(":%d", listenPort))
		if err != nil {
			return err
		}

		return nil
	})

	err = g.Wait()
	if err != nil {
		panic(err)
	}
}
