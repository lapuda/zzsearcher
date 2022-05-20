package scheduler

import (
	"zzsearcher/src/engine/types"
)

type WorkerChan chan types.Request

type Scheduler interface {
	ReadyNotifier
	Submit(request types.Request)
	WorkerChan() chan types.Request
	Run()
}

type ReadyNotifier interface {
	WorkerReady(chan types.Request)
}
