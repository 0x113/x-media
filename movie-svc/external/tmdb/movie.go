package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/0x113/x-media/movie-svc/common"
	"github.com/0x113/x-media/movie-svc/httpclient"
	"github.com/0x113/x-media/movie-svc/models"
)

// TMDbAPIClient contains method to operate with the TMDb API
type TMDbAPIClient struct {
	Client httpclient.HTTPClient
}

// GetTMDbMovieInfo calls the TMDb API and returns new movie info
func (t *TMDbAPIClient) GetTMDbMovieInfo(title, lang string) (*models.TMDbQueryMovie, error) {
	queryTitle := url.QueryEscape(title)
	apiUrl := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s&language=%s", common.Config.TMDbAPIKey, queryTitle, lang)
	// request
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}
	// response
	res, err := t.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Couldn't get movie info: wrong status code; wanted %d, got %d", http.StatusOK, res.StatusCode)
	}

	// decode the response
	tmdbQueryRes := new(models.TMDbQueryResponse)
	if err := json.NewDecoder(res.Body).Decode(tmdbQueryRes); err != nil {
		return nil, err
	}
	// get title of the first result
	if len(tmdbQueryRes.Results) == 0 {
		return nil, fmt.Errorf("Unable to find movies for title: %s", title)
	}

	return tmdbQueryRes.Results[0], nil
}

// GetTMDbGenres calls the TMDb API and returns list of genres
func (t *TMDbAPIClient) GetTMDbGenres(lang string) ([]*models.TMDbGenre, error) {
	apiUrl := fmt.Sprintf("https://api.themoviedb.org/3/genre/movie/list?api_key=%s&language=%s", common.Config.TMDbAPIKey, lang)
	// request
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}
	// response
	res, err := t.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Couldn't get genres: wrong status code; wanted %d, got %d", http.StatusOK, res.StatusCode)
	}

	// decode the response
	var response struct {
		Genres []*models.TMDbGenre `json:"genres"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Genres, nil
}
