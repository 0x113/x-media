package tvmaze_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/0x113/x-media/tvshow/external/tvmaze"
	"github.com/0x113/x-media/tvshow/mocks"

	"github.com/stretchr/testify/assert"
)

func TestSuccessGetTVmazeTVShowInfo(t *testing.T) {
	client := &mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			json := `[{
		"score": 25.39485,
		"show": {
			"id": 526,
			"url": "http://www.tvmaze.com/shows/526/the-office",
			"name": "The Office",
			"type": "Scripted",
			"language": "English",
			"genres": [
				"Comedy"
			],
			"status": "Ended",
			"runtime": 30,
			"premiered": "2005-03-24",
			"officialSite": "http://www.nbc.com/the-office",
			"schedule": {
				"time": "21:00",
				"days": [
					"Thursday"
				]
			},
			"rating": {
				"average": 8.5
			},
			"weight": 97,
			"network": {
				"id": 1,
				"name": "NBC",
				"country": {
					"name": "United States",
					"code": "US",
					"timezone": "America/New_York"
				}
			},
			"webChannel": null,
			"externals": {
				"tvrage": 6061,
				"thetvdb": 73244,
				"imdb": "tt0386676"
			},
			"image": {
				"medium": "http://static.tvmaze.com/uploads/images/medium_portrait/85/213184.jpg",
				"original": "http://static.tvmaze.com/uploads/images/original_untouched/85/213184.jpg"
			},
			"summary": "One of the best tv shows, no doubt",
			"updated": 1583654209,
			"_links": { "self": {
					"href": "http://api.tvmaze.com/shows/526"
				},
				"previousepisode": {
					"href": "http://api.tvmaze.com/episodes/711203"
				}
			}
		}
	},
	{
	}]
			`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}

	tvMazeInfo, err := tvmaze.GetTVmazeTVShowInfo(client, "The Office")
	assert.Nil(t, err)
	assert.Equal(t, "The Office", tvMazeInfo.Show.Name)
	assert.NotNil(t, tvMazeInfo)
}

func TestFailGetTVmazeTVShowInfo(t *testing.T) {
	client := &mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(nil),
			}, nil
		},
	}

	tvMazeInfo, err := tvmaze.GetTVmazeTVShowInfo(client, "Wrong status code")
	assert.NotNil(t, err)
	assert.Nil(t, tvMazeInfo)
}

func TestFailDoFuncErrorGetTVmazeTVShowInfo(t *testing.T) {
	client := &mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{}, errors.New("This method shouldn't be called")
		},
	}

	tvMazeInfo, err := tvmaze.GetTVmazeTVShowInfo(client, "DoFunc error")
	assert.NotNil(t, err)
	assert.Nil(t, tvMazeInfo)
}
