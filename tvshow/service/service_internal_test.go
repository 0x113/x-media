package service

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/0x113/x-media/tvshow/common"

	"github.com/stretchr/testify/assert"
)

func TestGetDirectories(t *testing.T) {
	var expectedTempDirs []string

	// NOTE: this probably could be done in more elegant way
	// TODO: change in the future
	// create temp dirs
	threeShowsDir, err := ioutil.TempDir("", "three-shows-*")
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %v", err)
	}

	// create three temp dirs in the three-shows-dir
	for i := 1; i <= 3; i++ {
		tvTempDir, err := ioutil.TempDir(threeShowsDir, fmt.Sprintf("tvshow-%d-*", i))
		if err != nil {
			t.Fatalf("Unable to create temporary tv show directory: %v", err)
		}
		expectedTempDirs = append(expectedTempDirs, tvTempDir)
	}

	// two shows and one file temporary dir
	twoShowsOneFileDir, err := ioutil.TempDir("", "two-shows-one-file-*")
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %v", err)
	}

	// create two temp dirs in the two-shows-one-file directory
	for i := 1; i <= 2; i++ {
		tvTempDir, err := ioutil.TempDir(twoShowsOneFileDir, fmt.Sprintf("tvshow-%d-*", i))
		if err != nil {
			t.Fatalf("Unable to create temporary tv show directory: %v", err)
		}
		expectedTempDirs = append(expectedTempDirs, tvTempDir)
	}

	// create temp file in the two-shows-one-file directory
	_, err = ioutil.TempFile(twoShowsOneFileDir, "temp-file-*")
	if err != nil {
		t.Fatalf("Unable to create temporary file: %v", err)
	}

	testCases := []struct {
		name         string
		tvShowDirs   []string
		expectedDirs []string
	}{
		{
			name:         "Success",
			tvShowDirs:   []string{threeShowsDir, twoShowsOneFileDir},
			expectedDirs: expectedTempDirs,
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
