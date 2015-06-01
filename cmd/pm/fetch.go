package main

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/nesv/pm"
	"github.com/spf13/cobra"
)

var FetchCmd = &cobra.Command{
	Use:   "fetch [url]",
	Short: "Fetch a remote package, and store it in the local cache",
	Long: `
pm currently supports fetching packages with the following URL schemes:

    * http
    * https
    * file

You can supply more than one package URL to the fetch command.
	`,
	Run: runFetch,
}

func runFetch(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalln("no package URLs specified")
	}

	// Check to see if the package cache directory exists, and if it
	// doesn't, then create it.
	cacheDir, err := os.Open(rootCacheDir)
	if err != nil && os.IsNotExist(err) {
		if e := os.MkdirAll(rootCacheDir, 0755); e != nil {
			log.Fatalln(err)
		}
	}
	cacheDir.Close()

	for _, urlStr := range args {
		if err := fetch(urlStr); err != nil {
			log.Fatalln(err)
		}
	}
}

func fetch(urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Fatalln(err)
	}

	destPath := filepath.Join(rootCacheDir, filepath.Base(u.Path))
	if f, err := os.Open(destPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed while checking cache for %q", filepath.Base(u.Path))
	} else if err != nil && os.IsNotExist(err) {
		log.Println("fetching", u.Path)

		r, err := pm.Fetch(urlStr)
		if err != nil {
			log.Fatalln(err)
		}
		defer r.Close()

		dest, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create %q", destPath)
		}
		defer dest.Close()

		if _, err := io.Copy(dest, r); err != nil {
			return fmt.Errorf("failed to write file %q", destPath)
		}
	} else {
		f.Close()
		log.Printf("using %q from cache", filepath.Base(u.Path))
	}

	return nil
}
