package cmd

import (
	"fmt"
	versionInfoCobra "github.com/ngyewch/go-versioninfo/cobra"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", appName),
		Short: "strongSwan Admin UI",
		RunE:  help,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func help(cmd *cobra.Command, args []string) error {
	err := cmd.Help()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	versionInfoCobra.AddVersionCmd(rootCmd, nil)
}

func initConfig() {
	// do nothing
}
