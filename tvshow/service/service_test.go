package service_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/0x113/x-media/tvshow/common"
	"github.com/0x113/x-media/tvshow/mocks"
	"github.com/0x113/x-media/tvshow/models"
	"github.com/0x113/x-media/tvshow/service"
	"github.com/0x113/x-media/tvshow/utils"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// TVShowServiceTestSuite define the suite
type TVShowServiceTestSuite struct {
	suite.Suite
	tvShowRepo    *mocks.MockTVShowRepository
	tvShowService service.TVShowService
	client        utils.HttpClient
}

// SetupTest initiates mocked database and new tv show service
func (suite *TVShowServiceTestSuite) SetupTest() {
	suite.tvShowRepo = mocks.NewMockTVShowRepository()
	logrus.SetOutput(ioutil.Discard) // disable logrus
}

// TestTVShowServiceTestSuite runs test suite
func TestTVShowServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TVShowServiceTestSuite))
}

func (suite *TVShowServiceTestSuite) TestSave() {
	suite.client = &mocks.MockClient{}
	suite.tvShowService = service.NewTVShowService(suite.client, suite.tvShowRepo)
	testCases := []struct {
		name    string
		tvShow  *models.TVShow
		wantErr bool
	}{
		{
			name: "Success",
			tvShow: &models.TVShow{
				Name:      "The Office",
				Language:  "English",
				Genres:    []string{"Comedy"},
				Runtime:   30,
				Premiered: "2005-03-24",
				Rating:    8.5,
				PosterURL: "https://static.tvmaze.com/uploads/images/original_untouched/85/213184.jpg",
				Summary:   "One of the best tv shows ever",
				DirPath:   "testdata/three_shows/The_Office",
			},
			wantErr: false,
		},
		{
			name: "Existing show",
			tvShow: &models.TVShow{
				Name:      "BoJack Horseman",
				Language:  "English",
				Genres:    []string{"Comedy", "Drama"},
				Runtime:   25,
				Premiered: "2014-08-22",
				Rating:    8.1,
				PosterURL: "https://static.tvmaze.com/uploads/images/original_untouched/236/590384.jpg",
				Summary:   "Meet the most beloved sitcom horse of the '90s, 20 years later.",
				DirPath:   "testdata/three_shows/BoJack Horseman",
			},
			wantErr: true,
		},
		{
			name:    "Invalid struct",
			tvShow:  &models.TVShow{},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			err := suite.tvShowService.Save(tt.tvShow)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *TVShowServiceTestSuite) TestUpdateAllTVShows() {
	// set up config
	suite.client = &mocks.MockClient{
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
	suite.tvShowService = service.NewTVShowService(suite.client, suite.tvShowRepo)
	common.Config = &common.Configuration{
		TVShowDirectories: []string{"testdata/three_shows/"},
	}

	_, errMap := suite.tvShowService.UpdateAllTVShows()
  expectedErrMap := map[string]string{}
	suite.Equal(expectedErrMap, errMap)
}

func (suite *TVShowServiceTestSuite) TestGetTVShowByName() {
	suite.client = &mocks.MockClient{}
	suite.tvShowService = service.NewTVShowService(suite.client, suite.tvShowRepo)
	testCases := []struct {
		name           string
		tvShowName     string
		expectedTVShow *models.TVShow
		wantErr        bool
	}{
		{
			name:       "Success",
			tvShowName: "BoJack Horseman",
			expectedTVShow: &models.TVShow{
				Name:      "BoJack Horseman",
				Language:  "English",
				Genres:    []string{"Comedy", "Drama"},
				Runtime:   25,
				Premiered: "2014-08-22",
				Rating:    8.1,
				PosterURL: "https://static.tvmaze.com/uploads/images/original_untouched/236/590384.jpg",
				Summary:   "Meet the most beloved sitcom horse of the '90s, 20 years later.",
				DirPath:   "testdata/three_shows/BoJack Horseman",
			},
			wantErr: false,
		},
		{
			name:           "Non-existent tv show",
			tvShowName:     "The Office",
			expectedTVShow: nil,
			wantErr:        true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			tvShow, err := suite.tvShowService.GetTVShowByName(tt.tvShowName)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			suite.Equal(tt.expectedTVShow, tvShow)
		})
	}
}

func (suite *TVShowServiceTestSuite) TestGetAllTVShows() {
	suite.client = &mocks.MockClient{}
	suite.tvShowService = service.NewTVShowService(suite.client, suite.tvShowRepo)

	testCases := []struct {
		name            string
		expectedTVShows []*models.TVShow
		wantErr         bool
	}{
		{
			name: "Success",
			expectedTVShows: []*models.TVShow{
				&models.TVShow{
					Name:      "BoJack Horseman",
					Language:  "English",
					Genres:    []string{"Comedy", "Drama"},
					Runtime:   25,
					Premiered: "2014-08-22",
					Rating:    8.1,
					PosterURL: "https://static.tvmaze.com/uploads/images/original_untouched/236/590384.jpg",
					Summary:   "Meet the most beloved sitcom horse of the '90s, 20 years later.",
					DirPath:   "testdata/three_shows/BoJack Horseman",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			tvShows, err := suite.tvShowService.GetAllTVShows()
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			suite.Equal(tt.expectedTVShows, tvShows)
		})
	}
}
