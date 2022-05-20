package scheduler

import (
	"zzsearcher/src/engine/types"
)

type SimpleScheduler struct {
	workChan chan types.Request
}

func (s *SimpleScheduler) Submit(
	r types.Request) {
	go func() { s.workChan <- r }()
}

func (s *SimpleScheduler) WorkerChan() chan types.Request {
	return s.workChan
}

func (s *SimpleScheduler) Run() {
	s.workChan = make(chan types.Request)
}

func (s *SimpleScheduler) WorkerReady(w chan types.Request) {
}
