package synq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/SYNQfm/helpers/common"
)

const (
	DEFAULT_V2_URL = "http://b9n2fsyd6jbfihx82.stoplight-proxy.io"
)

type ApiV2 struct {
	*BaseApi
	User        string
	Password    string
	TokenExpiry time.Time
}

type VideoList struct {
	Videos []VideoV2 `json:"data"`
}

type ErrorRespV2 struct {
	Message string `json:"message"`
}

type LoginResp struct {
	Token       string    `json:"jwt"`
	TokenExpiry time.Time `json:"exp"`
}

func (a ApiV2) Version() string {
	return "v2"
}

func NewV2(token string, timeouts ...time.Duration) ApiV2 {
	base := NewBase(token, timeouts...)
	base.Url = DEFAULT_V2_URL
	return ApiV2{BaseApi: &base}
}

func (a *ApiV2) handleAuth(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+a.GetKey())
}

func (a ApiV2) getBaseUrl() string {
	return a.GetUrl() + "/v2"
}

func (a *ApiV2) CreateAccount(name string, type_ string) string {
	return ""
}

func (a *ApiV2) makeRequest(method string, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return req, err
	}
	if method == "POST" {
		if strings.Contains(url, "/login") {
			req.Header.Add("content-type", "application/x-www-form-urlencoded")
		} else {
			req.Header.Add("content-type", "application/json")
		}
	}
	a.handleAuth(req)
	return req, nil
}

func (a ApiV2) ParseError(status int, bytes []byte) error {
	if status == 404 {
		return errors.New("404 Item not found")
	}
	type Resp struct {
		Message string `json:"message"`
	}
	resp := Resp{}
	err := json.Unmarshal(bytes, &resp)
	if err != nil {
		return common.NewError("could not parse error %d : %s", status, string(bytes))
	}
	msg := resp.Message
	if msg == "" {
		msg = fmt.Sprintf("Failed with status %d", status)
	}
	return errors.New(msg)
}

func (a *ApiV2) handleGet(url string, v interface{}) error {
	body := bytes.NewBufferString("")
	req, err := a.makeRequest("GET", url, body)
	if err != nil {
		return err
	}
	return handleReq(a, req, v)
}

func Login(user, password string, serverUrl ...string) (ApiV2, error) {
	var api ApiV2
	resp, err := login(user, password, serverUrl...)
	if err != nil {
		return api, err
	}
	api = NewV2(resp.Token)
	api.TokenExpiry = resp.TokenExpiry
	api.User = user
	api.Password = password
	return api, nil
}

func login(user, password string, serverUrl ...string) (LoginResp, error) {
	var r LoginResp
	var u string
	if len(serverUrl) > 0 {
		u = serverUrl[0]
	} else {
		u = DEFAULT_V2_URL
	}
	u = u + "/v2/login"
	form := url.Values{}
	form.Add("email", user)
	form.Add("password", password)
	resp, e := http.PostForm(u, form)
	if e != nil {
		return r, e
	}
	if resp.StatusCode != 200 {
		return r, common.NewError("error getting login %d", resp.StatusCode)
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(bytes, &r)
	return r, nil
}

func (a *ApiV2) Create(body ...[]byte) (VideoV2, error) {
	resp := VideoResp{}
	video := VideoV2{}
	url := a.getBaseUrl() + "/videos"
	buf := bytes.NewBufferString("")
	if len(body) > 0 {
		buf.Write(body[0])
	}
	req, err := a.makeRequest("POST", url, buf)
	if err != nil {
		return video, err
	}
	err = handleReq(a, req, &resp)
	if err != nil {
		return video, err
	}
	video = resp.Video
	video.Api = a
	return video, nil
}

func (a *ApiV2) GetVideos(accountId string) (videos []VideoV2, err error) {
	var resp VideoList
	path := "/videos"
	if accountId != "" {
		path = "/accounts/" + accountId + path
	}
	url := a.getBaseUrl() + path
	req, err := a.makeRequest("GET", url, nil)
	if err != nil {
		return videos, err
	}
	err = handleReq(a, req, &resp)
	if err != nil {
		return videos, err
	}
	for _, v := range resp.Videos {
		v.Api = a
		videos = append(videos, v)
	}
	return videos, nil
}

// Helper function to get details for a video, will create video object
func (a *ApiV2) GetVideo(id string) (video VideoV2, err error) {
	var resp VideoResp
	if id == "" || (len(id) != 32 && len(id) != 36) {
		return video, errors.New("video id is blank")
	}
	uuid := common.ConvertToUUIDFormat(id)
	url := a.getBaseUrl() + "/videos/" + uuid
	req, err := a.makeRequest("GET", url, nil)
	if err != nil {
		return video, err
	}
	err = handleReq(a, req, &resp)
	if err != nil {
		return video, err
	}
	video = resp.Video
	video.Api = a
	return video, nil
}

func (a *ApiV2) GetAssetList() ([]Asset, error) {
	list := AssetList{}
	url := a.getBaseUrl() + "/assets"
	err := a.handleGet(url, &list)
	return list.Assets, err
}
