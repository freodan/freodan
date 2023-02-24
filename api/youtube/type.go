package youtube

type initialData struct {
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Title    string `json:"title"`
					Selected bool   `json:"selected"`
					Content  struct {
						RichGridRenderer struct {
							Contents []contents `json:"contents"`
						} `json:"richGridRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
	Header struct {
		C4TabbedHeaderRenderer struct {
			SubscriberCountText struct {
				SimpleText string `json:"simpleText"`
			} `json:"subscriberCountText"`
		} `json:"c4TabbedHeaderRenderer"`
	} `json:"header"`
	Metadata struct {
		ChannelMetadataRenderer struct {
			Avatar struct {
				Thumbnails []struct {
					Url    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"avatar"`
			Description string `json:"description"`
			ExternalId  string `json:"externalId"`
			Title       string `json:"title"`
		} `json:"channelMetadataRenderer"`
	} `json:"metadata"`
}

type contents struct {
	RichItemRenderer struct {
		Content struct {
			VideoRenderer    videoRenderer    `json:"videoRenderer"`
			ReelItemRenderer reelItemRenderer `json:"reelItemRenderer"`
		} `json:"content"`
	} `json:"richItemRenderer"`
}

type reelItemRenderer struct {
	VideoId  string `json:"videoId"`
	Headline struct {
		SimpleText string `json:"simpleText"`
	} `json:"headline"`
	Thumbnail struct {
		Thumbnails []struct {
			Url    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"thumbnails"`
	} `json:"thumbnail"`
	ViewCountText struct {
		SimpleText string `json:"simpleText"`
	} `json:"viewCountText"`
	NavigationEndpoint struct {
		ReelWatchEndpoint struct {
			Overlay struct {
				ReelPlayerOverlayRenderer struct {
					ReelPlayerHeaderSupportedRenderers struct {
						ReelPlayerHeaderRenderer struct {
							TimestampText struct {
								SimpleText string `json:"simpleText"`
							} `json:"timestampText"`
						} `json:"reelPlayerHeaderRenderer"`
					} `json:"reelPlayerHeaderSupportedRenderers"`
				} `json:"reelPlayerOverlayRenderer"`
			} `json:"overlay"`
		} `json:"reelWatchEndpoint"`
	} `json:"navigationEndpoint"`
}

type videoRenderer struct {
	VideoId   string `json:"videoId"`
	Thumbnail struct {
		Thumbnails []struct {
			Url    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"thumbnails"`
	} `json:"thumbnail"`
	ThumbnailOverlays []struct {
		ThumbnailOverlayTimeStatusRenderer struct {
			Style string `json:"style"`
			Text  struct {
				SimpleText string `json:"simpleText"`
			} `json:"text"`
		} `json:"thumbnailOverlayTimeStatusRenderer"`
	} `json:"thumbnailOverlays"`
	Title struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"title"`
	DescriptionSnippet struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"descriptionSnippet"`
	PublishedTimeText struct {
		SimpleText string `json:"simpleText"`
	} `json:"publishedTimeText"`
	LengthText struct {
		SimpleText string `json:"simpleText"`
	} `json:"lengthText"`
	ViewCountText struct {
		SimpleText string `json:"simpleText"`
		Runs       []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"viewCountText"`
	UpcomingEventData struct {
		StartTime         string `json:"startTime"`
		UpcomingEventText struct {
			Runs []struct {
				Text string `json:"text"`
			} `json:"runs"`
		} `json:"upcomingEventText"`
	} `json:"upcomingEventData"`
}

type searchInitialData struct {
	Contents struct {
		TwoColumnSearchResultsRenderer struct {
			PrimaryContents struct {
				SectionListRenderer struct {
					Contents []struct {
						ItemSectionRenderer struct {
							Contents []struct {
								ChannelRenderer channelRenderer `json:"channelRenderer"`
							} `json:"contents"`
						} `json:"itemSectionRenderer"`
					} `json:"contents"`
				} `json:"sectionListRenderer"`
			} `json:"primaryContents"`
		} `json:"twoColumnSearchResultsRenderer"`
	} `json:"contents"`
}

type channelRenderer struct {
	ChannelId          string `json:"channelId"`
	DescriptionSnippet struct {
		Runs []struct {
			Text string `json:"text"`
			Bold bool   `json:"bold"`
		} `json:"runs"`
	} `json:"descriptionSnippet"`
	Title struct {
		SimpleText string `json:"simpleText"`
	} `json:"title"`
	Thumbnail struct {
		Thumbnails []struct {
			Url    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"thumbnails"`
	} `json:"thumbnail"`
	VideoCountText struct {
		SimpleText string `json:"simpleText"`
	} `json:"videoCountText"`
}
