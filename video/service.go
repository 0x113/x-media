package video

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
	log "github.com/sirupsen/logrus"
)

type VideoService interface {
	Save() error
	AllMovies() ([]*Movie, error)
}

type videoService struct {
	repo VideoRepository
}

func NewVideoService(repo VideoRepository) VideoService {
	return &videoService{
		repo,
	}
}

func (s *videoService) Save() error {
	log.Infoln("Updating movie database...")
	videos, err := s.getVideos("/home/xa0s/Downloads/Movies/") // TODO: use env variable
	if err != nil {
		log.Error("Unable to get movie info")
		return err
	}
	for _, v := range videos {
		movie, err := s.getMovieInfo(v)
		if err != nil || movie == nil {
			continue
		}
		s.repo.SaveMovie(movie)
	}
	log.Infoln("The movie database has been updated.")
	return nil
}

func (s *videoService) AllMovies() ([]*Movie, error) {
	return s.repo.FindAllMovies()
}

func (s *videoService) getVideos(videoDirPath string) ([]string, error) {

	/* Get movies from disk (mkv & mp4 files).*/
	var videos []string
	files, err := ioutil.ReadDir(videoDirPath)
	if err != nil {
		return videos, err
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "mkv") || strings.HasSuffix(f.Name(), "mp4") {
			videos = append(videos, f.Name())
		}
	}

	return videos, nil
}

func (s *videoService) getMovieInfo(movieFileName string) (*Movie, error) {
	toRemove := []string{".NSB", ".mkv", ".mp4"}
	var toSearch = s.removeFromArray(movieFileName, toRemove)

	/* Get movie info from filmweb.pl TODO: allow user to choose other service*/
	url := "https://filmweb.pl/search?q=" + toSearch

	url = strings.Replace(url, ".mp4", "", -1)

	res, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(res)

	/* Get movie card and check for errors. */
	movieCard := doc.Find("div", "class", "wrapperContent__content")
	if movieCard.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie card")
		return nil, err
	}
	movieCard = movieCard.Find("ul", "class", "resultsList")
	if movieCard.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find results list")
		return nil, err
	}
	movieCard = movieCard.Find("li")
	if movieCard.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie")
		return nil, err
	}

	/* Get movie title */
	titleHTML := movieCard.Find("data")
	if titleHTML.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie title")
		return nil, err
	}
	title := titleHTML.Attrs()["data-title"]

	/* Get movie release date */
	movieReleaseDateHTML := movieCard.Find("div")
	if movieReleaseDateHTML.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie release date")
		return nil, err
	}
	movieReleaseDate := movieReleaseDateHTML.Attrs()["data-release"]

	/* Get movie duration */
	movieDurationHTML := movieCard.Find("div", "class", "filmPreview__filmTime")
	if movieDurationHTML.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie duration")
		return nil, err
	}
	movieDuration := movieDurationHTML.Text()

	/* Get movie rate */
	movieRateHTML := movieCard.Find("div", "class", "filmPreview__rateBox")
	if movieRateHTML.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie rate")
		return nil, err
	}
	movieRate := movieRateHTML.Attrs()["data-rate"]
	// convert movie rate to float
	movieRateFloat, err := strconv.ParseFloat(movieRate, 64)
	if err != nil {
		return nil, err
	}

	/* Get movie director */
	movieDirectorHTML := movieCard.Find("div", "class", "filmPreview__info--directors")
	if movieDurationHTML.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie directors")
		return nil, err
	}
	movieDirectorHTML = movieDirectorHTML.Find("a")
	if movieDurationHTML.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie directors")
		return nil, err
	}
	movieDirector := movieDirectorHTML.Attrs()["title"]

	/* Get movie genre */
	movieGenreHTML := movieCard.Find("div", "class", "filmPreview__info--genres")
	if movieGenreHTML.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie genre")
		return nil, err
	}
	movieGenreHTML = movieGenreHTML.Find("ul")
	if movieGenreHTML.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie genre")
		return nil, err
	}
	movieGenreHTML = movieGenreHTML.Find("a")
	if movieGenreHTML.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie genre")
		return nil, err
	}
	movieGenre := movieGenreHTML.Text()

	/* Get movie poster */
	moviePosterHTML := movieCard.Find("img", "class", "filmPoster__image")
	if moviePosterHTML.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie poster")
		return nil, err
	}
	moviePoster := moviePosterHTML.Attrs()["data-src"]
	moviePoster = strings.Replace(moviePoster, "6.jpg", "3.jpg", -1)

	/* Get movie details (description) */
	detailsLinkHTML := movieCard.Find("a", "class", "filmPreview__link")
	if detailsLinkHTML.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find details link")
		return nil, err
	}
	// Scrape details page
	detailsURL := detailsLinkHTML.Attrs()["href"]
	detailsRes, err := soup.Get("https://filmweb.pl" + detailsURL)
	if err != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Unable to open details page")
		return nil, err
	}
	detailsDoc := soup.HTMLParse(detailsRes)
	// Get movie description
	descriptionHTML := detailsDoc.Find("div", "class", "filmPlot")
	if descriptionHTML.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie description")
		return nil, err
	}
	descriptionHTML = descriptionHTML.Find("p")
	if descriptionHTML.Error != nil {
		log.WithFields(log.Fields{"movie": toSearch}).Error("Cannot find movie description")
		return nil, err
	}
	description := descriptionHTML.Text()

	movie := Movie{
		Title:       title,
		Description: description,
		Director:    movieDirector,
		Genre:       movieGenre,
		Duration:    movieDuration,
		Rate:        movieRateFloat,
		ReleaseDate: movieReleaseDate,
		PosterPath:  moviePoster,
	}
	return &movie, nil
}

func (s *videoService) removeFromArray(str string, toRemove []string) string {
	for _, x := range toRemove {
		if strings.Contains(str, x) {
			str = strings.Replace(str, x, "", -1)
		}
	}
	return str
}
