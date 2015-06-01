package main

import (
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/nesv/pm"
	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install [url|path]",
	Short: "Install a package at the specified URL, or path",
	Long: `
When a package URL is provided, the install command will fetch the remote
package (unless the URL starts with "file://"), cache it, then extract the
archive to $BASE_DIR/$NAME-$VERSION.

After the package has been retrieved, and unpacked, the install command will
create symbolic links in $BIN_DIR that point to the binaries provided in the
package.

Calling "pm install ..." is equivalent to running the following commands:

    $ pm fetch <url>...
    $ pm unpack <package>
    $ pm link <package>

where <package> is of the format "name-version" (e.g. "foo-1.0.0").
`,
	Run: runInstall,
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

		// Divine the package name and version from the URL path.
		//
		// The parts of the package's basename will be useful in a bit,
		// but for now, this is also being leveraged as a means of
		// making sure that the package's basename is in the correct
		// format.
		pkgNameParts := strings.SplitN(filepath.Base(u.Path), "-", 3)
		if len(pkgNameParts) < 3 {
			// Something is messed up with the package name.
			log.Fatalln("package name is malformed")
		}

		if err := fetch(urlStr); err != nil {
			log.Fatalln(err)
		}

		cachedPath := filepath.Join(rootCacheDir, filepath.Base(u.Path))
		if err := unpack(cachedPath); err != nil {
			log.Fatalln(err)
		}

		linkedFiles, err := linkBinaries(
			rootBaseDir,
			strings.Join(pkgNameParts[0:2], pm.PackageFieldSeparator),
			rootBinDir,
		)
		if err != nil {
			log.Fatalln(err)
		}

		if Verbose {
			for link, target := range linkedFiles {
				log.Printf("linked %q -> %q", link, target)
			}
		}
	}
}
