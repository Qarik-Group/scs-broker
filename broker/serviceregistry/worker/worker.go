package worker

import (
	"github.com/gammazero/workerpool"
)

func New() *workerpool.WorkerPool {
	return workerpool.New(5)
}
