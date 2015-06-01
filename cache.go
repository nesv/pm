package pm

import (
	"fmt"
	"path/filepath"
	"strings"
)

func ListCachedPackagesSlice(cacheDir string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(cacheDir, "*.tar.gz"))
	if err != nil {
		return nil, err
	}

	var pkgs = make([]string, 0)

	for _, m := range matches {
		parts := strings.Split(filepath.Base(m), "-")
		if len(parts) == 4 {
			pkgs = append(pkgs, fmt.Sprintf("%s-%s", parts[0], parts[1]))
		}
	}

	return pkgs, nil
}
