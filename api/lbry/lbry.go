// Package lbry implements API and parser of contents on LBRY.
package lbry

import (
	"errors"
	"fmt"
	"main/types"
	"main/worker"
	"strconv"
)

var workerPool *worker.WorkerPool

func Init(wp *worker.WorkerPool) {
	workerPool = wp
}

func ListVideo(channelId string) (videos []types.Video, err error) {
	results, err := claimSearchStream(channelId)
	if err != nil {
		err = errors.New("lbry: failed to list video " + err.Error())
		return
	}

	videos = make([]types.Video, 0, len(results.Result.Items))

	for _, e := range results.Result.Items {
		v := types.Video{
			Host:          "lbry",
			Type:          "video",
			VideoId:       e.CanonicalUrl,
			Thumbnail:     e.Value.Thumbnail.Url,
			Title:         e.Value.Title,
			Description:   e.Value.Description,
			PublishedTime: 0,
			Length:        e.Value.Video.Duration,
			Views:         -1,
			ChannelId:     channelId,
			ChannelTitle:  e.SigningChannel.Value.Title,
			ChannelAvatar: e.SigningChannel.Value.Thumbnail.Url,
			Sources: map[string]string{
				"lbry": e.CanonicalUrl,
			},
		}
		if timestamp, err := strconv.Atoi(e.Value.ReleaseTime); err == nil {
			v.PublishedTime = int64(timestamp)
		}
		if e.Value.Video.Width == 0 {
			v.Type = "liveStream"
		}
		videos = append(videos, v)
	}

	return
}

func SearchChannel(query string) (channels []types.Channel, err error) {
	channel, err := odyseeSearchChannel(query)
	if err != nil {
		return
	}

	channels = make([]types.Channel, 0, len(channel))
	urls := make([]string, len(channel))
	for _, e := range channel {
		urls = append(urls, fmt.Sprintf("lbry://%s#%s", e.Name, e.ClaimId))
	}

	channels, err = resolveChannels(urls)

	return
}

func YtResolve(videoIds []string, channelIds []string) (results YtResolveResult, err error) {
	response, err := ytResolve(videoIds, channelIds)
	if err != nil {
		err = errors.New("lbry: ytResolve failed " + err.Error())
		return
	}

	results = response.Data

	return
}

func YtResolveVideo(videoIds []string) (videos map[string]types.Video, err error) {
	videos = make(map[string]types.Video)

	lbryResult, err := YtResolve(videoIds, nil)

	if err == nil {
		for k, v := range lbryResult.Videos {
			videos[k] = types.Video{
				Host:    "lbry",
				VideoId: "lbry://" + v,
			}
		}
	}

	return
}

func YtResolveChannel(channelIds []string) (channels map[string]types.Channel, err error) {
	ytResolveResult, err := YtResolve(nil, channelIds)

	lbryUrls := make([]string, 0, len(channelIds))
	for _, id := range ytResolveResult.Channels {
		if len(id) > 0 {
			lbryUrls = append(lbryUrls, "lbry://"+id)
		}
	}

	lbryChannels, _ := resolveChannels(lbryUrls)
	lbryChannelsMap := make(map[string]types.Channel, len(channelIds))
	for _, c := range lbryChannels {
		lbryChannelsMap[c.ChannelId] = c
	}

	channels = make(map[string]types.Channel, len(channelIds))
	for _, y := range channelIds {
		if l, ok := lbryChannelsMap["lbry://"+ytResolveResult.Channels[y]]; ok {
			channels[y] = l
		}
	}

	return
}
