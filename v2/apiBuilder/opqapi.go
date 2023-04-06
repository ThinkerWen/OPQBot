package apiBuilder

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"log"
	"net/url"
	"strconv"
)

type DoApi interface {
	Do(ctx context.Context) error
	DoAndResponse(ctx context.Context) (*Response, error)
}
type Builder struct {
	qqBot      int64
	url        string
	path       *string
	CgiCmd     *string     `json:"CgiCmd,omitempty"`
	CgiRequest *CgiRequest `json:"CgiRequest,omitempty"`
}
type CgiRequest struct {
	CommandId  *int    `json:"CommandId,omitempty"`
	FilePath   *string `json:"FilePath,omitempty"`
	FileUrl    *string `json:"FileUrl,omitempty"`
	ToUin      *int64  `json:"ToUin,omitempty"`
	ToType     *int    `json:"ToType,omitempty"`
	Content    *string `json:"Content,omitempty"`
	SubMsgType *int    `json:"SubMsgType,omitempty"`
	Images     []*File `json:"Images,omitempty"`
	Uid        *string `json:"Uid,omitempty"`
	AtUinLists []struct {
		QQUin *int64 `json:"QQUin,omitempty"`
	} `json:"AtUinLists,omitempty"`
}

func (b *Builder) BuildStringBody() (string, error) {
	body, err := json.Marshal(b)
	return string(body), err
}

func (b *Builder) Do(ctx context.Context) error {
	r, err := b.DoAndResponse(ctx)
	if err != nil {
		return err
	}
	if !r.Ok() {
		return fmt.Errorf(r.ErrorMsg())
	}
	return nil
}
func (b *Builder) DoAndResponse(ctx context.Context) (*Response, error) {
	body, err := b.BuildStringBody()
	if err != nil {
		return nil, err
	}
	var u string
	if b.path != nil {
		u, _ = url.JoinPath(b.url, *b.path)
	} else {
		u, _ = url.JoinPath(b.url, "/v1/LuaApiCaller")
	}
	resp, err := req.SetContext(ctx).SetQueryParam("funcname", "MagicCgiCmd").SetQueryParam("qq", strconv.FormatInt(b.qqBot, 10)).SetBodyJsonString(body).Post(u)
	if err != nil {
		return nil, err
	}
	r := NewResponse(resp.Bytes())
	log.Println(resp.String())
	if !r.Ok() {
		return nil, fmt.Errorf(r.ErrorMsg())
	}
	return r, nil
}

func NewApi(url string, botQQ int64) IMainFunc {
	return &Builder{qqBot: botQQ, url: url}
}