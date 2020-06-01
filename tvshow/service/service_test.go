package service_test

import (
	"io/ioutil"
	"testing"

	"github.com/0x113/x-media/tvshow/mocks"
	"github.com/0x113/x-media/tvshow/models"
	"github.com/0x113/x-media/tvshow/service"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// TVShowServiceTestSuite define the suite
type TVShowServiceTestSuite struct {
	suite.Suite
	tvShowRepo    *mocks.MockTVShowRepository
	tvShowService service.TVShowService
}

// SetupTest initiates mocked database and new tv show service
func (suite *TVShowServiceTestSuite) SetupTest() {
	suite.tvShowRepo = mocks.NewMockTVShowRepository()
	suite.tvShowService = service.NewTVShowService(suite.tvShowRepo)
	logrus.SetOutput(ioutil.Discard) // disable logrus
}

// TestTVShowServiceTestSuite runs test suite
func TestTVShowServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TVShowServiceTestSuite))
}

func (suite *TVShowServiceTestSuite) TestSave() {
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
			},
			wantErr: false,
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
