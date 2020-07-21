package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0x113/x-media/movie-svc/common"
	"github.com/0x113/x-media/movie-svc/httpclient"
	"github.com/0x113/x-media/movie-svc/mocks"
	"github.com/0x113/x-media/movie-svc/service"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// MovieHandlerTestSuite defines the test suite for the movie handler
type MovieHandlerTestSuite struct {
	suite.Suite
	router          *echo.Echo
	httpClient      httpclient.HTTPClient
	movieRepository *mocks.MockMovieRepository
	movieService    service.MovieService
}

// SetupTest initiates mocked database, config, router and disables the logrus output
func (suite *MovieHandlerTestSuite) SetupTest() {
	common.Config = &common.Configuration{
		TMDbAPIKey: "fake-key",
	}
	suite.router = echo.New()
	suite.movieRepository = mocks.NewMockMovieRepository()
	logrus.SetOutput(ioutil.Discard)

	// create temporary directries and files
	tmpdir, err := ioutil.TempDir("", "update-all-test-1-*")
	suite.Nil(err)

	tmpdir2, err := ioutil.TempDir("", "update-all-test-2-*")
	suite.Nil(err)

	_, err = ioutil.TempFile(tmpdir, "The.Godfather-*.mp4")
	suite.Nil(err)
	_, err = ioutil.TempFile(tmpdir, "Heat.1995-*.mkv")
	suite.Nil(err)

	_, err = ioutil.TempFile(tmpdir2, "Inception.2010-*.mkv")
	suite.Nil(err)
	_, err = ioutil.TempFile(tmpdir2, "txt-file.*.txt")

	common.Config = &common.Configuration{
		MovieDirectories: []string{tmpdir, tmpdir2},
	}
}

// TestMovieHandlerTestSuite runs the test suite
func TestMovieHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(MovieHandlerTestSuite))
}

func (suite *MovieHandlerTestSuite) TestUpdateAllMovies() {
	testCases := []struct {
		name               string
		json               string
		expectedStatusCode int
		doFunc             func(req *http.Request) (*http.Response, error)
		wantErr            bool
	}{
		{
			name:               "Success; en-language",
			json:               `{"language": "en"}`,
			expectedStatusCode: 200,
			doFunc: func(req *http.Request) (*http.Response, error) {
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
			name:               "Binding error; empty request body",
			json:               ``,
			expectedStatusCode: 400,
			wantErr:            true,
		},
	}

	for _, tt := range testCases {
		suite.httpClient = &mocks.MockClient{tt.doFunc}
		suite.movieService = service.NewMovieService(suite.movieRepository, suite.httpClient)
		h := movieHandler{suite.movieService}

		suite.Run(tt.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/movies/update/all", strings.NewReader(tt.json))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := suite.router.NewContext(req, rec)

			err := h.UpdateAllMovies(c)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			suite.Equal(tt.expectedStatusCode, rec.Code)
		})
	}
}

func (suite *MovieHandlerTestSuite) TestGetAllMovies() {
	// setup
	suite.httpClient = &mocks.MockClient{}
	suite.movieService = service.NewMovieService(suite.movieRepository, suite.httpClient)
	h := movieHandler{suite.movieService}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/movies/get/all", nil)
	rec := httptest.NewRecorder()
	c := suite.router.NewContext(req, rec)

	err := h.GetAllMovies(c)
	suite.Nil(err)
	suite.Equal(200, rec.Code)
}

func (suite *MovieHandlerTestSuite) TestGetMovieByID() {
	// setup
	suite.httpClient = &mocks.MockClient{}
	suite.movieService = service.NewMovieService(suite.movieRepository, suite.httpClient)
	h := movieHandler{suite.movieService}

	testCases := []struct {
		name               string
		id                 string
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:               "Success",
			id:                 "507f1f77bcf86cd799439011",
			expectedStatusCode: 200,
			wantErr:            false,
		},
		{
			name:               "Service error; non-existent movie or incorrect object id",
			id:                 "507f1f77bcf86cd799439010",
			expectedStatusCode: 500,
			wantErr:            true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := suite.router.NewContext(req, rec)
			c.SetPath("/api/v1/movies/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			err := h.GetMovieByID(c)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			suite.Equal(tt.expectedStatusCode, rec.Code)
		})
	}
}
