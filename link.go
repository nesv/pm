package pm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Link(baseDir, pkg, binDir string) (map[string]string, error) {
	pkg = strings.Replace(pkg, PackageFieldSeparator, string(os.PathSeparator), -1)
	metadataPath := filepath.Join(baseDir, pkg, "metadata.json")
	m, err := LoadMetadata(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("error loading %s: %v", metadataPath, err)
	}

	var linkedFiles = make(map[string]string)

	for _, bin := range m.Binaries {
		link := filepath.Join(binDir, filepath.Base(bin))
		targ := filepath.Join(baseDir, pkg, bin)

		// If we get here, it means we were able to open a file
		// at the intended path of our symlink.
		//
		// What we must do at this point, is check to see if the file
		// is a symlink, and if it is, blow it away, and create out own.
		//
		// If the file is not a symlink, then bail out.
		info, err := os.Lstat(link)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("error getting info on file %q", link)
		} else if err != nil && os.IsNotExist(err) {
			// There is no file there. Create the symlink!
			if err := createSymlink(targ, link); err != nil {
				return nil, fmt.Errorf("error symlinking %q -> %q",
					link, targ)
			}

			linkedFiles[link] = targ

			// Blow through the rest of this loop iteration.
			continue
		}

		if info.Mode()&os.ModeSymlink == 0 {
			return nil, fmt.Errorf("error: %q is not a symlink", link)
		}

		info, err = os.Stat(link)
		if err != nil {
			return nil, fmt.Errorf("error: getting info on file %q", link)
		}

		// Alright. Destroy it.
		if err := os.Remove(link); err != nil {
			return nil, fmt.Errorf("error removing symlink %q: %v", link, err)
		}

		// Create the new symlink.
		if err := createSymlink(targ, link); err != nil {
			return nil, fmt.Errorf("error: symlinking %q -> %q",
				link, targ)
		}

		linkedFiles[link] = targ
	}

	return linkedFiles, nil
}

func createSymlink(target, link string) error {
	var err error

	if target, err = filepath.Abs(target); err != nil {
		return err
	}

	if link, err = filepath.Abs(link); err != nil {
		return err
	}

	return os.Symlink(target, link)
}
