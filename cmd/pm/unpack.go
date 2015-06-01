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

var UnpackCmd = &cobra.Command{
	Use:   "unpack [package]",
	Short: "Unpack a cached version of a package",
	Run:   runUnpack,
}

func runUnpack(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalln("not enough arguments")
	}

	pkgFilename := strings.Join(
		[]string{args[0], runtime.GOOS, runtime.GOARCH},
		pm.PackageFieldSeparator,
	)
	pkgPath := filepath.Join(rootCacheDir, pkgFilename+".tar.gz")

	if err := unpack(pkgPath); err != nil {
		log.Fatalln(err)
	}
}

func unpack(pkgPath string) error {
	pkgFilename := filepath.Base(pkgPath)

	f, err := os.Open(pkgPath)
	if err != nil {
		return fmt.Errorf("error: package file %q is not cached", pkgFilename)
	}
	defer f.Close()

	if Verbose {
		log.Println("unpacking", pkgFilename)
	}

	unpackedFiles, err := pm.Unpack(rootBaseDir, f)
	if err != nil {
		return err
	}

	if Verbose {
		for _, fname := range unpackedFiles {
			log.Println("unpacked", fname)
		}
	}

	return nil
}
