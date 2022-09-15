package odysee

import (
	"encoding/json"
	"errors"
	"github.com/alist-org/alist/v3/drivers/base"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
	"strconv"
	"time"
)

// do others that not defined in Driver interface

func (d *Odysee) request(pathname string, method string, param map[string]any, callback base.ReqCallback, resp interface{}) ([]byte, error) {
	u := "https://api.na-backend.odysee.com/api/v1/proxy"

	req := base.RestyClient.SetTimeout(time.Second * 1000).R()
	req.SetQueryParam("m", pathname)
	req.SetHeader("x-lbry-auth-token", d.AuthToken)
	req.SetBody(Request{
		Method: pathname,
		Params: param,
	})

	if callback != nil {
		callback(req)
	}

	req.SetResult(resp)
	var e Resp[map[string]any]
	req.SetError(&e)
	res, err := req.Execute(method, u)
	if err != nil {
		return nil, err
	}
	return res.Body(), nil
}

func (d *Odysee) listChannel() ([]ChannelItem, error) {
	res := make([]ChannelItem, 0)
	var resp Resp[ChannelItems]
	_, err := d.request("channel_list", http.MethodPost, map[string]any{
		"page":      1,
		"page_size": 99999,
		"resolve":   true,
	}, nil, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != (Error{}) {
		return nil, errors.New(strconv.Itoa(resp.Error.Code) + resp.Error.Message)
	}
	res = append(res, resp.Result.Items...)

	return res, nil
}

func (d *Odysee) listChannelFile(id string, page int) ([]ChannelItem, error) {
	res := make([]ChannelItem, 0)
	var resp Resp[ChannelItems]
	_, err := d.request("claim_list", http.MethodPost, map[string]any{
		"page_size":  20,
		"page":       page,
		"channel_id": id,
	}, nil, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != (Error{}) {
		return nil, errors.New(strconv.Itoa(resp.Error.Code) + resp.Error.Message)
	}
	if resp.Result.TotalPages > page {
		resNext, _ := d.listChannelFile(id, page+1)
		res = append(res, resNext...)
	}

	res = append(res, resp.Result.Items...)

	result := make([]ChannelItem, 0)

	for i := range res {
		if res[i].ValueType == "stream" {
			result = append(result, res[i])
		}
	}
	return result, nil
}

func (d *Odysee) getFileDetail(uri string) (Detail, error) {
	var resp Resp[Detail]
	_, err := d.request("get", http.MethodPost, map[string]any{
		"uri":       uri,
		"save_file": false,
	}, nil, &resp)
	if err != nil {
		return Detail{}, err
	}
	if resp.Error != (Error{}) {
		return Detail{}, errors.New(strconv.Itoa(resp.Error.Code) + resp.Error.Message)
	}
	return resp.Result, nil
}

func (d *Odysee) DeleteStreamByClaimId(id string) error {
	var resp Resp[Detail]

	_, err := d.request("stream_abandon", http.MethodPost, map[string]any{
		"claim_id": id,
	}, nil, &resp)
	if err != nil {
		return errors.New(resp.Error.Message)
	}
	return nil
}

func (d *Odysee) upCommit(tempFile *os.File, stream model.FileStreamer) error {
	var resp Resp[Detail]
	_, err := d.request("publish", http.MethodPost, nil, func(req *resty.Request) {
		data := &Request{
			Method: "publish",
			Params: map[string]any{
				"name":          stream.GetName(),
				"title":         stream.GetName(),
				"description":   "",
				"bid":           "0.01000000",
				"thumbnail_url": "https://thumbs.odycdn.com/70a526314962c806435b8aab45f5e06e.webp",
				"release_time":  1663233797,
				"blocking":      false,
				"preview":       false,
				"license":       "None",
				"channel_id":    "f976dc1d600cd39d379eb71159472a4d054c628a",
				"file_path":     "__POST_FILE__",
				"optimize_file": false,
			},
		}
		dataStr, _ := json.Marshal(data)
		req.SetFormData(map[string]string{
			"json_payload": string(dataStr[:]),
		})
		req.SetFileReader("file", stream.GetName(), tempFile)
	}, &resp)
	if err != nil {
		return errors.New(resp.Error.Message)
	}
	return nil
}
