package odysee

import (
	"github.com/alist-org/alist/v3/internal/model"
	"strconv"
	"strings"
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
}
type Value struct {
	Source Source `json:"source"`
}
type Source struct {
	Name string `json:"name"`
	Size string `json:"size"`
}
type Detail struct {
	StreamingUrl string `json:"streaming_url"`
}

func fileToObj(f ChannelItem, level int) *model.ObjThumb {
	if level == 0 {
		return &model.ObjThumb{
			Object: model.Object{
				ID:       f.ClaimId,
				Name:     strings.ReplaceAll(f.Name, "@", ""),
				Size:     0,
				IsFolder: true,
			},
		}
	} else {
		size, _ := strconv.ParseInt(f.Value.Source.Size, 10, 64)
		return &model.ObjThumb{
			Object: model.Object{
				ID:       f.PermanentUrl,
				Name:     f.Value.Source.Name,
				Size:     size,
				IsFolder: false,
			},
		}
	}

}
