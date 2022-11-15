package fs

import (
	"os"
	"path"
	"path/filepath"

	"github.com/starkandwayne/scs-broker/config"
)

func ArtifactStat(filename string) (string, os.FileInfo, error) {
	artifact := filepath.Join(config.Parsed.ArtifactsDir, path.Base(filename))

	stat, err := os.Stat(artifact)

	return artifact, stat, err
}
