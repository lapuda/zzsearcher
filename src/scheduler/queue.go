package scheduler

import (
	"context"
	"fmt"
	"os"
	"time"
	"zzsearcher/src/engine/types"
)

type QueueScheduler struct {
	requestChan chan types.Request
	workChan    chan WorkerChan
	Ctx         *context.Context
}

func (q *QueueScheduler) Submit(r types.Request) {
	q.requestChan <- r
}

func (q *QueueScheduler) WorkerReady(w chan types.Request) {
	q.workChan <- w
}

func (q *QueueScheduler) WorkerChan() chan types.Request {
	return make(chan types.Request)
}

func (q *QueueScheduler) Run() {
	// 初始化
	q.requestChan = make(chan types.Request)
	q.workChan = make(chan WorkerChan)
	var isEmpty bool
	go func() {
		var requestQueue []types.Request
		var workQueue []WorkerChan
		for {
			var activeRequest types.Request
			var activeWorker WorkerChan
			// work空闲和Request
			if len(requestQueue) > 0 && len(workQueue) > 0 {
				activeRequest = requestQueue[0]
				activeWorker = workQueue[0]
			}
			// 多个独立的事件一起触发可能需要的是select
			select {
			case r := <-q.requestChan:
				isEmpty = false
				requestQueue = append(requestQueue, r)
			case w := <-q.workChan:
				workQueue = append(workQueue, w)
				isEmpty = false
			case activeWorker <- activeRequest: // 如果准备好的request发给了worker 就要从队列中删除
				requestQueue = requestQueue[1:]
				workQueue = workQueue[1:]
				isEmpty = false
			case <-time.Tick(5 * time.Second):
				fmt.Printf("Check quene has empty! requestQueue len is %d,workQueue is %d \n", len(requestQueue), len(workQueue))
				if len(requestQueue) == 0 {
					isEmpty = true
					fmt.Println("Quene is empty! set isEmpty is true!")
				}
				if len(requestQueue) == 0 && isEmpty == true {
					fmt.Println("Quene is empty and isEmpty has true! go dead!")
					os.Exit(0)
				}
			}
		}
	}()
}
