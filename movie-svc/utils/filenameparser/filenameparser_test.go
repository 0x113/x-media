package filenameparser_test

import (
	"testing"

	"github.com/0x113/x-media/movie-svc/utils/filenameparser"
	"github.com/stretchr/testify/assert"
)

func TestCreateTitle(t *testing.T) {
	testCases := []struct {
		name          string
		filepath      string
		expectedTitle string
		wantErr       bool
	}{
		{
			name:          "Slash at the end",
			filepath:      "/home/y0x/Videos/Heat.1995.NSB.mp4/",
			expectedTitle: "Heat",
			wantErr:       false,
		},
		{
			name:          "No slash at the end",
			filepath:      "/home/y0x/Videos/The Godfather 1972 H264 DVDRip Useless-Info.avi",
			expectedTitle: "The Godfather",
			wantErr:       false,
		},
		{
			name:          "Nested filepath",
			filepath:      "/home/y0x/Videos/[Useless.info]_The.Movie.Title.2005.Useless.avi",
			expectedTitle: "The Movie Title",
			wantErr:       false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			title, err := filenameparser.CreateTitle(tt.filepath)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.expectedTitle, title)
		})
	}
}
