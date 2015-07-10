package main

import (
	"log"
	"net/url"
	"path/filepath"
	"strings"

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

		if err := fetch(urlStr); err != nil {
			log.Fatalln(err)
		}

		cachedPath := filepath.Join(rootCacheDir, filepath.Base(u.Path))
		if err := unpack(cachedPath); err != nil {
			log.Fatalln(err)
		}

		pkgNameParts := strings.Split(args[0], "-")
		linkedFiles, err := linkBinaries(
			rootBaseDir,
			strings.Join(pkgNameParts[:len(pkgNameParts)-2], "-"),
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
