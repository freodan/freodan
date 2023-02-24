package lbry

type lbryClaim struct {
	CanonicalUrl   string `json:"canonical_url"`
	Name           string `json:"name"`
	Title          string `json:"title"`
	SigningChannel struct {
		Value struct {
			Title     string `json:"title"`
			Thumbnail struct {
				Url string `json:"url"`
			} `json:"thumbnail"`
		} `json:"value"`
	} `json:"signing_channel"`
	Value struct {
		Description string   `json:"description"`
		ReleaseTime string   `json:"release_time"`
		Title       string   `json:"title"`
		Tags        []string `json:"tags"`
		Thumbnail   struct {
			Url string `json:"url"`
		} `json:"thumbnail"`
		Video struct {
			Duration int `json:"duration"`
			Height   int `json:"height"`
			Width    int `json:"width"`
		} `json:"video"`
	} `json:"value"`
}

type lbryClaimSearchResponse struct {
	Result struct {
		Items []lbryClaim `json:"items"`
	} `json:"result"`
}

type lbryRequest struct {
	Method string            `json:"method"`
	Params lbryRequestParams `json:"params"`
}

type lbryRequestParams struct {
	Channel     string   `json:"channel,omitempty"`
	ChannelIds  []string `json:"channel_ids,omitempty"`
	ClaimType   []string `json:"claim_type,omitempty"`
	HasSource   bool     `json:"has_source,omitempty"`
	OrderBy     string   `json:"order_by,omitempty"`
	PageSize    int      `json:"page_size,omitempty"`
	ReleaseTime string   `json:"release_time,omitempty"`
	Urls        []string `json:"urls,omitempty"`
}

type lbryResolveResponse struct {
	Result map[string]lbryClaim `json:"result"`
}

type lbryYtResolveResponse struct {
	Data struct {
		Videos   map[string]string `json:"videos"`
		Channels map[string]string `json:"channels"`
	} `json:"data"`
}

type odyseeSearchResult struct {
	ClaimId string `json:"claimId"`
	Name    string `json:"name"`
}

type odyseeSubCountResult struct {
	Data []int64 `json:"data"`
}

type YtResolveResult struct {
	Videos   map[string]string `json:"videos"`
	Channels map[string]string `json:"channels"`
}
