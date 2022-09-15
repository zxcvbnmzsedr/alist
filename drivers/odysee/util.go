package odysee

import (
	"errors"
	"github.com/alist-org/alist/v3/drivers/base"
	"net/http"
	"strconv"
)

// do others that not defined in Driver interface

func (d *Odysee) request(pathname string, method string, param map[string]any, callback base.ReqCallback, resp interface{}) ([]byte, error) {
	u := "https://api.na-backend.odysee.com/api/v1/proxy"

	req := base.RestyClient.R()
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

func (d *Odysee) listChannelFile(id string) ([]ChannelItem, error) {
	res := make([]ChannelItem, 0)
	var resp Resp[ChannelItems]
	_, err := d.request("claim_search", http.MethodPost, map[string]any{
		"page_size": 50,
		"page":      1,
		"channel_ids": [1]string{
			id,
		},
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
