package service_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	"github.com/0x113/x-media/movie-svc/common"
	"github.com/0x113/x-media/movie-svc/httpclient"
	"github.com/0x113/x-media/movie-svc/mocks"
	"github.com/0x113/x-media/movie-svc/models"
	"github.com/0x113/x-media/movie-svc/service"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// MovieServiceTestSuite represents test suite for the movie service
type MovieServiceTestSuite struct {
	suite.Suite
	httpClient   httpclient.HTTPClient
	movieID      primitive.ObjectID
	movieRepo    *mocks.MockMovieRepository
	movieService service.MovieService
}

// SetupTest initiates mocked database and disables the logrus output
func (suite *MovieServiceTestSuite) SetupTest() {
	id, err := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011") // must be same like in the mocked database
	suite.Nil(err)
	suite.movieID = id
	common.Config = &common.Configuration{
		TMDbAPIKey: "fake-key",
	}
	suite.movieRepo = mocks.NewMockMovieRepository()
	logrus.SetOutput(ioutil.Discard)
}

// TestMovieServiceTestSuite runs the test suite for the movie service
func TestMoviceServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MovieServiceTestSuite))
}

func (suite *MovieServiceTestSuite) TestUpdateMovieByID() {
	testCases := []struct {
		name     string
		id       int
		lang     string
		filePath string
		doFunc   func(req *http.Request) (*http.Response, error)
		wantErr  bool
	}{
		{
			name:     "Existing movie - update",
			id:       949,
			lang:     "en",
			filePath: "/home/y0x/Videos/Heat.mp4",
			doFunc: func(req *http.Request) (*http.Response, error) {
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
}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
				}, nil
			},
			wantErr: false,
		},
		{
			name: "Failure - unable to get data from the TMDb API",
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       ioutil.NopCloser(nil),
				}, nil
			},
			wantErr: true,
		},
		{
			name: "New movie - save into the database",
			id:   524,
			doFunc: func(req *http.Request) (*http.Response, error) {
				json := `{
   "adult":false,
   "backdrop_path":"/pLR2O3dmA9xkCiPA26U7ErCUYSi.jpg",
   "belongs_to_collection":null,
   "budget":52000000,
   "genres":[
      {
         "id":80,
         "name":"Crime"
      }
   ],
   "homepage":"",
   "id":524,
   "imdb_id":"tt0112641",
   "original_language":"en",
   "original_title":"Casino",
   "overview":"In early-1970s Las Vegas, low-level mobster Sam \"Ace\" Rothstein gets tapped by his bosses to head the Tangiers Casino. At first, he's a great success in the job, but over the years, problems with his loose-cannon enforcer Nicky Santoro, his ex-hustler wife Ginger, her con-artist ex Lester Diamond and a handful of corrupt politicians put Sam in ever-increasing danger.",
   "popularity":18.997,
   "poster_path":"/4TS5O1IP42bY2BvgMxL156EENy.jpg",
   "production_companies":[
      {
         "id":33,
         "logo_path":"/8lvHyhjr8oUKOOy2dKXoALWKdp0.png",
         "name":"Universal Pictures",
         "origin_country":"US"
      },
      {
         "id":11583,
         "logo_path":null,
         "name":"Syalis DA",
         "origin_country":""
      },
      {
         "id":10898,
         "logo_path":null,
         "name":"Légende Entreprises",
         "origin_country":""
      },
      {
         "id":11584,
         "logo_path":null,
         "name":"De Fina-Cappa",
         "origin_country":""
      }
   ],
   "production_countries":[
      {
         "iso_3166_1":"FR",
         "name":"France"
      },
      {
         "iso_3166_1":"US",
         "name":"United States of America"
      }
   ],
   "release_date":"1995-11-22",
   "revenue":116112375,
   "runtime":179,
   "spoken_languages":[
      {
         "iso_639_1":"en",
         "name":"English"
      }
   ],
   "status":"Released",
   "tagline":"No one stays at the top forever.",
   "title":"Casino",
   "video":false,
   "vote_average":8.0,
   "vote_count":3304
}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(json))),
				}, nil
			},
			lang:     "en",
			filePath: "/home/y0x/Videos/Heat.mp4",
			wantErr:  false,
		},
	}

	var mutex sync.Mutex
	for _, tt := range testCases {
		suite.httpClient = &mocks.MockClient{tt.doFunc}
		suite.movieService = service.NewMovieService(suite.movieRepo, suite.httpClient)
		suite.Run(tt.name, func() {
			_, err := suite.movieService.UpdateMovieByID(tt.id, tt.lang, tt.filePath, &mutex) // NOTE: handle movie return
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *MovieServiceTestSuite) TestGetLocalTMDbID() {
	testCases := []struct {
		name       string
		filename   string
		expectedID int
		doFunc     func(req *http.Request) (*http.Response, error)
		wantErr    bool
	}{
		{
			name:       "Success - Heat",
			filename:   "Heat.mp4",
			expectedID: 949,
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
			name:       "Failure - unable to get info from the TMDb API",
			filename:   "test.mp4",
			expectedID: 0,
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       ioutil.NopCloser(nil),
				}, nil
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		suite.httpClient = &mocks.MockClient{tt.doFunc}
		suite.movieService = service.NewMovieService(suite.movieRepo, suite.httpClient)
		suite.Run(tt.name, func() {
			id, err := suite.movieService.GetLocalTMDbID(tt.filename)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			suite.Equal(tt.expectedID, id)
		})
	}

}

func (suite *MovieServiceTestSuite) TestUpdateAllMovies() {
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

	testCases := []struct {
		name           string
		expectedMovies map[string]string
		expectedErrors map[string]string
		doFunc         func(req *http.Request) (*http.Response, error)
	}{
		{
			name: "Success",
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
		},
	}

	for _, tt := range testCases {
		suite.httpClient = &mocks.MockClient{tt.doFunc}
		suite.movieService = service.NewMovieService(suite.movieRepo, suite.httpClient)
		suite.Run(tt.name, func() {
			updatedMovies, errors := suite.movieService.UpdateAllMovies("en")
			suite.NotNil(errors)
			suite.NotNil(updatedMovies)
		})
	}
}

func (suite *MovieServiceTestSuite) TestGetAll() {
	suite.httpClient = &mocks.MockClient{}
	suite.movieService = service.NewMovieService(suite.movieRepo, suite.httpClient)

	expectedMovies := []*models.Movie{
		&models.Movie{
			ID:               suite.movieID,
			TMDbID:           949,
			Title:            "Heat",
			Overview:         "Obsessive master thief, Neil McCauley leads a top-notch crew on various daring heists throughout Los Angeles while determined detective, Vincent Hanna pursues him without rest. Each man recognizes and respects the ability and the dedication of the other even though they are aware their cat-and-mouse game may end in violence.",
			OriginalTitle:    "Heat",
			OriginalLanguage: "en",
			ReleaseDate:      "1995-12-15",
			Genres: []string{
				"Action",
				"Crime",
				"Drama",
				"Thriller",
			},
			Rating:       7.9,
			Runtime:      170,
			BackdropPath: "/rfEXNlql4CafRmtgp2VFQrBC4sh.jpg",
			PosterPath:   "/rrBuGu0Pjq7Y2BWSI6teGfZzviY.jpg",
			DirPath:      "/home/y0x/Videos/Heat.1995.mp4",
		},
	}

	movies, err := suite.movieService.GetAllMovies()
	suite.Nil(err)
	suite.Equal(expectedMovies, movies)
}

func (suite *MovieServiceTestSuite) TestGetMovieByID() {
	testCases := []struct {
		name          string
		id            string
		expectedMovie *models.Movie
		wantErr       bool
	}{
		{
			name: "Success",
			id:   "507f1f77bcf86cd799439011",
			expectedMovie: &models.Movie{
				ID:               suite.movieID,
				TMDbID:           949,
				Title:            "Heat",
				Overview:         "Obsessive master thief, Neil McCauley leads a top-notch crew on various daring heists throughout Los Angeles while determined detective, Vincent Hanna pursues him without rest. Each man recognizes and respects the ability and the dedication of the other even though they are aware their cat-and-mouse game may end in violence.",
				OriginalTitle:    "Heat",
				OriginalLanguage: "en",
				ReleaseDate:      "1995-12-15",
				Genres: []string{
					"Action",
					"Crime",
					"Drama",
					"Thriller",
				},
				Rating:       7.9,
				Runtime:      170,
				BackdropPath: "/rfEXNlql4CafRmtgp2VFQrBC4sh.jpg",
				PosterPath:   "/rrBuGu0Pjq7Y2BWSI6teGfZzviY.jpg",
				DirPath:      "/home/y0x/Videos/Heat.1995.mp4",
			},
			wantErr: false,
		},
		{
			name:          "Incorrect object id",
			id:            "123",
			expectedMovie: nil,
			wantErr:       true,
		},
		{
			name:          "No movie with provided id in the database",
			id:            "507f1f77bcf86cd799439010",
			expectedMovie: nil,
			wantErr:       true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			movie, err := suite.movieService.GetMovieByID(tt.id)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			suite.Equal(tt.expectedMovie, movie)
		})
	}
}
