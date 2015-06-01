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

	parts := strings.SplitN(args[0], pm.PackageFieldSeparator, 2)
	if len(parts) < 2 {
		log.Fatalln("invalid package name: must be in the format <name>-<version>")
	}

	pkgs, err := pm.ListLinkedPackages(rootBaseDir, rootBinDir)
	if err != nil {
		log.Fatalln(err)
	}

	if err := unlinkPkg(pkgs, parts[0], parts[1]); err != nil {
		log.Fatalln(err)
	}
}

func unlinkPkg(linkedPkgs map[string]string, name, version string) error {
	if vsn, linked := linkedPkgs[name]; !linked || version != vsn {
		return fmt.Errorf("package is not linked")
	}

	unlinked, err := pm.Unlink(rootBaseDir, rootBinDir, name, version)
	if err != nil {
		return err
	}

	if Verbose {
		for _, link := range unlinked {
			fmt.Println("removed link", link)
		}
		fmt.Println("unlinked", strings.Join([]string{name, version}, pm.PackageFieldSeparator))
	}

	return nil
}
