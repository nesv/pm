package main

import (
	"log"
	"net/url"
	"path/filepath"

	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install [url|path]",
	Short: "Install a package at the specified URL, or path",
	Long: `When a package URL is provided, the install command will fetch
	the remote package (unless the URL starts with "file://"), cache it,
	then extract the archive to $BASE_DIR/$NAME-$VERSION.

	If you do not want to have the retrieved package cached locally, you
	can specify the "--no-cache" option. When this option is specified,
	the retrieved package will be stored in memory before being unpacked.

	After the package has been retrieved, and unpacked, the install command
	will create symbolic links in $BIN_DIR that point to the
	binaries provided in the package.

	Calling "pm install ..." is equivalent to running the following
	commands:

		$ pm fetch <url>...
		$ pm unpack $NAME $VERSION
		$ pm link $NAME $VERSION
	`,
	Run: runInstall,
}

var (
	installNoCache bool
)

func init() {
	InstallCmd.Flags().BoolVarP(&installNoCache, "no-cache", "", false, "Do not cache the retrieved package")
}

func runInstall(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalln("too few arguments")
	}

	for _, urlStr := range args {
		u, err := url.Parse(urlStr)
		if err != nil {
			log.Fatalln(err)
		}

		if err := fetch(urlStr); err != nil {
			log.Fatalln(err)
		}

		cachedPath := filepath.Join(rootCacheDir, filepath.Base(u.Path))
		if err := unpack(cachedPath); err != nil {
			log.Fatalln(err)
		}

	}
}
