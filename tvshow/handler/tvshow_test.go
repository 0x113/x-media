package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0x113/x-media/tvshow/common"
	"github.com/0x113/x-media/tvshow/mocks"
	"github.com/0x113/x-media/tvshow/service"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logrus.SetOutput(ioutil.Discard) // disable logging for tests
}

func TestGetTVShow(t *testing.T) {
	// setup
	client := &mocks.MockClient{}
	tvShowRepo := mocks.NewMockTVShowRepository()
	tvShowService := service.NewTVShowService(client, tvShowRepo)
	e := echo.New()

	testCases := []struct {
		name               string
		json               string
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:               "Success",
			json:               `{"name": "BoJack Horseman"}`,
			expectedStatusCode: 200,
			wantErr:            false,
		},
		{
			name:               "Invalid JSON",
			json:               ``,
			expectedStatusCode: 400,
			wantErr:            true,
		},
		{
			name:               "Non-existent show",
			json:               `{"name": "Silicon Valley"}`,
			expectedStatusCode: 404,
			wantErr:            true,
		},
	}

	handler := tvShowHandler{tvShowService}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// prepare request
			req := httptest.NewRequest(http.MethodPost, "/api/v1/tvshows/get", strings.NewReader(tt.json))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := handler.GetTVShow(c)
			// check error
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.expectedStatusCode, rec.Code)
		})
	}
}

func TestGetAllTVShows(t *testing.T) {
	// setup
	client := &mocks.MockClient{}
	tvShowRepo := mocks.NewMockTVShowRepository()
	tvShowService := service.NewTVShowService(client, tvShowRepo)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tvshows/get/all", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := tvShowHandler{tvShowService}

	if assert.NoError(t, handler.GetAllTVShows(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestUpdateAllTVShows(t *testing.T) {
	// setup
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
	tvShowRepo := mocks.NewMockTVShowRepository()
	tvShowService := service.NewTVShowService(client, tvShowRepo)
	common.Config = &common.Configuration{
		TVShowDirectories: []string{"../service/testdata/three_shows/"},
	}
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tvshows/update/all", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := tvShowHandler{tvShowService}

	if assert.NoError(t, handler.UpdateAllTVShows(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
