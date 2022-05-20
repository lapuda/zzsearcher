package fatcher

import (
	"errors"
	"fmt"
	"github.com/jeanphorn/log4go"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"time"
	"zzsearcher/src/util"
)

var Logger = log4go.NewDefaultLogger(log4go.FINE)

var UA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_5_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36"

func SkipFetcher(url string, needCache bool) ([]byte, []byte, error) {
	return nil, nil, errors.New(fmt.Sprintf("Fetch URL %s is skiped!", url))
}

var rateLimit = time.Tick(200 * time.Millisecond)

// 返回
func Fetch(url string, needCache bool) ([]byte, []byte, error) {
	var all []byte
	// 判断缓存
	cache := util.Cache{url}
	if cache.IsCached() && needCache {
		all = cache.GetCacheContent()
	} else {
		<-rateLimit
		cookieJar, _ := cookiejar.New(nil)
		client := &http.Client{Jar: cookieJar}
		request, _ := http.NewRequest("GET", url, nil)
		//request.Header.Add("Referer","https://www.cocomanhua.com/17449/")

		request.Header.Add("User-Agent", UA)
		var clist []*http.Cookie
		clist = append(clist, &http.Cookie{
			Name: "__",
		})
		cookieJar.SetCookies(request.URL, clist)

		resp, err := client.Do(request)
		if err != nil {
			return nil, nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, nil, fmt.Errorf("Request Body Exception! code is : %d", resp.StatusCode)
		}
		all, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}
	}

	// 缓存
	if needCache && !cache.IsCached() {
		cache.Cache(all)
	}

	/* 获取标题 */
	re := regexp.MustCompile("<title>([^>]*)</title>")
	matches := re.FindAllSubmatch(all, -1)
	if len(matches) < 1 {
		return nil, nil, errors.New("no title")
	}
	/* 获取标题 end */
	return all, matches[0][1], nil
}

func FetchBrowser(url string) ([]byte, []byte, error) {
	const (
		browseDriverPath = "browse_driver/chromedriver"
		port             = 9515
	)
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService(browseDriverPath, port, opts...)
	if nil != err {
		return nil, nil, err
	}
	defer service.Stop()

	//链接本地的浏览器 chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	//禁止图片加载，加快渲染速度
	imgCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
		"profile.managed_default_content_settings.css":    2,
		//"profile.managed_default_content_settings.javascript": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imgCaps,
		Path:  "",
		Args: []string{
			"--headless",
			"--start-maximized",
			//"--window-size=1200x60000",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
			"--disable-gpu",
			"--disable-impl-side-painting",
			"--disable-gpu-sandbox",
			"--disable-accelerated-2d-canvas",
			"--disable-accelerated-jpeg-decoding",
			//"--test-type=ui",
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)
	// 调起chrome浏览器
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		fmt.Println("connect to the webDriver faild", err.Error())
		panic(err)
	}
	defer wd.Quit()
	if err := wd.Get(url); err != nil {
		panic(err)
	}

	//e,err :=wd.FindElement(selenium.ByCSSSelector,"div.show-more")
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	if err1 := e.Click();err1 == nil {
	//		time.Sleep(5 * time.Second)
	//	}
	//}
	wd.ExecuteScript("javascriptshowPic();", nil)
	wd.SetPageLoadTimeout(time.Duration(time.Second * 20))

	//wd.ResizeWindow("",1000,10000)

	content, err := wd.PageSource()

	filename := "tmp/" + util.MD5(url) + ".html"
	util.FileStore(filename, []byte(content))
	//os.Exit(0)
	if err != nil {
		panic(err)
	}

	/* 获取标题 */
	re := regexp.MustCompile("<title>([^>]*)</title>")
	matches := re.FindAllSubmatch([]byte(content), -1)
	/* 获取标题 end */

	return []byte(content), matches[0][1], nil
}

func FetcherFunc(funcName string) func(string, bool) ([]byte, []byte, error) {
	switch funcName {
	case "simple_fetcher":
		return Fetch
	default:
		return SkipFetcher
	}
}
