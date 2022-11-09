package utilities

import "gopkg.in/yaml.v2"

func NewManifest(appName string, webCount uint) *Manifest {
	application := &ManifestApplication{
		Name:         appName,
		Instances:    webCount,
		DefaultRoute: false,
	}

	return &Manifest{
		Applications: []*ManifestApplication{
			application,
		},
	}
}

type Manifest struct {
	Applications []*ManifestApplication `yaml:"applications"`
}

func (m *Manifest) Bytes() ([]byte, error) {
	return yaml.Marshal(m)
}

type ManifestApplication struct {
	Name         string `yaml:"name"`
	Instances    uint   `yaml:"instances"`
	DefaultRoute bool   `yaml:"default_route"`
}
