package httpartifacttransport

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/starkandwayne/config-server-broker/broker"
	"github.com/starkandwayne/config-server-broker/config"
)

type HttpArtifactTransport struct {
	Config config.Config
	Logger lager.Logger
	Client *http.Client
}

func NewHttpArtifactTransport(config config.Config, logger lager.Logger) HttpArtifactTransport {
	return HttpArtifactTransport{
		Config: config,
		Logger: logger,
	}
}

func (transport *HttpArtifactTransport) EnableHttpFileTransport() {
	t := &http.Transport{}
	os.Mkdir(broker.ArtifactsDir, 0777)
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("./"+broker.ArtifactsDir)))
	transport.Client = &http.Client{Transport: t}
}

func (transport *HttpArtifactTransport) DownloadArtifact(filename string, url string) error {

	if transport.Client == nil {
		transport.Client = &http.Client{}
	}

	resp, err := transport.Client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	os.Mkdir(broker.ArtifactsDir, 0777)
	// Create the file
	out, err := os.Create("./" + broker.ArtifactsDir + "/" + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	num, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	transport.Logger.Info(fmt.Sprintf("Wrote: %d bytes", num))

	fi, err := os.Stat("./" + broker.ArtifactsDir + "/" + filename)
	if err != nil {
		return err
	}

	transport.Logger.Info(fmt.Sprintf("Filename: %s Size: %d", fi.Name(), fi.Size()))

	return err
}
