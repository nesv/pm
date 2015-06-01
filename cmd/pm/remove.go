package main

import (
	"log"
	"strings"

	"github.com/nesv/pm"
	"github.com/spf13/cobra"
)

var RemoveCmd = &cobra.Command{
	Use:   "remove [package...]",
	Short: "Unlink one or more packages, and remove their unpacked files",
	Long: `
This command will unlink a package and remove any of the unpacked files in
BASE_DIR. It is roughly equivalent to calling "pm unlink <package>", then
recursively removing BASEDIR/<name>/<version>.

If you also want to remove the cached package file, specify the "-p, --purge"
flag.
`,
	Run: runRemove,
}

var (
	RemovePurge bool
)

func init() {
	RemoveCmd.Flags().BoolVarP(&RemovePurge, "purge", "p", false, "Purge the cached package file")
}

func runRemove(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalln("not enough arguments")
	}

	linkedPkgs, err := pm.ListLinkedPackages(rootBaseDir, rootBinDir)
	if err != nil {
		log.Fatalln(err)
	}

	for _, pkg := range args {
		parts := strings.SplitN(pkg, pm.PackageFieldSeparator, 2)
		if len(parts) < 2 {
			log.Fatalln("invalid package name: must be in the format <name>-<version>")
		}

		if err := unlinkPkg(linkedPkgs, parts[0], parts[1]); err != nil {
			log.Fatalln(err)
		}

		if err := cleanUnlinkedPkg(parts[0], parts[1]); err != nil {
			log.Fatalln(err)
		}

		if RemovePurge {
			if err := cleanCachedPkg(parts[0], parts[1]); err != nil {
				log.Fatalln(err)
			}
		}
	}
}
