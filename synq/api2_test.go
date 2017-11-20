package synq

import (
	"strings"
	"testing"

	"github.com/SYNQfm/SYNQ-Golang/test_helper"
	"github.com/stretchr/testify/require"
)

var testAssetId string
var testVideoIdV2 string

const (
	TEST_AUTH = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwczovL3Rlc3QuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfDU3MjE4MjFiM2ExYWFmYmUxNTlkZGE2NSIsImF1ZCI6InRESzZBdUF0QVc0ckFySzhOSTltMXdJRW5WQU9RcjUxIiwiZXhwIjoxNDkzNDM5NTExLCJpYXQiOjE0NjE4MTcxMTF9.29JkFxoHqCRPIH2wVbT-ZNIMBK8xXLwkjbLmyWxpquE"
)

func init() {
	testAssetId = test_helper.ASSET_ID
	testVideoIdV2 = test_helper.V2_VIDEO_ID
	test_helper.SetSampleDir(sampleDir)
}

func setupTestApiV2(key string) ApiV2 {
	api := NewV2(key)
	url := test_helper.SetupServer("v2")
	api.SetUrl(url)
	return api
}

func TestMakeReq2(t *testing.T) {
	assert := require.New(t)
	api := setupTestApiV2("fake")
	body := strings.NewReader("")
	req, err := api.makeRequest("POST", "url", body)
	assert.Nil(err)
	assert.Equal("POST", req.Method)
	req, err = api.makeRequest("GET", "url", body)
	assert.Nil(err)
	assert.Equal("GET", req.Method)
}

func TestCreate2(t *testing.T) {
	assert := require.New(t)
	api := setupTestApiV2("fake")
	_, err := api.Create()
	assert.NotNil(err)
	api.Key = TEST_AUTH
	video, err := api.Create()
	assert.Nil(err)
	assert.Equal(testVideoIdV2, video.Id)
}

func TestGet2(t *testing.T) {
	assert := require.New(t)
	api := setupTestApiV2(TEST_AUTH)
	_, err := api.GetVideo("")
	assert.NotNil(err)
	video, err := api.GetVideo(testVideoIdV2)
	assert.Nil(err)
	assert.Equal(testVideoIdV2, video.Id)
	assert.Len(video.Assets, 1)
	assert.Equal(testAssetId, video.Assets[0].Id)
}
