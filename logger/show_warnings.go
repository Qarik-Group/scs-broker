package logger

import (
	"fmt"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/lager"
)

func ShowWarnings(warnings ccv3.Warnings, subject interface{}) {
	Info(
		fmt.Sprintf(
			"NOTICE: %d warning(s) were detected!",
			len(warnings),
		),
		lager.Data{"Subject": subject},
	)

	for warn := range warnings {
		w := warnings[warn]
		Info(fmt.Sprintf("Warning(#%d): %s ", warn+1, w))
	}
}
