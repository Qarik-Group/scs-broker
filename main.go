package main

import (
	"net/http"
	"os"
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
	"github.com/starkandwayne/config-server-broker/broker"
	"github.com/starkandwayne/config-server-broker/config"
	"github.com/starkandwayne/config-server-broker/httpartifacttransport"
)

var brokerLogger lager.Logger
var httpTransport httpartifacttransport.HttpArtifactTransport

func main() {

	brokerLogger = lager.NewLogger("config-server-broker")
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	brokerConf, err := config.ParseConfig()
	if err != nil {
		brokerLogger.Fatal("Reading config from env", err, lager.Data{
			"broker-config-environment-variable": config.ConfigEnvVarName,
		})
	}

	brokerLogger.Info("preparing transport")
	httpTransport = httpartifacttransport.NewHttpArtifactTransport(brokerConf, brokerLogger)

	brokerLogger.Info("downloading-artifact")
	url := brokerConf.ConfigServerDownloadURI

	if strings.HasPrefix(url, "file://") {
		httpTransport.EnableHttpFileTransport()
	}

	err = httpTransport.DownloadArtifact("spring-cloud-config-server.jar", url)
	if err != nil {
		brokerLogger.Fatal("Error downloading config-server jar", err, lager.Data{"uri": url})
	}
	brokerLogger.Info("download-Complete")
	brokerLogger.Info("starting")

	serviceBroker := &broker.ConfigServerBroker{
		Config: brokerConf,
		Logger: brokerLogger,
	}

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: brokerConf.Auth.Username,
		Password: brokerConf.Auth.Password,
	}

	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerCredentials)
	http.Handle("/", brokerAPI)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	brokerLogger.Info("listening", lager.Data{"port": port})
	brokerLogger.Fatal("http-listen", http.ListenAndServe("0.0.0.0:"+port, nil))
}

/*
func downloadArtifact(filename string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	os.Mkdir(broker.ArtifactsDir, 0777)
	// Create the file
	out, err := os.Create(broker.ArtifactsDir + "/" + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	num, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	brokerLogger.Info(fmt.Sprintf("Wrote: %d bytes", num))

	fi, err := os.Stat(broker.ArtifactsDir + "/" + filename)
	if err != nil {
		return err
	}

	brokerLogger.Info(fmt.Sprintf("Filename: %s Size: %d", fi.Name(), fi.Size()))

	return err
}
*/
