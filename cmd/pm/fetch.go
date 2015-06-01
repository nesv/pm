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

	// Check to see if the package cache directory exists, and if it
	// doesn't, then create it.
	if fi, err := os.Stat(rootCacheDir); err != nil && os.IsNotExist(err) {
		if Verbose {
			log.Println("cache directory does not exist; creating it")
		}

		if err := os.MkdirAll(rootCacheDir, 0755); err != nil {
			log.Fatalln("failed to create cache directory:", err)
		}
	} else if !fi.IsDir() {
		log.Fatalln(rootCacheDir, "exists, but is not a directory")
	}

	destPath := filepath.Join(rootCacheDir, filepath.Base(u.Path))
	if f, err := os.Open(destPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed while checking cache for %q", filepath.Base(u.Path))
	} else if err != nil && os.IsNotExist(err) {
		if Verbose {
			log.Println("fetching", u.String())
		}

		r, err := pm.Fetch(urlStr)
		if err != nil {
			log.Fatalln(err)
		}
		defer r.Close()

		dest, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create %q: %v", destPath, err)
		}
		defer dest.Close()

		if _, err := io.Copy(dest, r); err != nil {
			return fmt.Errorf("failed to write file %q", destPath)
		}
	} else {
		f.Close()
		if Verbose {
			log.Printf("%s is already cached", u.Path)
		}
	}

	return nil
}
