package utilities

import (
	"fmt"
)

func NodeNames(instanceID string, desiredCount int) []string {
	names := make([]string, 0)

	for len(names) < desiredCount {
		names = append(names, fmt.Sprintf("service-registry-%s-%05d", instanceID, len(names)))
	}

	return names
}
