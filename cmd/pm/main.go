package main

import (
	"log"

	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(0)

	rootCmd := &cobra.Command{
		Use:   "pm",
		Short: "pm is a brutally-simple package manager",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	rootCmd.AddCommand(
		BuildCmd,
	)

	rootCmd.Execute()
}
