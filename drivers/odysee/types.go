package odysee

import (
	"github.com/alist-org/alist/v3/internal/model"
	"strconv"
	"time"
)

type Request struct {
	JsonRpc string         `json:"jsonrpc" default:"2.0"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
}

type Resp[T interface{}] struct {
	JsonRpc string `json:"jsonrpc"`
	Error   Error  `json:"error"`
	Result  T      `json:"result"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ChannelItems struct {
	Items      []ChannelItem `json:"items"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalItems int           `json:"total_items"`
	TotalPages int           `json:"total_pages"`
}

type ChannelItem struct {
	PermanentUrl string `json:"permanent_url"`
	Name         string `json:"name"`
	ClaimId      string `json:"claim_id"`
	CanonicalUrl string `json:"canonical_url"`
	ShortUrl     string `json:"short_url"`
	Value        Value  `json:"value"`
	ValueType    string `json:"value_type"`
	Timestamp    int64  `json:"timestamp"`
}
type Thumbnail struct {
	Url string `json:"url"`
}

func (c *ChannelItem) GetSize() int64 {
	size, _ := strconv.ParseInt(c.Value.Source.Size, 10, 64)
	return size
}
func (c *ChannelItem) GetName() string {
	if c.ValueType == "stream" {
		return c.Value.Source.Name
	}
	if c.Value.Title != "" {
		return c.Value.Title
	}
	return c.Name
}
func (c *ChannelItem) ModTime() time.Time {
	return time.Unix(c.Timestamp, 0)
}
func (c *ChannelItem) IsDir() bool {
	return c.ValueType != "stream"
}
func (c *ChannelItem) GetID() string {
	if c.ValueType == "stream" {
		return c.PermanentUrl
	}
	if c.ValueType == "collection" {
		return "collection_" + c.ClaimId
	}
	return "channel_" + c.ClaimId
}
func (c *ChannelItem) GetPath() string {
	return c.PermanentUrl
}
func (c *ChannelItem) Thumb() string {
	return c.Value.Thumbnail.Url
}

type Value struct {
	Thumbnail Thumbnail `json:"thumbnail"`
	Source    Source    `json:"source"`
	Title     string    `json:"title"`
}
type Source struct {
	Name string `json:"name"`
	Size string `json:"size"`
}
type Detail struct {
	StreamingUrl string `json:"streaming_url"`
}

func fileToObj(f ChannelItem) *model.ObjThumb {
	return &model.ObjThumb{
		Object: model.Object{
			ID:       f.GetID(),
			Name:     f.GetName(),
			Size:     f.GetSize(),
			IsFolder: f.IsDir(),
			Modified: f.ModTime(),
		},
		Thumbnail: model.Thumbnail{Thumbnail: f.Thumb()},
	}

}
