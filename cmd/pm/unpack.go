package main

import (
	"log"

	"github.com/spf13/cobra"
)

var UnpackCmd = &cobra.Command{
	Use:   "unpack [NAME-VERSION]",
	Short: "Unpack a cached version of a package",
	Run:   runUnpack,
}

func runUnpack(cmd *cobra.Command, args []string) {
	log.Fatalln("not implemented")
}
