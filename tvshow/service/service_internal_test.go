package service

import (
	"testing"

	"github.com/0x113/x-media/tvshow/common"

	"github.com/stretchr/testify/assert"
)

func TestGetDirectories(t *testing.T) {
	testCases := []struct {
		name         string
		tvShowDirs   []string
		expectedDirs []string
	}{
		{
			name:         "Success",
			tvShowDirs:   []string{"testdata/three_shows", "testdata/two_shows_one_file"},
			expectedDirs: []string{"BoJack Horseman", "The_Office", "Trailer.Park.Boys", "Rick and Morty", "The.Sopranos"},
		},
		{
			name:         "Non-existent directory",
			tvShowDirs:   []string{"testdata/no-such-dir"},
			expectedDirs: []string{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// set up config
			common.Config = &common.Configuration{
				TVShowDirectories: tt.tvShowDirs,
			}

			tvShowDirs := getDirectories()
			assert.Equal(t, tt.expectedDirs, tvShowDirs)
		})
	}
}

func TestCreateName(t *testing.T) {
	dirNames := map[string]string{
		"The_Office/":       "The Office",
		"Trailer.Park.Boys": "Trailer Park Boys",
		"Rick and Morty":    "Rick and Morty",
	}

	for name, expected := range dirNames {
		actual := createName(name)
		assert.Equal(t, expected, actual)
	}
}
