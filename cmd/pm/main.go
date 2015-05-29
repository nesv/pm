package main

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	rootBaseDir  string
	rootCacheDir string
	rootBinDir   string
	Verbose      bool
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

	rootCmd.PersistentFlags().StringVarP(&rootBaseDir, "base-dir", "D", "/usr/local/lib", "Base working directory")
	rootCmd.PersistentFlags().StringVarP(&rootCacheDir, "cache-dir", "C", "/usr/local/var/pm/cache", "Location where pm will cache downloaded packages")
	rootCmd.PersistentFlags().StringVarP(&rootBinDir, "bin-dir", "B", "/usr/local/bin", "Directory where symlinks will be created for installed packages")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Print more output")

	rootCmd.AddCommand(
		BuildCmd,
		FetchCmd,
		InstallCmd,
		LinkCmd,
		UnpackCmd,
	)

	rootCmd.Execute()
}
