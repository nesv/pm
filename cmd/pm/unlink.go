package main

import (
	"log"

	"github.com/spf13/cobra"
)

var UnlinkCmd = &cobra.Command{
	Use:   "unlink [name]",
	Short: "Unlink a package",
	Run:   runUnlink,
}

func runUnlink(cmd *cobra.Command, args []string) {
	log.Fatalln("not implemented")
}
