package main

import (
	"io"
	"net/http"
	"strings"
)

var hostWhitelist = []string{
	"i.ytimg.com/",
	"thumbs.odycdn.com/",
	"spee.ch/",
	"yt3.ggpht.com/",
	"yt3.googleusercontent.com",
	"thumbnails.lbry.com/",
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	if len(url) < 8 {
		w.WriteHeader(403)
		w.Write([]byte("403 Forbidden"))
		return
	}
	url = url[7:]
	isInWhitelist := false
	for _, host := range hostWhitelist {
		if strings.HasPrefix(url, host) {
			isInWhitelist = true
		}
	}
	if !isInWhitelist {
		w.WriteHeader(403)
		w.Write([]byte("403 Forbidden"))
		return
	}
	resp, err := http.Get("https://" + url)
	if err == nil {
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
