package types

type Config struct {
	DownloadPath string   `toml:"download_path"`
	Url          string   `toml:"url"`
	Grpc         GrpcConf `toml:"grpc"`
	IsUHD        bool     `toml:"isUHD"`
}
type GrpcConf struct {
	Addr string
	Port string
}

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
	Desc    string      `json:"desc"`
	Video   AwemeVideo  `json:"video"`
	Author  AwemeAuthor `json:"author"`
	AwemeId string      `json:"aweme_id"`
}

type AwemeVideo struct {
	PlayAddr AwemePlayAddr `json:"play_addr"`
}

type AwemeAuthor struct {
	Nickname string `json:"nickname"`
}

type AwemePlayAddr struct {
	UrlList []string `json:"url_list"`
	Uri     string   `json:"uri"`
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

type Video struct {
	DownUrl string
	Name    string
	//视频信息
	Uri string
}
