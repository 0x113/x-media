package scandir

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// GetFiles scans all of the files within the directory given
// by the file extension. If extensions are nil it returns
// all files from the given directory.
//
// GetFiles returns list of file paths and error.
func GetFiles(dir string, extensions []string) ([]string, error) {
	// check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("Directory %s does not exist", dir)
	}
	// append slash to dir name if doesn't have
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filePaths []string
	for _, f := range files {
		if extensions == nil {
			filePaths = append(filePaths, dir+f.Name())
			continue
		} else if !f.IsDir() && hasSuffix(f.Name(), extensions) {
			filePaths = append(filePaths, dir+f.Name())
		}
	}

	return filePaths, nil
}

// hasSuffix checks if the given string has one off the
// given suffixes
func hasSuffix(str string, arr []string) bool {
	for _, s := range arr {
		if strings.HasSuffix(str, s) {
			return true
		}
	}
	return false
}
