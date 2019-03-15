package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Cache struct {
	LastModified string
	Values       *map[string]string
}

func (b *BotConfig) getCacheFileName(name string) string {
	return fmt.Sprintf("%s/%s", b.CacheDir, name)
}

func (b *BotConfig) initCache() {
	tempDir, err := ioutil.TempDir("", "bipboupcache")
	if err != nil {
		panic(err)
	}

	b.CacheDir = tempDir
}

func (b *BotConfig) createCacheIfNeeded() {
	fileInfo, err := os.Stat(b.CacheDir)
	if err != nil || !fileInfo.IsDir() {
		b.ClearCache(true)
	}
}

func (b *BotConfig) ClearCache(leaveDirOpt ...bool) {
	os.RemoveAll(b.CacheDir)
	if len(leaveDirOpt) > 0 && leaveDirOpt[0] {
		b.initCache()
	}
}

func (b *BotConfig) GetCache(name string) (cache *Cache) {
	err := json.Unmarshal(b.GetCacheRaw(name), &cache)
	if err != nil {
		cache = nil
	}

	return
}

func (b *BotConfig) GetCacheRaw(name string) []byte {
	fileHandler, err := os.Open(b.getCacheFileName(name))
	defer fileHandler.Close()
	if err != nil {
		return []byte{}
	}

	bytes, err := ioutil.ReadAll(fileHandler)
	return bytes
}

func (b *BotConfig) SetCache(name string, lastModified string, values *map[string]string) {
	// First we check the cacheDir still exists. If not, we create a new one.
	b.createCacheIfNeeded()

	fileHandler, err := os.OpenFile(b.getCacheFileName(name), os.O_CREATE|os.O_WRONLY, 0644)
	defer fileHandler.Close()
	if err != nil {
		return
	}

	encoder := json.NewEncoder(fileHandler)
	encoder.Encode(&Cache{LastModified: lastModified, Values: values})
}

func (b *BotConfig) SetCacheRaw(name string, values []byte) {
	// First we check the cacheDir still exists. If not, we create a new one.
	b.createCacheIfNeeded()

	fileHandler, err := os.OpenFile(b.getCacheFileName(name), os.O_CREATE|os.O_WRONLY, 0644)
	defer fileHandler.Close()
	if err != nil {
		return
	}

	fileHandler.Write(values)
}
