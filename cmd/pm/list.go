package main

import (
	"fmt"
	"log"

	"github.com/nesv/pm"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list [filter]",
	Short: "Generate lists of information about packages",
	Long: `
The following is a list of supported values for FILTER:

    * linked
    * cached
    * unpacked
`,
	Run: runList,
}

func runList(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalln("not enough arguments")
	}

	var pkgs []string
	var err error

	switch filter := args[0]; filter {
	case "linked":
		pkgs, err = pm.ListLinkedPackagesSlice(rootBaseDir, rootBinDir)
	case "cached":
		pkgs, err = pm.ListCachedPackagesSlice(rootCacheDir)
	case "unpacked":
		pkgs, err = pm.ListUnpackedPackagesSlice(rootBaseDir)
	default:
		log.Fatalln("unknown filter:", filter)
	}

	if err != nil {
		log.Fatalln(err)
	}

	for _, pkg := range pkgs {
		fmt.Println(pkg)
	}
}
