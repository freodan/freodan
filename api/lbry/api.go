package lbry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main/types"
	"net/http"
	"net/url"
	"strings"
)

const lbryApiUrl = "https://api.na-backend.odysee.com/api/v1/proxy"

func claimSearchStream(channelId string) (results lbryClaimSearchResponse, err error) {
	urlString := lbryApiUrl
	body := lbryRequest{
		Method: "claim_search",
		Params: lbryRequestParams{
			Channel: channelId,
			// ChannelIds: []string{
			// 	channelId,
			// },
			ClaimType: []string{
				"stream",
				"repost",
			},
			HasSource: true,
			OrderBy:   "release_time",
			PageSize:  30,
			// ReleaseTime: "<1671242700",
		},
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		err = errors.New("lbry: failed to construct request body " + err.Error())
		return
	}

	resp, err := workerPool.Post(urlString, "application/json", bytes.NewReader(bodyJson))
	respBodyBytes := resp.Response.Bytes

	if err != nil {
		err = errors.New("lbry: fetch failed " + err.Error())
		return
	}

	err = json.Unmarshal(respBodyBytes, &results)

	return
}

func resolve(urls []string) (results lbryResolveResponse, err error) {
	urlString := lbryApiUrl
	body := lbryRequest{
		Method: "resolve",
		Params: lbryRequestParams{
			Urls: urls,
		},
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		err = errors.New("lbry: failed to construct request body " + err.Error())
		return
	}

	resp, err := workerPool.Post(urlString, "application/json", bytes.NewReader(bodyJson))
	respBodyBytes := resp.Response.Bytes

	if err != nil {
		err = errors.New("lbry: fetch failed " + err.Error())
		return
	}

	err = json.Unmarshal(respBodyBytes, &results)

	return
}

func resolveChannels(urls []string) (channels []types.Channel, err error) {
	channels = make([]types.Channel, 0, len(urls))

	resolve, err := resolve(urls)
	if err != nil {
		return
	}

	for _, url := range urls {
		c, ok := resolve.Result[url]
		if ok && len(c.CanonicalUrl) > 0 {
			channels = append(channels, types.Channel{
				Host:                   "lbry",
				ChannelAvatar:          c.Value.Thumbnail.Url,
				ChannelDescription:     c.Value.Description,
				ChannelId:              url,
				ChannelSubscriberCount: 0,
				ChannelTitle:           c.Value.Title,
				Alternatives:           make([]types.Channel, 0),
			})
		}
	}

	return
}

func odyseeSearchChannel(query string) (channels []odyseeSearchResult, err error) {
	urlString := fmt.Sprintf("https://lighthouse.odysee.tv/search?s=%s&size=20&from=0&claimType=channel&nsfw=false", url.QueryEscape(query))

	resp, err := workerPool.Get(urlString)
	respBodyBytes := resp.Response.Bytes
	if err != nil {
		err = errors.New("odysee: fetch failed " + err.Error())
		return
	}

	err = json.Unmarshal(respBodyBytes, &channels)

	return
}

// Auth token required
func odyseeSubCount(claimIds []string) (counts []int64, err error) {
	urlString := "https://api.odysee.com/subscription/sub_count"

	resp, err := http.Post(urlString, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("claim_id=%s", url.QueryEscape(strings.Join(claimIds, ",")))))
	if err != nil {
		err = errors.New("odysee: fetch failed " + err.Error())
		return
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("odysee: failed to read response body " + err.Error())
		return
	}

	var results odyseeSubCountResult
	err = json.Unmarshal(respBodyBytes, &results)

	if err == nil {
		counts = results.Data
	}

	return
}

func ytResolve(videoIds []string, channelIds []string) (results lbryYtResolveResponse, err error) {
	urlString := fmt.Sprintf("https://api.lbry.com/yt/resolve?video_ids=%s&channel_ids=%s", strings.Join(videoIds, ","), strings.Join(channelIds, ","))

	resp, err := workerPool.Get(urlString)
	respBodyBytes := resp.Response.Bytes
	if err != nil {
		err = errors.New("lbry: fetch failed " + err.Error())
		return
	}

	err = json.Unmarshal(respBodyBytes, &results)
	return
}
