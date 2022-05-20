package util

import (
	"fmt"
	"os"
)

const (
	CACHE_DIR = "cache"
	CACHE_EXT = ".html"
)

type Cache struct {
	URL string
}

func (c *Cache) FileKey() string {
	return MD5(c.URL)
}

func (c *Cache) GetCachePath() string {
	return CACHE_DIR + "/" + c.FileKey() + CACHE_EXT
}

func (c *Cache) IsCached() bool {
	info, err := os.Stat(c.GetCachePath())
	if info != nil {
		return true
	}
	return os.IsExist(err)
}

func (c *Cache) Cache(data []byte) {
	FileStore(c.GetCachePath(), data)
}

func (c *Cache) GetCacheContent() []byte {
	Logger.Info(fmt.Sprintf("%s Fetch From Cache %s", c.URL, c.GetCachePath()))
	return FileReStore(c.GetCachePath())
}
