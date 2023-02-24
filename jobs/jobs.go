// Package jobs provides an interface of JobRequest and implementations of it.
package jobs

import (
	"encoding/json"
	"fmt"
	"main/cache"
	"main/types"
	"sync"
	"time"
)

const jobRequestBuffer = 128

var jobListener = make(map[string][]chan types.JobResult, jobRequestBuffer)
var lock sync.RWMutex

func Request(jReq types.JobRequest) (c chan types.JobResult) {
	c = make(chan types.JobResult, 1)

	jKey := jReq.Key()

	lock.RLock()
	data, hit := cache.GetCache(jKey)
	lock.RUnlock()

	if hit {
		c <- data.Data
	} else {
		lock.Lock()
		data, hit = cache.GetCache(jKey)
		if hit {
			c <- data.Data
		} else {
			if _, ok := jobListener[jKey]; !ok {
				go execute(jReq)
			}
			jobListener[jKey] = append(jobListener[jKey], c)
		}
		lock.Unlock()
	}
	return
}

func execute(jReq types.JobRequest) {
	jKey := jReq.Key()
	jResult := jReq.Run()

	lock.Lock()
	cache.SetCache(jKey, jResult, jResult.ExpiryTime)
	for _, c := range jobListener[jKey] {
		c <- jResult
	}
	delete(jobListener, jKey)
	lock.Unlock()
}

func setDefaultHeader(jr *types.JobResult) {
	jr.StatusCode = 200
	jr.Headers = make(map[string]string)
	jr.Headers["Access-Control-Allow-Origin"] = "*"
	jr.Headers["Cache-Control"] = fmt.Sprintf("public, max-age=%d", jr.ExpiryTime)
	jr.Headers["Content-Type"] = "application/json; charset=utf-8"

	gmt := time.FixedZone("GMT", 0)
	jr.Headers["Date"] = time.Now().In(gmt).Format("Mon, 02 Jan 2006 15:04:05 MST")
}

func writeJsonResp(body any, expiryTime int64, err error) (result types.JobResult) {
	var jsonBytes []byte
	if err == nil {
		jsonBytes, err = json.Marshal(&jsonResponse{
			Data:  body,
			Error: nil,
		})
	}

	if err == nil {
		result.ExpiryTime = expiryTime
		result.Body = jsonBytes
		setDefaultHeader(&result)
	} else {
		errString := err.Error()
		jsonBytes, err = json.Marshal(&jsonResponse{
			Data:  nil,
			Error: &errString,
		})
		result.ExpiryTime = 0
		result.Body = jsonBytes
		setDefaultHeader(&result)
		result.StatusCode = 503
		// result.Body = []byte(err.Error())
	}
	return
}
