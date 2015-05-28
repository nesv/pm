package pm

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadMetadata(pth string) (*Metadata, error) {
	f, err := os.Open(pth)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var m Metadata

	dec := json.NewDecoder(f)
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

type Metadata struct {
	Architecture string              `json:"architecture"`
	Platform     string              `json:"platform"`
	Description  string              `json:"description"`
	Name         string              `json:"name"`
	Maintainer   string              `json:"maintainer"`
	SourceURL    string              `json:"source_url"`
	Resources    map[string][]string `json:"resources"`
	Binaries     []string            `json:"binaries"`
	Version      string              `json:"version"`
	Checksums    []string            `json:"checksums"`
}

func (m *Metadata) PackageName() string {
	return fmt.Sprintf("%s-%s-%s-%s",
		m.Name,
		m.Version,
		m.Platform,
		m.Architecture,
	)
}
