package utilities

import (
	"fmt"
	"strings"
)

type fallable interface {
	Failure() error
}

func CollectedErrors(nodes []fallable, msg string) error {
	c := make([]string, 0)

	for _, node := range nodes {
		err := node.Failure()
		if err != nil {
			c = append(c, err.Error())
		}
	}

	if len(c) > 0 {
		return fmt.Errorf("%s: %s", msg, strings.Join(c, ","))
	}

	return nil
}
