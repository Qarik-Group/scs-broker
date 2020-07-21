package main

import (
	"io"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
	"github.com/starkandwayne/config-server-broker/broker"
	"github.com/starkandwayne/config-server-broker/config"
)

func main() {
	brokerLogger := lager.NewLogger("config-server-broker")
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	brokerLogger.Info("Downloading artifact")
	url := "https://github.com/starkandwayne/spring-cloud-config-server/releases/download/v3.0.0/spring-cloud-config-server.jar"
	downloadArtifact("spring-cloud-config-server.jar", url)

	brokerLogger.Info("Download Complete")
	brokerLogger.Info("Starting Config Server broker")

	brokerConf, err := config.ParseConfig()
	if err != nil {
		brokerLogger.Fatal("Reading config from env", err, lager.Data{
			"broker-config-environment-variable": config.ConfigEnvVarName,
		})
	}

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

	brokerLogger.Fatal("http-listen", http.ListenAndServe("0.0.0.0:8080", nil))
}

func downloadArtifact(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	os.Mkdir(broker.ArtifactsDir, 0777)
	// Create the file
	out, err := os.Create(broker.ArtifactsDir + "/" + filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
