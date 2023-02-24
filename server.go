package main

import (
	"fmt"
	"log"
	"main/api/lbry"
	"main/api/youtube"
	"main/jobs"
	"main/types"
	"main/worker"
	"net/http"
	"strings"
)

func main() {
	conf, err := readConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	lbryProxy := make([]string, 0)
	for _, v := range conf.Worker.Lbry {
		fmt.Printf(`[CONFIG] Added Lbry worker. Proxy: "%s", Concurrent: %d`+"\n", v.Proxy, v.Concurrent)
		for i := 0; i < v.Concurrent; i++ {
			lbryProxy = append(lbryProxy, v.Proxy)
		}
	}

	lbryWp, err := worker.NewWorkerPool(lbryProxy)
	if err != nil {
		log.Fatal(err)
	}
	lbry.Init(lbryWp)

	youtubeProxy := make([]string, 0)
	for _, v := range conf.Worker.Youtube {
		fmt.Printf(`[CONFIG] Added YouTube worker. Proxy: "%s", Concurrent: %d`+"\n", v.Proxy, v.Concurrent)
		for i := 0; i < v.Concurrent; i++ {
			youtubeProxy = append(youtubeProxy, v.Proxy)
		}
	}

	youtubeWp, err := worker.NewWorkerPool(youtubeProxy)
	if err != nil {
		log.Fatal(err)
	}
	youtube.Init(youtubeWp)

	cs := make([]worker.CaptchaSolver, 0)
	for _, v := range conf.CaptchaSolver {
		if v.Service == "2captcha" {
			cs = append(cs, &worker.TwoCaptcha{
				Host: v.Host,
				Key:  v.Key,
			})
		} else {
			fmt.Printf(`[CONFIG] Unsupported captcha solving service: %s`, v.Service)
			continue
		}
		fmt.Printf(`[CONFIG] Added captcha solving service: %s`+"\n", v.Service)
	}
	worker.CaptchaSolverInit(cs)

	fmt.Printf(`[CONFIG] HTTP server listening at: "%s"`+"\n", conf.Listen)

	http.HandleFunc("/api/search/channel", searchChannelHandler)
	http.HandleFunc("/api/list/video", listVideoHandler)
	http.HandleFunc("/api/ytResolve/video", ytResolveVideoHandler)
	http.HandleFunc("/api/ytResolve/channel", ytResolveChannelHandler)
	http.HandleFunc("/proxy/", proxyHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	err = http.ListenAndServe(conf.Listen, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func listVideoHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	channelId := r.URL.Query().Get("channel_id")
	params := r.URL.Query().Get("params")

	jReq := jobs.JobRequestListVideo{
		Host:      host,
		ChannelId: channelId,
		Params:    params,
	}

	jResult := <-jobs.Request(&jReq)

	writeResponse(w, jResult)
}

func ytResolveVideoHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	videoIds := r.URL.Query().Get("videoIds")

	jReq := jobs.JobRequestResolveYouTubeVideo{
		Host:     host,
		VideoIds: strings.Split(videoIds, ","),
	}

	jResult := <-jobs.Request(&jReq)

	writeResponse(w, jResult)
}

func ytResolveChannelHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	channelIds := r.URL.Query().Get("channelIds")

	jReq := jobs.JobRequestResolveYouTubeChannel{
		Host:       host,
		ChannelIds: strings.Split(channelIds, ","),
	}

	jResult := <-jobs.Request(&jReq)

	writeResponse(w, jResult)
}

func searchChannelHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	query := r.URL.Query().Get("query")

	jReq := jobs.JobRequestSearchChannel{
		Host:  host,
		Query: query,
	}

	jResult := <-jobs.Request(&jReq)

	writeResponse(w, jResult)
}

func writeResponse(w http.ResponseWriter, r types.JobResult) {
	for k, v := range r.Headers {
		w.Header().Add(k, v)
	}
	w.WriteHeader(r.StatusCode)
	w.Write(r.Body)
}
