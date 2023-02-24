package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/types"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const youtubeHost = "https://www.youtube.com"

func getChannel(url string) (channel []types.Channel, err error) {
	matched, err := regexp.MatchString(`^https?:\/\/(www\.)?youtube\.com\/((@[^\/]*)|(c\/[^\/]*)|(channel\/[^\/]*))`, url)
	if !matched || err != nil {
		err = errors.New("youtube: invalid url")
		return
	}

	r, err := workerPool.Get(url)
	respBody := r.Response.String

	if err != nil {
		err = errors.New("youtube: fetch failed " + err.Error())
		return
	}

	var jsonString string

	jsonString, err = extractJson(respBody)
	if err != nil {
		return
	}

	var data initialData
	json.Unmarshal([]byte(jsonString), &data)

	var channelAvatar = ""
	var channelId = data.Metadata.ChannelMetadataRenderer.ExternalId
	var channelTitle = data.Metadata.ChannelMetadataRenderer.Title
	var channelDescription = data.Metadata.ChannelMetadataRenderer.Description
	var channelSubscriberCount = countBuilder(data.Header.C4TabbedHeaderRenderer.SubscriberCountText.SimpleText)

	if len(data.Metadata.ChannelMetadataRenderer.Avatar.Thumbnails) > 0 {
		channelAvatar = data.Metadata.ChannelMetadataRenderer.Avatar.Thumbnails[len(data.Metadata.ChannelMetadataRenderer.Avatar.Thumbnails)-1].Url
	}

	channel = []types.Channel{
		{
			Host:                   "youtube",
			ChannelAvatar:          channelAvatar,
			ChannelDescription:     channelDescription,
			ChannelId:              channelId,
			ChannelSubscriberCount: channelSubscriberCount,
			ChannelTitle:           channelTitle,
			Alternatives:           make([]types.Channel, 0),
		},
	}

	return
}

func getVideos(channelId string, method string, params string) (videos []types.Video, err error) {
	var url string

	if method == "page" {
		url = fmt.Sprintf(youtubeHost+"/channel/%s/%s", channelId, params)
	} else if method == "rss" {
		url = fmt.Sprintf(youtubeHost+"/feeds/videos.xml?channel_id=%s", channelId)
		err = errors.New("youtube: method \"rss\" is not implemented")
		return
	} else if method == "api" {
		err = errors.New("youtube: method \"api\" is not implemented")
		return
	} else {
		err = errors.New("youtube: unknown method")
		return
	}

	r, err := workerPool.Get(url)
	respBody := r.Response.String

	if err != nil {
		err = errors.New("youtube: fetch failed " + err.Error())
		return
	}

	if method == "page" || method == "api" {
		var jsonString string

		if method == "page" {
			jsonString, err = extractJson(respBody)
			if err != nil {
				return
			}
		} else {
			jsonString = respBody
		}

		var data initialData
		err = json.Unmarshal([]byte(jsonString), &data)
		if err != nil {
			err = errors.New("youtube: failed to unmarshal json " + err.Error())
			return
		}

		videos, err = initialDataParser(data)
	}

	return
}

func searchChannel(query string) (channels []types.Channel, err error) {
	urlString := fmt.Sprintf(youtubeHost+"/results?sp=%s&search_query=%s", "EgIQAg==", url.QueryEscape(query))

	r, err := workerPool.Get(urlString)
	respBody := r.Response.String

	if err != nil {
		err = errors.New("youtube: fetch failed " + err.Error())
		return
	}

	var jsonString string

	jsonString, err = extractJson(respBody)
	if err != nil {
		return
	}

	var data searchInitialData
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		err = errors.New("youtube: failed to unmarshal json " + err.Error())
		return
	}

	if len(data.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents) > 0 {
		contents := data.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents[0].ItemSectionRenderer.Contents
		channels = make([]types.Channel, 0, len(contents))

		for i := 0; i < len(contents); i++ {
			if len(contents[i].ChannelRenderer.ChannelId) < 1 {
				continue
			}

			var channelAvatar = ""
			var channelId = contents[i].ChannelRenderer.ChannelId
			var channelTitle = contents[i].ChannelRenderer.Title.SimpleText
			var channelDescription = ""
			var channelSubscriberCount = countBuilder(contents[i].ChannelRenderer.VideoCountText.SimpleText)

			if len(contents[i].ChannelRenderer.Thumbnail.Thumbnails) > 0 {
				channelAvatar = contents[i].ChannelRenderer.Thumbnail.Thumbnails[len(contents[i].ChannelRenderer.Thumbnail.Thumbnails)-1].Url
			}

			for _, e := range contents[i].ChannelRenderer.DescriptionSnippet.Runs {
				channelDescription += e.Text
			}

			channel := types.Channel{
				Host:                   "youtube",
				ChannelAvatar:          channelAvatar,
				ChannelDescription:     channelDescription,
				ChannelId:              channelId,
				ChannelSubscriberCount: channelSubscriberCount,
				ChannelTitle:           channelTitle,
				Alternatives:           make([]types.Channel, 0),
			}

			channels = append(channels, channel)
		}
	}

	return
}

func initialDataParser(data initialData) (videos []types.Video, err error) {
	var contents []contents

	var channelAvatar = ""
	var channelId = data.Metadata.ChannelMetadataRenderer.ExternalId
	var channelTitle = data.Metadata.ChannelMetadataRenderer.Title

	if len(data.Metadata.ChannelMetadataRenderer.Avatar.Thumbnails) > 0 {
		channelAvatar = data.Metadata.ChannelMetadataRenderer.Avatar.Thumbnails[len(data.Metadata.ChannelMetadataRenderer.Avatar.Thumbnails)-1].Url
	}

	for i := 0; i < len(data.Contents.TwoColumnBrowseResultsRenderer.Tabs); i++ {
		if data.Contents.TwoColumnBrowseResultsRenderer.Tabs[i].TabRenderer.Selected == true {
			contents = data.Contents.TwoColumnBrowseResultsRenderer.Tabs[i].TabRenderer.Content.RichGridRenderer.Contents
			videos = make([]types.Video, 0, len(contents))
		}
	}

	for i := 0; i < len(contents); i++ {
		if len(contents[i].RichItemRenderer.Content.VideoRenderer.VideoId) > 0 {
			video, err := videoParser(contents[i].RichItemRenderer.Content.VideoRenderer)
			if err != nil {
				continue
			}

			video.ChannelId = channelId
			video.ChannelTitle = channelTitle
			video.ChannelAvatar = channelAvatar

			videos = append(videos, video)
		}

		if len(contents[i].RichItemRenderer.Content.ReelItemRenderer.VideoId) > 0 {
			video, err := reelParser(contents[i].RichItemRenderer.Content.ReelItemRenderer)
			if err != nil {
				continue
			}

			video.ChannelId = channelId
			video.ChannelTitle = channelTitle
			video.ChannelAvatar = channelAvatar

			videos = append(videos, video)
		}
	}

	return
}

func reelParser(reel reelItemRenderer) (video types.Video, err error) {
	video.Host = "youtube"
	video.Type = "short"
	video.VideoId = reel.VideoId
	video.Thumbnail = ""
	video.Title = reel.Headline.SimpleText
	video.Description = ""
	video.PublishedTime = 0
	video.Length = 0
	video.Views = 0
	video.Sources = make(map[string]string)
	video.Sources["youtube"] = reel.VideoId

	if len(reel.Thumbnail.Thumbnails) > 0 {
		video.Thumbnail = reel.Thumbnail.Thumbnails[len(reel.Thumbnail.Thumbnails)-1].Url
	}

	video.Views = countBuilder(reel.ViewCountText.SimpleText)

	// Sometime YouTube doesn't provide publish timestamp. Could be AB testing.
	if len(reel.NavigationEndpoint.ReelWatchEndpoint.Overlay.ReelPlayerOverlayRenderer.ReelPlayerHeaderSupportedRenderers.ReelPlayerHeaderRenderer.TimestampText.SimpleText) > 0 {
		video.PublishedTime = videoPublishedTimeBuilder(reel.NavigationEndpoint.ReelWatchEndpoint.Overlay.ReelPlayerOverlayRenderer.ReelPlayerHeaderSupportedRenderers.ReelPlayerHeaderRenderer.TimestampText.SimpleText)
	} else {
		video.PublishedTime = -1
	}

	return
}

func videoParser(vr videoRenderer) (video types.Video, err error) {
	video.Host = "youtube"
	video.Type = "video"
	video.VideoId = vr.VideoId
	video.Thumbnail = ""
	video.Title = ""
	video.Description = ""
	video.PublishedTime = 0
	video.Length = 0
	video.Views = 0
	video.Sources = make(map[string]string)
	video.Sources["youtube"] = vr.VideoId

	if len(vr.Thumbnail.Thumbnails) > 0 {
		video.Thumbnail = vr.Thumbnail.Thumbnails[len(vr.Thumbnail.Thumbnails)-1].Url
	}
	if len(vr.Title.Runs) > 0 {
		video.Title = vr.Title.Runs[0].Text
	}
	if len(vr.DescriptionSnippet.Runs) > 0 {
		video.Description = vr.DescriptionSnippet.Runs[0].Text
	}

	// View count builder
	if views, err := strconv.Atoi(stripText2Number(vr.ViewCountText.SimpleText)); err == nil {
		video.Views = int64(views)
	}

	video.Length = videoLengthBuilder(vr.LengthText.SimpleText)
	video.PublishedTime = videoPublishedTimeBuilder(vr.PublishedTimeText.SimpleText)

	// Video type
	publishedTimeSlice := strings.Split(vr.PublishedTimeText.SimpleText, " ")
	if len(publishedTimeSlice) > 0 && strings.HasPrefix(publishedTimeSlice[0], "Streamed") {
		video.Type = "stream"
	}

	if len(vr.ViewCountText.Runs) > 1 {
		if vr.ThumbnailOverlays[0].ThumbnailOverlayTimeStatusRenderer.Style == "LIVE" {
			video.Type = "liveStream"
			if l, err := strconv.ParseInt(stripText2Number(vr.ViewCountText.Runs[0].Text), 10, 64); err == nil {
				video.Views = l
			}
		}
	}

	if len(vr.ThumbnailOverlays) > 1 {
		if strings.Contains(vr.ThumbnailOverlays[0].ThumbnailOverlayTimeStatusRenderer.Text.SimpleText, "PREMIERE") {
			video.Type = "livePremiere"
			if len(vr.ViewCountText.Runs) > 0 {
				if l, err := strconv.ParseInt(stripText2Number(vr.ViewCountText.Runs[0].Text), 10, 64); err == nil {
					video.Views = l
				}
			}
		}
	}

	if startTime, err := strconv.ParseInt(vr.UpcomingEventData.StartTime, 10, 64); err == nil && startTime > 0 {
		video.PublishedTime = startTime

		if len(vr.UpcomingEventData.UpcomingEventText.Runs) > 0 {
			if strings.HasPrefix(vr.UpcomingEventData.UpcomingEventText.Runs[0].Text, "Premiere") {
				video.Type = "upcomingPremiere"
			} else if strings.HasPrefix(vr.UpcomingEventData.UpcomingEventText.Runs[0].Text, "Scheduled") {
				video.Type = "upcomingStream"
			}
		}
	}

	return
}

func extractJson(html string) (jsonString string, err error) {
	const substrStart = ">var ytInitialData = {"
	const substrEnd = ";</script>"
	dataStart := strings.Index(html, substrStart)
	if dataStart < 0 {
		err = errors.New("youtube: ytInitialData not found")
		return
	}
	dataStart = dataStart + len(substrStart) - 1
	dataEnd := strings.Index(html[dataStart:], substrEnd)
	if dataEnd < 0 {
		err = errors.New("youtube: ytInitialData not found")
		return
	}
	jsonString = html[dataStart : dataStart+dataEnd]
	return
}

// Remove all non-numeric symbol in the string.
func stripText2Number(input string) (output string) {
	var sb strings.Builder
	for i := 0; i < len(input); i++ {
		if input[i] >= '0' && input[i] <= '9' || input[i] == '.' {
			sb.WriteByte(input[i])
		}
	}
	output = sb.String()
	return
}

func countBuilder(input string) (count int64) {
	inputSlice := strings.Split(input, " ")
	countUnits := map[byte]float64{
		'K': 1000.0,
		'M': 1000000.0,
		'B': 1000000000.0,
	}
	if c, err := strconv.ParseFloat(stripText2Number(inputSlice[0]), 64); err == nil {
		if len(inputSlice[0]) > 1 {
			if s, ok := countUnits[inputSlice[0][len(inputSlice[0])-1]]; ok == true {
				count = int64(c * s)
			}
		}
	}
	return
}

func videoLengthBuilder(input string) (lengthInSeconds int) {
	lengthInSeconds = 0
	lengthUnitsInSeconds := []int{1, 60, 3600}
	lengthUnits := strings.Split(input, ":")
	for j := len(lengthUnits) - 1; j >= 0; j-- {
		if l, err := strconv.Atoi(lengthUnits[j]); err == nil {
			lengthInSeconds += l * lengthUnitsInSeconds[len(lengthUnits)-j-1]
		}
	}
	return
}

func videoPublishedTimeBuilder(input string) (timestamp int64) {
	unixTimeNow := time.Now().Unix()
	publishedTimeUnitsInSeconds := map[string]int{
		"second":  1,
		"seconds": 1,
		"minute":  60,
		"minutes": 60,
		"hour":    3600,
		"hours":   3600,
		"day":     86400,
		"days":    86400,
		"week":    604800,
		"weeks":   604800,
		"month":   2592000,
		"months":  2592000,
		"year":    31536000,
		"years":   31536000,
	}
	publishedTimeSlice := strings.Split(input, " ")
	publishedTimeOffset := 0
	if len(publishedTimeSlice) >= 3 {
		if l, err := strconv.Atoi(publishedTimeSlice[len(publishedTimeSlice)-3]); err == nil {
			if s, ok := publishedTimeUnitsInSeconds[publishedTimeSlice[len(publishedTimeSlice)-2]]; ok == true {
				publishedTimeOffset += l * s
			}
		}
	}
	timestamp = unixTimeNow - int64(publishedTimeOffset)
	return
}
