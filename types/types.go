// Package types provides types shared across the project.
package types

type JobRequest interface {
	Key() string
	Run() JobResult
}

type JobResult struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
	ExpiryTime int64 // How long the cached result should be served in seconds.
}

type Channel struct {
	Host                   string    `json:"host"`
	ChannelAvatar          string    `json:"channelAvatar"`
	ChannelDescription     string    `json:"channelDescription"`
	ChannelId              string    `json:"channelId"`
	ChannelSubscriberCount int64     `json:"channelSubscriberCount"`
	ChannelTitle           string    `json:"channelTitle"`
	Alternatives           []Channel `json:"alternatives"`
}

type Video struct {
	Host          string            `json:"host"`
	Type          string            `json:"type"`
	VideoId       string            `json:"videoId"`
	Thumbnail     string            `json:"thumbnail"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	PublishedTime int64             `json:"publishedTime"`
	Length        int               `json:"length"`
	Views         int64             `json:"views"`
	ChannelId     string            `json:"channelId"`
	ChannelTitle  string            `json:"channelTitle"`
	ChannelAvatar string            `json:"channelAvatar"`
	Sources       map[string]string `json:"sources"`
}
