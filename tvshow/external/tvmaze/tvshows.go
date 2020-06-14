package tvmaze

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/0x113/x-media/tvshow/models"
	"github.com/0x113/x-media/tvshow/utils"

	log "github.com/sirupsen/logrus"
)

// TODO: remove logging and move to the service

// GetTVmazeTVShowInfo calls TVmaze api and returns new TVmaze object
func GetTVmazeTVShowInfo(client utils.HttpClient, title string) (*models.TVmazeTVShow, error) {
	query := url.QueryEscape(title)
	apiUrl := fmt.Sprintf("https://api.tvmaze.com/search/shows?q=%s", query)
	// request
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Debugf("Unable to prepare request for show [%s]; err: %v", title, err)
		return nil, err
	}
	// response
	res, err := client.Do(req)
	if err != nil {
		log.Debugf("Unable to send request to the TVmaze api; err: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	// check status code
	if res.StatusCode != http.StatusOK {
		log.Debugf("Expected status code: %d; got: %d", http.StatusOK, res.StatusCode)
		return nil, fmt.Errorf("Expected 200 status code, got %d", res.StatusCode)
	}
	// decode
	tvMazeResponse := []*models.TVmazeTVShow{}
	if err := json.NewDecoder(res.Body).Decode(&tvMazeResponse); err != nil {
		log.Debugf("Unable to decode TVmaze info for show[%s]; err: %v", title, err)
		return nil, err
	}

	var tvMazeInfo *models.TVmazeTVShow
	if len(tvMazeResponse) >= 1 {
		tvMazeInfo = tvMazeResponse[0] // get first match
	}

	return tvMazeInfo, nil

}
