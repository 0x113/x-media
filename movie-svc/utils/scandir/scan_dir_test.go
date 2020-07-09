package scandir_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/0x113/x-media/movie-svc/utils/scandir"
	"github.com/stretchr/testify/assert"
)

func TestGetFiles(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "get-files-test")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpdir)

	// create temp mp4 file
	tmpMp4File, err := ioutil.TempFile(tmpdir, "video-file.*.mp4")
	assert.Nil(t, err)

	// test all files
	files, err := scandir.GetFiles(tmpdir, nil)
	assert.Nil(t, err)
	assert.Equal(t, files, []string{tmpMp4File.Name()})

	// non existent dir
	files, err = scandir.GetFiles("./no-such-dir/", nil)
	assert.NotNil(t, err)
	assert.Nil(t, files)

	// test files by extensions
	tmpTxtFile, err := ioutil.TempFile(tmpdir, "txt-file.*.txt")
	assert.Nil(t, err)

	files, err = scandir.GetFiles(tmpdir, []string{".txt"})
	assert.Nil(t, err)
	assert.Equal(t, files, []string{tmpTxtFile.Name()})

	// test no files in the given directory
	files, err = scandir.GetFiles(tmpdir, []string{".te"})
	assert.Nil(t, err)
	assert.Equal(t, []string(nil), files)
}
