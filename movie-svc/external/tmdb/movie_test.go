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

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// TMDbAPIClientTestSuite defines the test suite
type TMDbAPIClientTestSuite struct {
	suite.Suite
}

// SetupTest initiates the fake api key and disables the logrus output
func (suite *TMDbAPIClientTestSuite) SetupTest() {
	common.Config = &common.Configuration{
		TMDbAPIKey: "fake-key",
	}
	logrus.SetOutput(ioutil.Discard)
}

// TestTMDbAPIClientTestSuite runs the test suite
func TestTMDbAPIClientTestSuite(t *testing.T) {
	suite.Run(t, new(TMDbAPIClientTestSuite))
}

func (suite *TMDbAPIClientTestSuite) TestGetTMDbQueryMovieInfo() {
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
			name: "HTTP client error",
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Body: ioutil.NopCloser(nil),
				}, errors.New("HTTP client error")
			},
			wantErr: true,
		},
		{
			name: "Wrong response status code",
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
		{
			name: "Decoding error",
			DoFunc: func(req *http.Request) (*http.Response, error) {
				json := `this should be json`
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
			movieQ, err := tmdbApiClient.GetTMDbQueryMovieInfo("Heat", "en")
			if tt.wantErr {
				suite.NotNil(err)
				suite.Nil(movieQ)
			} else {
				suite.Nil(err)
				suite.NotNil(movieQ)
			}
		})
	}
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
   "adult":false,
   "backdrop_path":"/rfEXNlql4CafRmtgp2VFQrBC4sh.jpg",
   "belongs_to_collection":null,
   "budget":60000000,
   "genres":[
      {
         "id":28,
         "name":"Action"
      },
      {
         "id":80,
         "name":"Crime"
      },
      {
         "id":18,
         "name":"Drama"
      },
      {
         "id":53,
         "name":"Thriller"
      }
   ],
   "homepage":"",
   "id":949,
   "imdb_id":"tt0113277",
   "original_language":"en",
   "original_title":"Heat",
   "overview":"Obsessive master thief, Neil McCauley leads a top-notch crew on various daring heists throughout Los Angeles while determined detective, Vincent Hanna pursues him without rest. Each man recognizes and respects the ability and the dedication of the other even though they are aware their cat-and-mouse game may end in violence.",
   "popularity":23.175,
   "poster_path":"/rrBuGu0Pjq7Y2BWSI6teGfZzviY.jpg",
   "production_companies":[
      {
         "id":508,
         "logo_path":"/7PzJdsLGlR7oW4J0J5Xcd0pHGRg.png",
         "name":"Regency Enterprises",
         "origin_country":"US"
      },
      {
         "id":675,
         "logo_path":null,
         "name":"Forward Pass",
         "origin_country":"US"
      },
      {
         "id":174,
         "logo_path":"/IuAlhI9eVC9Z8UQWOIDdWRKSEJ.png",
         "name":"Warner Bros. Pictures",
         "origin_country":"US"
      }
   ],
   "production_countries":[
      {
         "iso_3166_1":"US",
         "name":"United States of America"
      }
   ],
   "release_date":"1995-12-15",
   "revenue":187436818,
   "runtime":170,
   "spoken_languages":[
      {
         "iso_639_1":"en",
         "name":"English"
      },
      {
         "iso_639_1":"es",
         "name":"Español"
      }
   ],
   "status":"Released",
   "tagline":"A Los Angeles Crime Saga",
   "title":"Heat",
   "video":false,
   "vote_average":7.9,
   "vote_count":4110
}
`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
				}, nil
			},
			wantErr: false,
		},
		{
			name: "HTTP client error",
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Body: ioutil.NopCloser(nil),
				}, errors.New("HTTP client error")
			},
			wantErr: true,
		},
		{
			name: "Wrong response status code",
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       ioutil.NopCloser(nil),
				}, nil
			},
			wantErr: true,
		},
		{
			name: "Decoding error",
			DoFunc: func(req *http.Request) (*http.Response, error) {
				json := `this should be json`
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
			movie, err := tmdbApiClient.GetTMDbMovieInfo(949, "en")
			if tt.wantErr {
				suite.NotNil(err)
				suite.Nil(movie)
			} else {
				suite.Nil(err)
				suite.NotNil(movie)
				suite.Equal("Heat", movie.Title)
			}
		})
	}
}
