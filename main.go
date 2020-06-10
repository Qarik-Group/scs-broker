package main

import (
	"fmt"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	// "github.com/pivotal-cf/brokerapi"
	// "github.com/starkandwayne/config-server-broker/broker"
)

func main() {
	brokerLogger := lager.NewLogger("config-server-broker")
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	brokerLogger.Info("Starting Config Server broker")

	// serviceBroker := &broker.ConfigServerBroker{}

	// brokerCredentials := brokerapi.BrokerCredentials{
	// 	Username: "admin",
	// 	Password: "admin",
	// }

	// brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerCredentials)
	// http.Handle("/", HelloServer)
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)

	// brokerLogger.Fatal("http-listen", http.ListenAndServe("0.0.0.0:8080", nil))
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
