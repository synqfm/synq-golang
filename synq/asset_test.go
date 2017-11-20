package synq

import (
	"log"
	"net/http"
	"testing"

	"github.com/SYNQfm/SYNQ-Golang/test_helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ASSET_TYPE     = "mp4"
	ASSET_CREATED  = "created"
	ASSET_UPLOADED = "uploaded"
	ASSET_LOCATION = "https://s3.amazonaws.com/synq-jessica/uploads/01/82/01823629bcf24c34b714ae21e1a4647f/01823629bcf24c34b714ae21e1a4647f.mp4"
)

func setupTestVideoV2() VideoV2 {
	api := setupTestApiV2(TEST_AUTH)
	video, _ := api.Create()
	url := test_helper.SetupServer("v2")
	video.Api.SetUrl(url)
	return video
}

func handleAsset(w http.ResponseWriter, r *http.Request) {
}

func TestGetAssetList(t *testing.T) {
	log.Println("Testing GetAssetList")
	assert := require.New(t)
	api := setupTestApiV2(TEST_AUTH)
	assets, err := api.GetAssetList()
	assert.Nil(err)
	assert.NotEmpty(assets)
	assert.Equal(testVideoIdV2, assets[0].VideoId)
	assert.Equal(testAssetId, assets[0].Id)
	assert.Equal(ASSET_TYPE, assets[0].Type)
	assert.Equal(ASSET_LOCATION, assets[0].Location)
	assert.Equal(ASSET_CREATED, assets[0].State)
}

func TestGetVideoAssetList(t *testing.T) {
	log.Println("Testing GetVideoAssetList")
	assert := require.New(t)
	video := setupTestVideoV2()
	err := video.GetVideoAssetList()
	assert.Nil(err)
	assert.NotEmpty(video.Assets)
	assert.Equal(testVideoIdV2, video.Assets[0].VideoId)
	assert.Equal(testAssetId, video.Assets[0].Id)
	assert.Equal(ASSET_TYPE, video.Assets[0].Type)
	assert.Equal(ASSET_LOCATION, video.Assets[0].Location)
	assert.Equal(ASSET_CREATED, video.Assets[0].State)
}

func TestGetAsset(t *testing.T) {
	log.Println("Testing GetAsset")
	assert := assert.New(t)
	video := setupTestVideoV2()
	asset, err := video.GetAsset(testAssetId)
	assert.Equal(testVideoIdV2, asset.VideoId)
	assert.Equal(ASSET_TYPE, asset.Type)
	assert.Equal(ASSET_UPLOADED, asset.State)
	assert.Equal(ASSET_LOCATION, asset.Location)
	assert.Nil(err)
}

func TestUpdate(t *testing.T) {
	log.Println("Testing Update")
	assert := assert.New(t)
	video := setupTestVideoV2()
	asset, _ := video.GetAsset(testAssetId)
	assert.NotEmpty(asset)
	asset.State = ASSET_UPLOADED
	err := asset.Update()
	assert.Nil(err)
	asset, _ = video.GetAsset(testAssetId)
	assert.Equal(ASSET_UPLOADED, asset.State)
}

func TestDelete(t *testing.T) {
	log.Println("Testing Delete")
	assert := assert.New(t)
	video := setupTestVideoV2()
	asset, _ := video.GetAsset(testAssetId)
	assert.NotEmpty(asset)
	err := asset.Delete()
	assert.Nil(err)
}

func TestHandleAssetReq(t *testing.T) {
	log.Println("Testing TestHandleAssetReq")
	assert := assert.New(t)
	video := setupTestVideoV2()
	asset, _ := video.GetAsset(testAssetId)
	ogAsset := asset
	url := video.Api.getBaseUrl() + "/assets/" + testAssetId
	err := asset.handleAssetReq("GET", url, nil)
	assert.Nil(err)
	assert.Equal(ogAsset, asset)
}
