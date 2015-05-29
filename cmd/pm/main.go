package main

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	rootBaseDir  string
	rootCacheDir string
)

func main() {
	log.SetFlags(0)

	rootCmd := &cobra.Command{
		Use:   "pm",
		Short: "A brutally-simple package manager",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&rootBaseDir, "base-dir", "d", "/usr/local/lib", "Base working directory")
	rootCmd.PersistentFlags().StringVarP(&rootCacheDir, "cache-dir", "C", "/usr/local/var/pm/cache", "Location where pm will cache downloaded packages")

	rootCmd.AddCommand(
		BuildCmd,
		InstallCmd,
		FetchCmd,
		UnpackCmd,
	)

	rootCmd.Execute()
}
