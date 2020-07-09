package filenameparser

import (
	"strings"

	parsetorrentname "github.com/middelink/go-parse-torrent-name"
)

// CreateTitle extracts the file name from the given file path.
// Then it tries to parse name to extract movie title.
func CreateTitle(filepath string) (string, error) {
	var filename string
	pathSlice := strings.Split(filepath, "/")

	if strings.HasSuffix(filepath, "/") {
		filename = pathSlice[len(pathSlice)-2]
	} else {
		filename = pathSlice[len(pathSlice)-1]
	}

	info, err := parsetorrentname.Parse(filename)
	if err != nil {
		return "", err
	}

	return info.Title, nil
}
