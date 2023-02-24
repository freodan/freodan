package worker

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Worker struct {
	client *http.Client
	lock   sync.Mutex
}

func (w *Worker) Do(req *http.Request) (resp *WorkerResponse) {
	w.lock.Lock()
	defer w.lock.Unlock()

	var r *http.Response
	resp = &WorkerResponse{}

	r, resp.Error = w.client.Do(req)
	if resp.Error != nil {
		resp.Error = errors.New("client: get failed " + resp.Error.Error())
		return
	}

	resp.Header = &r.Header
	resp.Request = r.Request
	resp.Response.Bytes, resp.Error = io.ReadAll(r.Body)
	if resp.Error != nil {
		resp.Error = errors.New("client: failed to read response body " + resp.Error.Error())
		return
	}

	currentUrl := resp.Request.URL
	resp.Response.String = string(resp.Response.Bytes)

	if strings.HasSuffix(currentUrl.Host, "google.com") && currentUrl.Path == "/sorry/index" {
		isRecaptcha, siteKey, dataS, q := isRecaptcha(resp.Response.String)
		fmt.Println(siteKey, dataS, q)
		if isRecaptcha {
			var gRecaptchaResponse string
			gRecaptchaResponse, resp.Error = solveRecaptcha("https://www.google.com/sorry/index", siteKey, dataS)
			if resp.Error != nil {
				resp.Error = errors.New("client: failed to solve recaptcha " + resp.Error.Error())
				return
			}

			r, resp.Error = w.client.Post(
				"https://"+currentUrl.Host+currentUrl.Path,
				"application/x-www-form-urlencoded",
				strings.NewReader(
					fmt.Sprintf(
						"g-recaptcha-response=%s&q=%s&continue=%s",
						url.QueryEscape(gRecaptchaResponse),
						url.QueryEscape(q),
						url.QueryEscape(req.URL.String()),
					),
				),
			)

			resp.Header = &r.Header
			resp.Request = r.Request
			resp.Response.Bytes, resp.Error = io.ReadAll(r.Body)
			if resp.Error != nil {
				resp.Error = errors.New("client: failed to read response body " + resp.Error.Error())
				return
			}

			resp.Response.String = string(resp.Response.Bytes)
		} else {
			resp.Error = errors.New("client: unknown response")
		}
	}
	return
}

type workerRequest struct {
	request         *http.Request
	responseChannel chan *WorkerResponse
}

type WorkerResponse struct {
	Error    error
	Header   *http.Header
	Request  *http.Request
	Response struct {
		String string
		Bytes  []byte
	}
}

type WorkerPool struct {
	requests chan *workerRequest
	workers  []*Worker
}

func (wp *WorkerPool) Get(url string) (resp *WorkerResponse, err error) {
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = errors.New("workerpool: failed to create http request " + err.Error())
		return
	}

	rc := make(chan *WorkerResponse)

	request := &workerRequest{
		request:         httpReq,
		responseChannel: rc,
	}

	wp.requests <- request

	resp = <-rc

	err = resp.Error

	return
}

func (wp *WorkerPool) Post(url string, contentType string, body io.Reader) (resp *WorkerResponse, err error) {
	httpReq, err := http.NewRequest("POST", url, body)
	if err != nil {
		err = errors.New("workerpool: failed to create http request " + err.Error())
		return
	}

	httpReq.Header.Set("Content-Type", contentType)

	rc := make(chan *WorkerResponse)

	request := &workerRequest{
		request:         httpReq,
		responseChannel: rc,
	}

	wp.requests <- request

	resp = <-rc

	err = resp.Error

	return
}

func (wp *WorkerPool) spawn() {
	for _, w := range wp.workers {
		go func(wp *WorkerPool, w *Worker) {
			for {
				r, ok := <-wp.requests
				if !ok {
					break
				}
				r.responseChannel <- w.Do(r.request)
			}
		}(wp, w)
	}
}
