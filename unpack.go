package pm

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
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
	var meta Metadata

	for {
		hdr, err := tr.Next()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(tr); err != nil {
			return nil, fmt.Errorf("failed to read %q from archive: %v", hdr.Name, err)
		}

		br := bytes.NewReader(buf.Bytes())

		if filepath.Base(hdr.Name) == "metadata.json" {
			dec := json.NewDecoder(br)
			if err := dec.Decode(&meta); err != nil {
				return nil, fmt.Errorf("failed to load metadata file from archive: %v", err)
			}

			br.Seek(0, 0)
		}

		tgt := filepath.Join(baseDir, hdr.Name)

		if err := os.MkdirAll(filepath.Dir(tgt), 0755); err != nil {
			return nil, err
		}

		f, err := os.Create(tgt)
		if err != nil {
			return nil, err
		}

		if _, err := io.Copy(f, br); err != nil {
			f.Close()
			return nil, err
		}

		files = append(files, tgt)

		f.Close()
	}

	// Fix the permissions on the unpacked binaries.
	for _, bin := range meta.Binaries {
		binPath := filepath.Join(baseDir, fmt.Sprintf("%s-%s", meta.Name, meta.Version), bin)
		if err := os.Chmod(binPath, 0755); err != nil {
			return nil, fmt.Errorf("failed to change permissions for %q: %v", binPath, err)
		}
	}

	return files, nil
}
