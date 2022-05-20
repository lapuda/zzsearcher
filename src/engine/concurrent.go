package engine

import (
	"zzsearcher/src/engine/types"
	"zzsearcher/src/scheduler"
)

type ConcurrentEngine struct {
	Scheduler scheduler.Scheduler
	WorkCount int
	NeedCache bool
	ItemChan  chan map[types.ParserType]interface{}
	mode      types.MODE
}

func (e *ConcurrentEngine) Run(seeds ...types.Request) {
	// 定义输出
	out := make(chan types.ParserResult)
	e.Scheduler.Run()
	// 起worker
	for i := 0; i < e.WorkCount; i++ {
		createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler, e.NeedCache)
	}
	// 提交种子
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}
	// 提交数据和提交种子
	for {
		result := <-out
		for _, item := range result.Items {
			go func() { e.ItemChan <- map[types.ParserType]interface{}{result.ParserType: item} }()
		}
		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}

}

/*创建worker*/
func createWorker(in chan types.Request, out chan types.ParserResult, ready scheduler.ReadyNotifier, needCache bool) {
	go func() {
		for {
			ready.WorkerReady(in)
			result, err := worker(<-in, needCache)
			if err != nil {
				Logger.Info(err)
				continue
			}
			out <- result
		}
	}()
}
