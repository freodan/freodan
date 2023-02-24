package jobs

import (
	"fmt"
	"main/api/lbry"
	"main/api/youtube"
	"main/types"
)

type jsonResponse struct {
	Data  any     `json:"data"`
	Error *string `json:"error"`
}

// All types here must implement interface JobRequest
//
// type JobRequest interface {
// 	Key() string
// 	Run() types.JobResult
// }

type JobRequestListVideo struct {
	Host      string
	ChannelId string
	Params    string
}

func (jr *JobRequestListVideo) Key() (key string) {
	key = fmt.Sprintf("%T_%s_%s_%s", jr, jr.Host, jr.ChannelId, jr.Params)
	return
}

func (req *JobRequestListVideo) Run() (result types.JobResult) {
	var videos []types.Video
	var err error

	if req.Host == "youtube" {
		videos, err = youtube.ListVideo(req.ChannelId, "page", req.Params)
	} else if req.Host == "lbry" {
		videos, err = lbry.ListVideo(req.ChannelId)
	}

	result = writeJsonResp(videos, 900, err) // Cache the result for 15 minutes

	return
}

type JobRequestSearchChannel struct {
	Host  string
	Query string
}

func (jr *JobRequestSearchChannel) Key() (key string) {
	key = fmt.Sprintf("%T_%s_%s", jr, jr.Host, jr.Query)
	return
}

func (req *JobRequestSearchChannel) Run() (result types.JobResult) {
	var channels []types.Channel
	var err error

	if req.Host == "youtube" {
		channels, err = youtube.SearchChannel(req.Query)
	} else if req.Host == "lbry" {
		channels, err = lbry.SearchChannel(req.Query)
	}

	result = writeJsonResp(channels, 604800, err) // Cache the result for 1 week

	return
}

type JobRequestResolveYouTubeVideo struct {
	Host     string
	VideoIds []string
}

func (jr *JobRequestResolveYouTubeVideo) Key() (key string) {
	key = fmt.Sprintf("%T_%s_%s", jr, jr.Host, jr.VideoIds)
	return
}

func (req *JobRequestResolveYouTubeVideo) Run() (result types.JobResult) {
	var videoMap map[string]types.Video
	var err error

	if req.Host == "lbry" {
		videoMap, err = lbry.YtResolveVideo(req.VideoIds)
	}

	result = writeJsonResp(videoMap, 900, err) // Cache the result for 15 minutes

	return
}

type JobRequestResolveYouTubeChannel struct {
	Host       string
	ChannelIds []string
}

func (jr *JobRequestResolveYouTubeChannel) Key() (key string) {
	key = fmt.Sprintf("%T_%s_%s", jr, jr.Host, jr.ChannelIds)
	return
}

func (req *JobRequestResolveYouTubeChannel) Run() (result types.JobResult) {
	var channelMap map[string]types.Channel
	var err error

	if req.Host == "lbry" {
		channelMap, err = lbry.YtResolveChannel(req.ChannelIds)
	}

	result = writeJsonResp(channelMap, 604800, err) // Cache the result for 1 week

	return
}
