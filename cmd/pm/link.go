package main

import (
	"log"

	"github.com/nesv/pm"
	"github.com/spf13/cobra"
)

var LinkCmd = &cobra.Command{
	Use:   "link [package]",
	Short: "Link the binaries of the specified package",
	Run:   runLink,
}

func runLink(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalln("not enough arguments")
	}

	linkedFiles, err := linkBinaries(rootBaseDir, args[0], rootBinDir)
	if err != nil {
		log.Fatalln(err)
	}

	if Verbose {
		for link, target := range linkedFiles {
			log.Printf("linked %q to %q", link, target)
		}
	}
}

func linkBinaries(baseDir, pkg, binDir string) (map[string]string, error) {
	return pm.Link(baseDir, pkg, binDir)
}
