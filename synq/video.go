package synq

import (
	"errors"
	"net/url"
	"time"
)

/*
{
  "input": {
    "url": "https://multicdn.synq.fm/projects/0a/bf/0abfe1b849154082993f2fce77a16fd9/uploads/videos/45/d4/45d4063d00454c9fb21e5186a09c3115.mp4",
    "width": 720,
    "height": 1280,
    "duration": 17.48,
    "file_size": 16706384,
    "framerate": 29.97,
    "uploaded_at": "2017-02-15T03:05:17.978Z"
  },
  "state": "uploaded",
  "player": {
    "views": 0,
    "embed_url": "https://player.synq.fm/embed/45d4063d00454c9fb21e5186a09c3115",
    "thumbnail_url": "https://multicdn.synq.fm/projects/0a/bf/0abfe1b849154082993f2fce77a16fd9/derivatives/thumbnails/45/d4/45d4063d00454c9fb21e5186a09c3115/0000360.jpg"
  },
  "outputs": {
    "hls": {
      "url": "https://multicdn.synq.fm/projects/0a/bf/0abfe1b849154082993f2fce77a16fd9/derivatives/videos/45/d4/45d4063d00454c9fb21e5186a09c3115/hls/45d4063d00454c9fb21e5186a09c3115_hls.m3u8",
      "state": "complete"
    },
    "mp4_360": {
      "url": "https://multicdn.synq.fm/projects/0a/bf/0abfe1b849154082993f2fce77a16fd9/derivatives/videos/45/d4/45d4063d00454c9fb21e5186a09c3115/mp4_360/45d4063d00454c9fb21e5186a09c3115_mp4_360.mp4",
      "state": "complete"
    },
    "mp4_720": {
      "url": "https://multicdn.synq.fm/projects/0a/bf/0abfe1b849154082993f2fce77a16fd9/derivatives/videos/45/d4/45d4063d00454c9fb21e5186a09c3115/mp4_720/45d4063d00454c9fb21e5186a09c3115_mp4_720.mp4",
      "state": "complete"
    },
    "mp4_1080": {
      "url": "https://multicdn.synq.fm/projects/0a/bf/0abfe1b849154082993f2fce77a16fd9/derivatives/videos/45/d4/45d4063d00454c9fb21e5186a09c3115/mp4_1080/45d4063d00454c9fb21e5186a09c3115_mp4_1080.mp4",
      "state": "complete"
    },
    "webm_720": {
      "url": "https://multicdn.synq.fm/projects/0a/bf/0abfe1b849154082993f2fce77a16fd9/derivatives/videos/45/d4/45d4063d00454c9fb21e5186a09c3115/webm_720/45d4063d00454c9fb21e5186a09c3115_webm_720.webm",
      "state": "complete"
    }
  },
  "userdata": {},
  "video_id": "45d4063d00454c9fb21e5186a09c3115",
  "created_at": "2017-02-15T03:01:16.767Z",
  "updated_at": "2017-02-15T03:06:31.794Z"
}
*/
type Video struct {
	Id        string                 `json:"video_id"`
	Outputs   map[string]interface{} `json:"outputs"`
	Player    map[string]interface{} `json:"player"`
	Input     map[string]interface{} `json:"intput"`
	State     string                 `json:"state"`
	Userdata  map[string]interface{} `json:"userdata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Api       Api
}

func New(key string) Video {
	api := Api{Key: key}
	api.Url = DEFAULT_URL
	api.Timeout = DEFAULT_TIMEOUT_MS
	return Video{Api: api}
}

func (v *Video) Create() error {
	if v.Id != "" {
		return errors.New("This video already has an Id (" + v.Id + "), can not create")
	}
	form := url.Values{}
	return v.Api.handlePost("create", form, v)
}

func (v *Video) Details() error {
	if v.Id == "" {
		return errors.New("There is no id associated with this video object")
	}
	form := url.Values{}
	form.Add("video_id", v.Id)
	return v.Api.handlePost("details", form, v)
}
