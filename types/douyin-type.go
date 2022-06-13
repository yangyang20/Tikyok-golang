package types

type VideoList struct {
	AwemeList []AwemeDetails `json:"aweme_list"`
	MaxCursor int            `json:"max_cursor"`
	MinCursor int            `json:"min_cursor"`
	HasMore   bool           `json:"has_more"`
}
