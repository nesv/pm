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
Each of the flags supported by this command will present a list of packages
using a logical AND. For example, if you provide the "--linked" and "--cached"
flags, the returned list of packages will be those that are cached and linked.
`,
	Run: runList,
}

var (
	ListLinked   bool
	ListCached   bool
	ListUnpacked bool
)

func init() {
	ListCmd.Flags().BoolVarP(&ListLinked, "linked", "i", false, `List all linked ("installed") packages`)
	ListCmd.Flags().BoolVarP(&ListCached, "cached", "c", false, `List all cached packages`)
	ListCmd.Flags().BoolVarP(&ListUnpacked, "unpacked", "x", false, `List all unpacked ("extracted") packages`)
}

func runList(cmd *cobra.Command, args []string) {
	var pkgList = make(map[string]struct{})

	if ListLinked {
		pkgs, err := pm.ListLinkedPackagesSlice(rootBaseDir, rootBinDir)
		if err != nil {
			log.Fatalln(err)
		}

		for _, p := range pkgs {
			pkgList[p] = struct{}{}
		}
	}

	if ListCached {
		pkgs, err := pm.ListCachedPackagesSlice(rootCacheDir)
		if err != nil {
			log.Fatalln(err)
		}

		for _, p := range pkgs {
			pkgList[p] = struct{}{}
		}
	}

	if ListUnpacked {
		pkgs, err := pm.ListUnpackedPackagesSlice(rootBaseDir)
		if err != nil {
			log.Fatalln(err)
		}

		for _, p := range pkgs {
			pkgList[p] = struct{}{}
		}
	}

	for pkg, _ := range pkgList {
		fmt.Println(pkg)
	}
}
