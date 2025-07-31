package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ngyewch/vpn-dashboard/resources"
	"github.com/ngyewch/vpn-dashboard/strongswan"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"net/http"
)

func doServe(ctx context.Context, cmd *cli.Command) error {
	listenAddr := cmd.String(flagListenAddr.Name)

	var g errgroup.Group

	strongswanClient, err := strongswan.NewClient()
	if err != nil {
		return err
	}

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

		reg := prometheus.NewPedanticRegistry()
		collector := strongswan.NewCollector(strongswanClient)
		reg.MustRegister(collector)
		e.GET("/metrics", echo.WrapHandler(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))

		rootFs, err := fs.Sub(resources.UiFS, "ui")
		if err != nil {
			return err
		}
		e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(rootFs))))

		err = e.Start(listenAddr)
		if err != nil {
			return err
		}

		return nil
	})

	err = g.Wait()
	if err != nil {
		return err
	}

	return nil
}
