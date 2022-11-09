package broker

import (
	"errors"

	cf "github.com/cloudfoundry-community/go-cfclient"
)

func (broker *SCSBroker) getProcessStatsByAppAndType(appGUID string, procType string) ([]cf.Stats, error) {
	stats := make([]cf.Stats, 0)

	community, err := broker.GetCommunity()
	if err != nil {
		return stats, err
	}

	procs, err := broker.getApplicationProcessesByType(appGUID, procType)
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
