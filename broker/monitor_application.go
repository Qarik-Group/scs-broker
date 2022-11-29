package broker

import (
	"time"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/lager"
	"github.com/cloudfoundry-community/go-cfclient"
)

func (broker *SCSBroker) MonitorApplicationStartup(cfClient *ccv3.Client, community *cfclient.Client, logger lager.Logger, appGUID string) (bool, error) {

	waittime := 30
	timepassed := 0

	for timepassed < waittime {
		time.Sleep(time.Second)
		timepassed += 1
		successStart, err := broker.checkApplicationStatus(cfClient, community, logger, appGUID)
		if err != nil {
			return false, err
		}
		if !successStart {
			return successStart, err
		}
	}

	return true, nil

}

func (broker *SCSBroker) checkApplicationStatus(cfClient *ccv3.Client, community *cfclient.Client, logger lager.Logger, appGUID string) (bool, error) {
	stats, err := getProcessStatsByAppAndType(cfClient, community, broker.Logger, appGUID, "web")
	if err != nil {
		return false, err
	}

	for _, stat := range stats {
		if stat.State == "CRASHED" {
			return false, err
		}
	}

	return true, nil
}
