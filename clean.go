package pm

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const PackageFieldSeparator = "-"

func ListUnpackedPackagesSlice(baseDir string) ([]string, error) {
	var err error

	baseDir, err = filepath.Abs(baseDir)
	if err != nil {
		return nil, err
	}

	matches, err := filepath.Glob(filepath.Join(baseDir, "*", "*", "metadata.json"))
	if err != nil {
		return nil, err
	}

	var pkgs = make([]string, 0)

	for _, mpath := range matches {
		m, err := LoadMetadata(mpath)
		if err != nil && os.IsNotExist(err) {
			continue
		}

		pkg := fmt.Sprintf("%s%s%s", m.Name, PackageFieldSeparator, m.Version)
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

// CleanUnlinkedPackages removes all of the unpacked archives that do not
// currently have a symbolic link in binDir targeting a binary within them.
//
// This function will return a list of the names and vesrions of unpacked
// archives, that were removed, in the format:
//
//     name/version
func CleanUnlinkedPackages(baseDir, binDir string) ([]string, error) {
	var removed = make([]string, 0)

	linkedPkgs, err := ListLinkedPackages(baseDir, binDir)
	if err != nil {
		return nil, err
	}

	baseDir, err = filepath.Abs(baseDir)
	if err != nil {
		return nil, err
	}

	binDir, err = filepath.Abs(binDir)
	if err != nil {
		return nil, err
	}

	matches, err := filepath.Glob(filepath.Join(baseDir, "*", "*"))
	if err != nil {
		return nil, err
	}

	log.Printf("DEBUG unpacked: %#v", matches)
	return nil, nil

	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-directories.
		if !info.IsDir() {
			return nil
		}

		trimmed := strings.TrimPrefix(path, baseDir)
		parts := strings.Split(trimmed, string(os.PathSeparator))
		if len(parts) == 2 {
			// We are only going to handle the cases where we have
			// as much of the pathname as we need.
			if vsn, ok := linkedPkgs[parts[0]]; ok && parts[1] == vsn {
				// This package is linked. Skip it!
				return nil
			}

			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("error recursively removing %q: %v", path, err)
			}

			removed = append(removed, trimmed)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return removed, nil
}

// ListLinkedPackagesSlice is a convenience function that creates a string slice
// from the map produced by ListLinkedPackages.
func ListLinkedPackagesSlice(baseDir, binDir string) ([]string, error) {
	pkgs, err := ListLinkedPackages(baseDir, binDir)
	if err != nil {
		return nil, err
	}

	var linked = make([]string, 0)
	for pkg, vsn := range pkgs {
		linked = append(linked, fmt.Sprintf("%s-%s", pkg, vsn))
	}

	return linked, nil
}

type DirtyLinkError struct {
	LinkPath, PackageName, Version, WantVersion string
}

func (e *DirtyLinkError) Error() string {
	return fmt.Sprintf(
		"package version mismatch: %q is from %v-%v, but we want version %v",
		e.LinkPath,
		e.PackageName,
		e.Version,
		e.WantVersion,
	)
}

// ListLinkedPackages returns a map of unpacked archives that currently have
// symbolic links in binDir targeting the binaries within them.
//
// The returned map holds the format:
//
//     map["package_name"] = "package_version"
//
func ListLinkedPackages(baseDir, binDir string) (map[string]string, error) {
	var pkgSet = make(map[string]string)

	var err error

	binDir, err = filepath.Abs(binDir)
	if err != nil {
		return nil, err
	}

	baseDir, err = filepath.Abs(baseDir)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(binDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories.
		if info.IsDir() {
			return nil
		}

		// Skip regular files.
		if info.Mode()&os.ModeSymlink != os.ModeSymlink {
			return nil
		}

		target, err := filepath.EvalSymlinks(path)
		if err != nil {
			return fmt.Errorf("failed to resolve symlink for %q: %v", path, err)
		}

		if !strings.HasPrefix(target, baseDir) {
			// It doesn't look like this symlink points to anything
			// in our base directory. Skip it.
			return nil
		}

		// If we get here, this means that the symlink points to
		// something in our base directory.
		//
		// Let's try and divine the package name and version, from
		// the target path by stripping the baseDir prefix off of
		// the target path,
		target = strings.TrimPrefix(target, baseDir+string(os.PathSeparator))
		parts := strings.SplitN(target, string(os.PathSeparator), 3)
		if len(parts) < 3 {
			// If we have fewer than three parts to the stripped
			// target path, then we are just going to skip over this
			// file.
			return nil
		}

		name := parts[0]
		version := parts[1]

		if vsn, ok := pkgSet[name]; !ok {
			// This is our first time stumbling across this package.
			pkgSet[name] = version
		} else if vsn != version {
			// The versions do not match! This probably means that
			// there was an unclean package unlinking at some point.
			// Error out.
			return &DirtyLinkError{
				LinkPath:    path,
				PackageName: name,
				Version:     vsn,
				WantVersion: version,
			}
		}

		return nil
	})

	return pkgSet, err
}
