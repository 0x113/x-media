package tmdb_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/0x113/x-media/movie-svc/common"
	"github.com/0x113/x-media/movie-svc/external/tmdb"
	"github.com/0x113/x-media/movie-svc/mocks"

	"github.com/stretchr/testify/suite"
)

// TMDbAPIClientTestSuite defines the test suite
type TMDbAPIClientTestSuite struct {
	suite.Suite
}

// SetupTest initiates the fake api key
func (suite *TMDbAPIClientTestSuite) SetupTest() {
	common.Config = &common.Configuration{
		TMDbAPIKey: "fake-key",
	}
}

// TestTMDbAPIClientTestSuite runs the test suite
func TestTMDbAPIClientTestSuite(t *testing.T) {
	suite.Run(t, new(TMDbAPIClientTestSuite))
}

func (suite *TMDbAPIClientTestSuite) TestGetTMDbMovieInfo() {
	testCases := []struct {
		name    string
		DoFunc  func(req *http.Request) (*http.Response, error)
		wantErr bool
	}{
		{
			name: "Success",
			DoFunc: func(req *http.Request) (*http.Response, error) {
				json := `{
   "page":1,
   "total_results":327,
   "total_pages":17,
   "results":[
      {
         "popularity":24.521,
         "id":949,
         "video":false,
         "vote_count":4103,
         "vote_average":7.9,
         "title":"Heat",
         "release_date":"1995-12-15",
         "original_language":"en",
         "original_title":"Heat",
         "genre_ids":[
            28,
            80,
            18,
            53
         ],
         "backdrop_path":"/rfEXNlql4CafRmtgp2VFQrBC4sh.jpg",
         "adult":false,
         "overview":"Obsessive master thief, Neil McCauley leads a top-notch crew on various daring heists throughout Los Angeles while determined detective, Vincent Hanna pursues him without rest. Each man recognizes and respects the ability and the dedication of the other even though they are aware their cat-and-mouse game may end in violence.",
         "poster_path":"/rrBuGu0Pjq7Y2BWSI6teGfZzviY.jpg"
      }
   ]
}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
				}, nil
			},
			wantErr: false,
		},
		{
			name: "Client err",
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Body: ioutil.NopCloser(nil),
				}, errors.New("Client error")
			},
			wantErr: true,
		},
		{
			name: "Unexpected response status code",
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       ioutil.NopCloser(nil),
				}, nil
			},
			wantErr: true,
		},
		{
			name: "Empty results list",
			DoFunc: func(req *http.Request) (*http.Response, error) {
				json := `{
   "page":1,
   "total_results":327,
   "total_pages":17,
   "results":[]}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
				}, nil
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		client := &mocks.MockClient{tt.DoFunc}
		tmdbApiClient := &tmdb.TMDbAPIClient{client}
		suite.Run(tt.name, func() {
			movieQ, err := tmdbApiClient.GetTMDbMovieInfo("Heat", "en")
			if tt.wantErr {
				suite.NotNil(err)
				suite.Nil(movieQ)
			} else {
				suite.Nil(err)
				suite.NotNil(movieQ)
				suite.Equal(movieQ.Title, "Heat")
			}
		})
	}
}
