package main

import (
	"log"

	"github.com/nesv/pm"
	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a package",
	Long: `By default, the build command will look for a metadata.json file
	in the current directory, but you can specify the path to the metadata
	file using the -m or --metadata flags.

	For example, assuming the following fields in your metadata file:

	{
		"name": "foo",
		"version": "0.1.0",
		"platform": "linux",
		"architecture": "amd64",
		"checksums": ["sha512"],
		...
	}

	...build will produce "foo-0.1.0-linux-amd64.tar.gz", and
	"foo-0.1.0-linux-amd64.tar.gz.sha512".
	`,
	Run: runBuild,
}

var (
	buildMetadataFile string
	buildOutputDir    string
)

func init() {
	BuildCmd.Flags().StringVarP(&buildMetadataFile, "metadata", "m", "metadata.json", "Path to the metadata.json file to use for the build")
	BuildCmd.Flags().StringVarP(&buildOutputDir, "output-dir", "o", ".", "Change the directory to put the resulting archive and checksum files in")
}

func runBuild(cmd *cobra.Command, args []string) {
	metadata, err := pm.LoadMetadata(buildMetadataFile)
	if err != nil {
		log.Fatalln(err)
	}

	if err := pm.Build(metadata, buildOutputDir); err != nil {
		log.Fatalln(err)
	}
}
