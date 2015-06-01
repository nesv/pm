package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nesv/pm"
	"github.com/spf13/cobra"
)

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean stuff up",
	Long: `
The "clean" command will only clean up the parts of itself that you tell it to.
Without any flags provided, "clean" will do nothing.

To have pm clean its cache, provide the "--cache" flag:

    $ pm clean --cache

To clean unpacked, unlinked archives, provide the "--unlinked" flag:

    $ pm clean --unlinked

For convenience, you can supply the "--all" flag, which is synonymous to
calling:

    $ pm clean --cache --unlinked
`,
	Run: runClean,
}

var (
	CleanAll      bool
	CleanCache    bool
	CleanUnlinked bool
)

func init() {
	CleanCmd.Flags().BoolVarP(&CleanAll, "all", "", false, "Clean everything")
	CleanCmd.Flags().BoolVarP(&CleanCache, "cache", "c", false, "Remove all cached files")
	CleanCmd.Flags().BoolVarP(&CleanUnlinked, "unlinked", "u", false, "Remove all unpacked, unlinked packages")
}

func runClean(cmd *cobra.Command, args []string) {
	if CleanAll {
		CleanCache = true
		CleanUnlinked = true
	}

	if CleanCache {
		cleanCache()
	}

	if CleanUnlinked {
		cleanUnlinked()
	}
}

func cleanCache() {
	cached, err := pm.ListCachedPackagesSlice(rootCacheDir)
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range cached {
		if Verbose {
			fmt.Println("removing", v, "from cache")
		}

		cachedFile := filepath.Join(
			rootCacheDir,
			fmt.Sprintf("%s-%s-%s.tar.gz", v, runtime.GOOS, runtime.GOARCH),
		)
		if err := os.Remove(cachedFile); err != nil {
			log.Fatalln(err)
		}
	}
}

func cleanUnlinked() {
	lp, err := pm.ListLinkedPackagesSlice(rootBaseDir, rootBinDir)
	if err != nil {
		log.Fatalln(err)
	}

	var linkedPkgs = make(map[string]struct{})
	for _, v := range lp {
		linkedPkgs[v] = struct{}{}
	}

	unpacked, err := pm.ListUnpackedPackagesSlice(rootBaseDir)
	if err != nil {
		log.Fatalln(err)
	}

	for _, u := range unpacked {
		if _, linked := linkedPkgs[u]; !linked {
			if Verbose {
				fmt.Println("removing", u)
			}

			parts := []string{rootBaseDir}
			parts = append(parts, strings.Split(u, pm.PackageFieldSeparator)...)
			pth := filepath.Join(parts...)
			if err := os.RemoveAll(pth); err != nil {
				log.Fatalln(err)
			}
		}
	}
}
