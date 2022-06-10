package types

type VideoItemList struct {
	ItemList   []VideoItem `json:"itemList"`
	HasMore    bool        `json:"hasMore"`
	StatusCode int         `json:"statusCode"`
	Cursor     string      `json:"cursor"`
}

type VideoItem struct {
	Id   string `json:"id"`
	Desc string `json:"desc"`
}

type VideoDetail struct {
	StatusCode   int            `json:"status_code"`
	AwemeDetails []AwemeDetails `json:"aweme_details"`
	StatusMsg    string         `json:"status_msg"`
}

type AwemeDetails struct {
	Desc   string      `json:"desc"`
	Video  AwemeVideo  `json:"video"`
	Author AwemeAuthor `json:"author"`
}

type AwemeVideo struct {
	PlayAddr AwemePlayAddr `json:"play_addr"`
}

type AwemeAuthor struct {
	Nickname string `json:"nickname"`
}

type AwemePlayAddr struct {
	UrlList []string `json:"url_list"`
}

type HomePageContent struct {
	UserPage struct {
		SecUid string `json:"secUid"`
	} `json:"UserPage"`
	ItemList struct {
		UserPost struct {
			Cursor string `json:"cursor"`
		} `json:"user-post"`
	} `json:"ItemList"`
}
