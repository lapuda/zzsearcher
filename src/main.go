package main

import (
	"context"
	"flag"
	"strconv"
	"zzsearcher/src/engine"
	"zzsearcher/src/engine/types"
	"zzsearcher/src/fatcher"
	"zzsearcher/src/parser"
	"zzsearcher/src/scheduler"
	"zzsearcher/src/store"
)

func main() {
	var num = flag.Int("conc_num", 10, "-conc_num 10")
	flag.Parse()
	//en := engine.SimpleEngine{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	en := engine.ConcurrentEngine{
		Scheduler: &scheduler.QueueScheduler{Ctx: &ctx},
		WorkCount: *num,
		NeedCache: true,
		//ItemChan:  store.ItemSaver(ctx), //mongodb
		ItemChan: store.ItemRemote(ctx), //remote
		//ItemChan: store.ItemPrinter(ctx), // 打印
	}
	//request := types.Request{
	//	Url: "https://m.xianman123.com/wuliandianfeng/",
	//	ParseFunc: parser.MangaParser,
	//	FetchFunc: fatcher.Fetch,
	//}

	//request := types.Request{
	//	Url: "https://m.xianman123.com/f-1-0-0-0-0-1-1.html",
	//	ParseFunc: parser.SeedListParser,
	//	FetchFunc: fatcher.Fetch,
	//}

	en.Run(makeSeed(10)...)
	//en.Run(makeDetailSeed())
	//en.Run(makeChapterSeed())
}

func makeDetailSeed() types.Request {
	return types.Request{

		Url:       "http://www.kumw5.com/16041/",
		ParseFunc: parser.MangaParser,
		FetchFunc: fatcher.Fetch,
	}
}

func makeChapterSeed() types.Request {
	return types.Request{
		Url:       "http://www.kumw5.com/mulu/24538/1-1.html",
		ParseFunc: parser.ChapterListParser,
		FetchFunc: fatcher.Fetch,
	}
}

func makeSeed(num int) []types.Request {
	var seed []types.Request
	for i := 1; i <= num; i++ {
		url := "http://www.kumw5.com/sort/13-" + strconv.Itoa(i) + ".html"
		seed = append(seed, types.Request{
			Url:       url,
			ParseFunc: parser.MangaListParser,
			FetchFunc: fatcher.Fetch,
		})
	}
	return seed
}
