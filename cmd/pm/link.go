package main

import (
	"log"

	"github.com/nesv/pm"
	"github.com/spf13/cobra"
)

var LinkCmd = &cobra.Command{
	Use:   "link [package] [version]",
	Short: "Link the binaries of the specified package and version",
	Run:   runLink,
}

func runLink(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		log.Fatalln("not enough arguments")
	}

	linkedFiles, err := linkBinaries(rootBaseDir, args[0], args[1], rootBinDir)
	if err != nil {
		log.Fatalln(err)
	}

	if Verbose {
		for link, target := range linkedFiles {
			log.Printf("linked %q to %q", link, target)
		}
	}
}

func linkBinaries(baseDir, name, version, binDir string) (map[string]string, error) {
	return pm.Link(baseDir, name, version, binDir)
}
