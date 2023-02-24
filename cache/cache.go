// Package cache implements a simple in-memory cache with TTL.
//
// This package is NOT thread-safe. Concurrent read and write might produce runtime error.
package cache

import (
	"main/types"
	"time"
)

const maxNumCacheEntries = 1000000

type CacheEntry struct {
	Data       types.JobResult
	CreateTime int64 // Unix timestamp
	ExpiryTime int64 // How long the cached result should be served in seconds.
}

var cache = make(map[string]CacheEntry, maxNumCacheEntries)
var cacheQueue = make([]string, 0, maxNumCacheEntries)

func GetCache(key string) (data CacheEntry, hit bool) {
	data, hit = cache[key]
	if time.Now().Unix() >= data.ExpiryTime {
		hit = false
		delete(cache, key)
	}
	return
}

func SetCache(key string, data types.JobResult, expiry int64) {
	if len(cacheQueue) >= maxNumCacheEntries {
		freeCache(1)
	}

	timestamp := time.Now().Unix()
	cache[key] = CacheEntry{
		Data:       data,
		CreateTime: timestamp,
		ExpiryTime: timestamp + expiry,
	}
	cacheQueue = append(cacheQueue, key)
	return
}

func freeCache(size int) {
	if len(cacheQueue) < size {
		size = len(cacheQueue)
	}
	for i := 0; i < size; i++ {
		delete(cache, cacheQueue[i])
	}
	cacheQueue = cacheQueue[size:]
}
