package synq

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/SYNQfm/SYNQ-Golang/test_server"
	"github.com/stretchr/testify/require"
)

func TestMarshalVideo(t *testing.T) {
	assert := require.New(t)
	video := setupTestVideoV2()
	video.Metadata = []byte(`{"test":true}`)
	bytes, err := json.Marshal(video)
	assert.Nil(err)
	v := VideoV2{}
	json.Unmarshal(bytes, &v)
	assert.Equal(video.Metadata, v.Metadata)
}

func TestVideoUpdate(t *testing.T) {
	assert := require.New(t)
	video := setupTestVideoV2()
	video.Metadata = json.RawMessage(`{"meta":"new"}`)
	video.Userdata = json.RawMessage(`{"user":"new"}`)
	video.CompletenessScore = 95.4
	// this is just fake, the updated value will come from a hard coded json
	err := video.Update()
	assert.Nil(err)
	assert.Equal(`{"type":"show"}`, string(video.Metadata))
	assert.Contains(string(video.Userdata), "test2")
	reqs, vals := test_server.GetReqs()
	assert.Len(reqs, 1)
	assert.Len(vals, 1)
	assert.Equal(`{"metadata":{"meta":"new"},"user_data":{"user":"new"},"completeness_score":95.4}`, vals[0].Get("body"))
}

func TestCreateAsset(t *testing.T) {
	log.Println("Testing CreateAsset")
	assert := require.New(t)
	video := setupTestVideoV2()
	asset, err := video.CreateAsset(ASSET_CREATED, ASSET_TYPE, ASSET_LOCATION)
	assert.Nil(err)
	assert.NotNil(asset.Id)
	assert.Equal(testAssetId, asset.Id)
}

func TestCreateOrUpdateAsset(t *testing.T) {
	log.Println("Testing CreateAsset")
	assert := require.New(t)
	video := setupTestVideoV2()
	asset := Asset{State: ASSET_CREATED, Type: ASSET_TYPE, Location: ASSET_LOCATION}
	err := video.CreateOrUpdateAsset(&asset)
	assert.Nil(err)
	assert.Equal(testAssetId, asset.Id)
	asset.State = ASSET_UPLOADED
	err = video.CreateOrUpdateAsset(&asset)
	reqs, vals := test_server.GetReqs()
	assert.Nil(err)
	assert.Len(reqs, 2)
	assert.Len(vals, 2)
	req := reqs[1]
	val := vals[1]
	assert.Equal("PUT", req.Method)
	a := Asset{}
	body := val.Get("body")
	json.Unmarshal([]byte(body), &a)
	assert.Equal(asset.State, a.State)
	assert.Equal(asset.Location, a.Location)
}

func TestCreateAssetForUpload(t *testing.T) {
	log.Println("Testing GettingAssetForUpload")
	assert := require.New(t)
	video := setupTestVideoV2()
	asset, err := video.CreateAssetForUpload("video/mp4")
	assert.Nil(err)
	assert.Len(video.Assets, 1)
	assert.Equal(asset.Id, video.Assets[0].Id)
	assert.Equal("uploads/9e/9d/9e9dc8c8f70541db88dab3034894deb9/01823629bcf24c34b714ae21e1a4647f.mp4", asset.UploadParameters.Key)
	assert.Equal("https://synq-bruce.s3.amazonaws.com", asset.UploadParameters.Action)
}

func TestAddAccount(t *testing.T) {
	assert := require.New(t)
	video := setupTestVideoV2()
	err := video.AddAccount(test_server.ACCOUNT_ID)
	assert.Nil(err)
	reqs, vals := test_server.GetReqs()
	assert.Len(reqs, 1)
	val := vals[0]
	body := val.Get("body")
	obj := struct {
		Accounts []Account `json:"video_accounts"`
	}{}
	json.Unmarshal([]byte(body), &obj)
	assert.Len(obj.Accounts, 1)
	assert.Equal(test_server.ACCOUNT_ID, obj.Accounts[0].Id)
}
