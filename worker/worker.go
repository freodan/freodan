// Package worker implements a HTTP client pool with proxy and captcha solver support.
package worker

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const workerRequestBuffer = 128

func NewWorker(proxyUrl string) (w *Worker, err error) {
	cookieJar, _ := cookiejar.New(nil)

	if len(proxyUrl) == 0 {
		w = &Worker{
			client: &http.Client{
				Jar: cookieJar,
			},
		}
		return
	}

	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		err = errors.New("http: invalid proxy address")
		return
	}
	transport := http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	w = &Worker{
		client: &http.Client{
			Jar:       cookieJar,
			Transport: &transport,
		},
	}

	return
}

func NewWorkerPool(proxyUrls []string) (pool *WorkerPool, err error) {
	pool = &WorkerPool{
		requests: make(chan *workerRequest, workerRequestBuffer),
	}
	for _, v := range proxyUrls {
		var w *Worker
		w, err = NewWorker(v)
		if err == nil {
			pool.workers = append(pool.workers, w)
		} else {
			err = errors.New(fmt.Sprintf("NewWorkerPool: failed to create worker %s %s", v, err.Error()))
			return
		}
	}
	pool.spawn()
	return
}
