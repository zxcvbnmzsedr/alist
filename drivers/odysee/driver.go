package odysee

import (
	"context"
	"errors"
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/pkg/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Odysee struct {
	model.Storage
	Addition
}

func (d *Odysee) Config() driver.Config {
	return config
}

func (d *Odysee) GetAddition() driver.Additional {
	return d.Addition
}

func (d *Odysee) Init(ctx context.Context, storage model.Storage) error {
	d.Storage = storage
	err := utils.Json.UnmarshalFromString(d.Storage.Addition, &d.Addition)
	if err != nil {
		return err
	}
	var resp Resp[ChannelItems]
	_, err = d.request("channel_list", http.MethodPost, map[string]any{
		"page":      1,
		"page_size": 99999,
		"resolve":   true,
	}, nil, &resp)
	if resp.Error != (Error{}) {
		return errors.New(strconv.Itoa(resp.Error.Code) + resp.Error.Message)
	}
	return err
}

func (d *Odysee) Drop(ctx context.Context) error {
	return nil
}

func (d *Odysee) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	id := dir.GetID()
	path := dir.GetPath()
	if id == "" && path == "" {
		// 一级查询频道
		files, err := d.listChannel()
		if err != nil {
			return nil, err
		}
		return utils.SliceConvert(files, func(src ChannelItem) (model.Obj, error) {
			return fileToObj(src, 0), nil
		})
	} else if id != "" && len(strings.Split(path, "/")) == 2 {
		// 二级查询频道文件列表
		files, err := d.listChannelFile(id, 1)
		if err != nil {
			return nil, err
		}
		return utils.SliceConvert(files, func(src ChannelItem) (model.Obj, error) {
			return fileToObj(src, 1), nil
		})
	}
	return nil, nil

}

func (d *Odysee) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	detail, err := d.getFileDetail(file.GetID())
	if err != nil {
		return nil, err
	}
	return &model.Link{
		Header: http.Header{
			"Referer": []string{"https://odysee.com/"},
		},
		URL: detail.StreamingUrl,
	}, nil
}

func (d *Odysee) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error {
	return nil
}

func (d *Odysee) Move(ctx context.Context, srcObj, dstDir model.Obj) error {
	return nil
}

func (d *Odysee) Rename(ctx context.Context, srcObj model.Obj, newName string) error {
	return nil
}

func (d *Odysee) Copy(ctx context.Context, srcObj, dstDir model.Obj) error {
	return nil
}

func (d *Odysee) Remove(ctx context.Context, obj model.Obj) error {
	id := obj.GetID()
	if strings.Contains(id, "#") {
		claimId := strings.Split(id, "#")[1]
		return d.DeleteStreamByClaimId(claimId)
	}
	return nil
}

func (d *Odysee) Put(ctx context.Context, dstDir model.Obj, stream model.FileStreamer, up driver.UpdateProgress) error {
	tempFile, err := utils.CreateTempFile(stream.GetReadCloser())
	if err != nil {
		return err
	}
	defer func() {
		_ = tempFile.Close()
		_ = os.Remove(tempFile.Name())
	}()
	return d.upCommit(tempFile, stream)
}

var _ driver.Driver = (*Odysee)(nil)
