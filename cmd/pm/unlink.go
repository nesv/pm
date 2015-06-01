package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/nesv/pm"
	"github.com/spf13/cobra"
)

var UnlinkCmd = &cobra.Command{
	Use:   "unlink [package]",
	Short: "Unlink a package",
	Run:   runUnlink,
}

func runUnlink(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalln("not enough arguments")
	}

	pkgs, err := pm.ListLinkedPackages(rootBaseDir, rootBinDir)
	if err != nil {
		log.Fatalln(err)
	}

	parts := strings.SplitN(args[0], pm.PackageFieldSeparator, 2)
	if len(parts) < 2 {
		log.Fatalln("invalid package name: must be in the format <name>-<version>")
	}

	if version, linked := pkgs[parts[0]]; !linked || parts[1] != version {
		log.Fatalln("package is not linked")
	}

	unlinked, err := pm.Unlink(rootBaseDir, rootBinDir, parts[0], parts[1])
	if err != nil {
		log.Fatalln(err)
	}

	if Verbose {
		for _, link := range unlinked {
			fmt.Println("removed link", link)
		}
		fmt.Println("unlinked", args[0])
	}
}
