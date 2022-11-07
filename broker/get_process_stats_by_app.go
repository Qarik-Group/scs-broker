package broker

import (
	"errors"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/lager"
	cf "github.com/cloudfoundry-community/go-cfclient"
)

func getProcessStatsByAppAndType(cfClient *ccv3.Client, community *cf.Client, logger lager.Logger, appGUID string, procType string) ([]cf.Stats, error) {
	stats := make([]cf.Stats, 0)

	procs, err := getApplicationProcessesByType(cfClient, logger, appGUID, procType)
	if err != nil {
		return stats, err
	}

	for _, proc := range procs {
		candidates, err := community.GetProcessStats(proc.GUID)
		if err != nil {
			continue
		}

		for _, stat := range candidates {
			stats = append(stats, stat)
		}
	}

	if len(stats) == 0 {
		return stats, errors.New("no stats found")
	}

	return stats, nil

}
