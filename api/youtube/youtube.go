// Package youtube implements API and parser of contents on YouTube.
package youtube

import (
	"main/types"
	"main/worker"
	"regexp"
)

var workerPool *worker.WorkerPool

func Init(wp *worker.WorkerPool) {
	workerPool = wp
}

func ListVideo(channelId string, method string, params string) (videos []types.Video, err error) {
	videos, err = getVideos(channelId, method, params)
	return
}

func SearchChannel(query string) (channels []types.Channel, err error) {
	matched, err := regexp.MatchString(`^https?:\/\/(www\.)?youtube\.com\/((@[^\/]*)|(c\/[^\/]*)|(channel\/[^\/]*))`, query)

	if matched {
		channels, err = getChannel(query)
	} else {
		channels, err = searchChannel(query)
	}

	return
}
