package pm

import (
	"fmt"
	"os"
	"path/filepath"
)

// Unlink removes all of the symbolic links for a package, and returns a string
// slice of all of the symlinks that were removed.
func Unlink(baseDir, binDir, name, version string) ([]string, error) {
	mpath := filepath.Join(baseDir, name, version, "metadata.json")
	meta, err := LoadMetadata(mpath)
	if err != nil {
		return nil, err
	}

	var unlinked = make([]string, len(meta.Binaries))

	for i, bin := range meta.Binaries {
		targetPath := filepath.Join(baseDir, name, version, bin)
		linkPath := filepath.Join(binDir, filepath.Base(bin))
		if tgt, err := filepath.EvalSymlinks(targetPath); err != nil {
			return nil, err
		} else {
			if tgt != targetPath {
				return nil, fmt.Errorf(
					"mismatched symlink target: got %q, wanted %q",
					targetPath,
					tgt,
				)
			}
		}

		unlinked[i] = linkPath

		if err := os.Remove(linkPath); err != nil {
			return nil, err
		}
	}

	return unlinked, nil
}
