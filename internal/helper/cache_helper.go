package helper

import (
	"fmt"

	"github.com/bluele/gcache"
)

type CacheHelper struct {
	gc gcache.Cache
}

func NewCacheHelper() *CacheHelper {
	gc := gcache.New(3).
		Build()
	return &CacheHelper{gc: gc}
}

// Untuk mengambil value dari dalam cache
func (ch CacheHelper) Get(key string) string {
	isKeyExists := ch.gc.Has(key)
	if !isKeyExists {
		PrintLog(fmt.Sprintf("Key %s belum tersedia didalam cache", key))
		return ""
	}

	value, err := ch.gc.Get(key)
	if err != nil {
		PrintLog(fmt.Sprintf("Gagal mengambil key %s dari dalam cache", key))
	}
	return fmt.Sprint(value)
}

// Untuk mengisi value kedalam cache
func (ch CacheHelper) Set(key string, value string) bool {
	err := ch.gc.Set(key, value)
	if err != nil {
		PrintLog(fmt.Sprintf("Gagal menyimpan key %s dengan value %s kedalam cache", key, value))
		return false
	}
	return true
}
