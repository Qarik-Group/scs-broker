package main

import (
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi/v7"
	"github.com/starkandwayne/scs-broker/broker"
	"github.com/starkandwayne/scs-broker/broker/configserver"
	"github.com/starkandwayne/scs-broker/broker/implementation"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry"
	"github.com/starkandwayne/scs-broker/config"
	"github.com/starkandwayne/scs-broker/logger"
)

var brokerLogger lager.Logger

func main() {
	brokerLogger = lager.NewLogger("scs-broker")
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	logger.Setup(brokerLogger)

	err := config.ParseConfig()
	if err != nil {
		logger.Fatal("Reading config from env", err, lager.Data{
			"broker-config-environment-variable": config.ConfigEnvVarName,
		})
	}

	logger.Info("starting")

	implementation.Register("configserver", configserver.New())
	implementation.Register("serviceregistry", serviceregistry.New())

	serviceBroker := broker.New()

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: config.Parsed.Auth.Username,
		Password: config.Parsed.Auth.Password,
	}

	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerCredentials)
	http.Handle("/", brokerAPI)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("listening", lager.Data{"port": port})
	logger.Fatal("http-listen", http.ListenAndServe("0.0.0.0:"+port, nil))
}
