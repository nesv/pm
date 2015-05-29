package pm

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Unpack(baseDir string, pkg io.Reader) ([]string, error) {
	if baseDir == "" {
		return nil, fmt.Errorf("pm: no base directory specified")
	}

	gz, err := gzip.NewReader(pkg)
	if err != nil {
		return nil, fmt.Errorf("pm: cannot unpack: package is not compressed with gzip")
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	var files = make([]string, 0)

	for {
		hdr, err := tr.Next()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		tgt := filepath.Join(baseDir, hdr.Name)

		if err := os.MkdirAll(filepath.Dir(tgt), 0755); err != nil {
			return nil, err
		}

		f, err := os.Create(tgt)
		if err != nil {
			return nil, err
		}

		if _, err := io.Copy(f, tr); err != nil {
			f.Close()
			return nil, err
		}

		files = append(files, tgt)

		f.Close()
	}

	return files, nil
}
