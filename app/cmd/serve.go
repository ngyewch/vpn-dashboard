package cmd

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/resources"
	"github.com/org-arl/cloud-infrastructure/software/vpn-dashboard/strongswan"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"net/http"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Serve",
		Args:  cobra.ExactArgs(0),
		RunE:  serve,
	}
)

func serve(cmd *cobra.Command, args []string) error {
	listenAddr, err := cmd.Flags().GetString("listen-addr")
	if err != nil {
		return err
	}

	var g errgroup.Group

	strongswanClient, err := strongswan.NewClient()
	if err != nil {
		panic(err)
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
		panic(err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("listen-addr", ":8080", "listen addr")
}
