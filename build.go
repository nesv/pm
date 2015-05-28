package pm

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Build(m *Metadata, outputDir string) error {
	tarGzPath := filepath.Join(outputDir, fmt.Sprintf("%s.tar.gz", m.PackageName()))
	tgz, err := os.Create(tarGzPath)
	if err != nil {
		return fmt.Errorf("pm: error creating %q: %v", tarGzPath, err)
	}

	gzw := gzip.NewWriter(tgz)
	defer gzw.Close()

	tarw := tar.NewWriter(gzw)
	defer tarw.Close()

	// Add the metadata file.
	if err := writeMetadata(tarw, m); err != nil {
		return err
	}

	// Add the binaries.
	for _, binPath := range m.Binaries {
		if err := tarAddBinary(tarw, m, binPath); err != nil {
			return err
		}
	}

	// Add resources (extra files).
	for baseDir, globs := range m.Resources {
		for _, g := range globs {
			resourceGlob := filepath.Join(baseDir, g)
			matches, err := filepath.Glob(resourceGlob)
			if err != nil {
				return fmt.Errorf("pm: bad globbing pattern: %q", resourceGlob)
			}

			for _, match := range matches {
				if err := writeResource(tarw, m, match); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func writeResource(tw *tar.Writer, m *Metadata, resourcePath string) error {
	log.Println("adding resource", resourcePath)

	f, err := os.Open(resourcePath)
	if err != nil {
		return fmt.Errorf("pm: error opening resource %q: %v", resourcePath, err)
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("pm: error getting info for file %q: %v", resourcePath, err)
	}

	hdr, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return fmt.Errorf("pm: error creating tar header for file %q: %v", resourcePath, err)
	}

	name := filepath.Join(
		fmt.Sprintf("%s-%s", m.Name, m.Version),
		filepath.Clean(resourcePath),
	)
	hdr.Name = name
	hdr.Uid = 0
	hdr.Gid = 0
	hdr.Mode = 0444

	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("pm: error writing resource tar header for %q: %v", resourcePath, err)
	}

	if _, err := io.Copy(tw, f); err != nil {
		return fmt.Errorf("pm: error writing resource %q to archive: %v", resourcePath, err)
	}

	return nil
}

func writeMetadata(tw *tar.Writer, m *Metadata) error {
	log.Println("adding metadata")

	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return fmt.Errorf("pm: error marshaling metadata: %v", err)
	}

	now := time.Now()

	name := filepath.Join(fmt.Sprintf("%s-%s", m.Name, m.Version), "metadata.json")
	hdr := &tar.Header{
		Name:       name,
		Size:       int64(len(b)),
		ModTime:    now,
		AccessTime: now,
		ChangeTime: now,
		Mode:       0444,
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("pm: error writing metadata tar header: %v", err)
	}

	if _, err := tw.Write(b); err != nil {
		return fmt.Errorf("pm: error writing metadata: %v", err)
	}

	return nil
}

func tarAddBinary(tw *tar.Writer, m *Metadata, binPath string) error {
	log.Println("adding binary", binPath)

	bf, err := os.Open(binPath)
	if err != nil {
		return fmt.Errorf("pm: error opening binary file %q: %v", binPath, err)
	}
	defer bf.Close()

	info, err := bf.Stat()
	if err != nil {
		return fmt.Errorf("pm: error getting stat info on binary file %q: %v", binPath, err)
	}

	var hdr *tar.Header
	if hdr, err = tar.FileInfoHeader(info, ""); err != nil {
		return fmt.Errorf("pm: error creating header for binary file %q: %v", binPath, err)
	}

	name := filepath.Join(
		fmt.Sprintf("%s-%s", m.Name, m.Version),
		filepath.Clean(binPath),
	)
	hdr.Name = name
	hdr.Uid = 0
	hdr.Gid = 0
	hdr.Mode = 0555

	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("pm: error writing tar header for %q: %v", binPath, err)
	}

	if _, err := io.Copy(tw, bf); err != nil {
		return fmt.Errorf("pm: error writing binary file %q: %v", binPath, err)
	}

	return nil
}
