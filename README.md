<p align="center"><img width="80%" src="https://imgur.com/O4m6Y0Y.png"></p>

<h4 align="center">Share videos in a fast, simple and nice way.</h4>

## :warning: UPDATE :warning:
I am currently rewriting this to the microservices architecture, so it doesn't work.Sorry for the inconvenience.

## How To Use
To clone and run this application, you'll need [git](https://git-scm.com/), [docker](https://docker.com/) and [docker-compose](https://docs.docker.com/compose/install/). Follow this steps to make this work:

1. Clone this repository and create `.env` file.
	```
	$ git clone https://github.com/0x113/x-media
	$ cd x-media
	$ touch .env
	```

2. Set envirionment varaibles in `.env` file:
	If you're using vim, type:
	`$ vi .env`
	And if you're using emacs, enter the following:
	`$ sudo apt-get remove emacs`
	`vi .env`
	* `VIDEO_DIR` - path to folder with mp4 files
	* `MOVIES_SUB_DIR` - path to folder with subtitles for videos
	* `JWT_SECRET` - key for generating jwt
	* `MYSQL_ROOT_PASSWORD` - MySQL root password

	Example:
	```
	VIDEO_DIR=/home/user/Movies
	MOVIES_SUB_DIR=/home/user/Movies/sub
	JWT_SECRET=TOP SECRET
	MYSQL_ROOT_PASSWORD=rootorsomethingbetter
	```

3. Run it.
	```
	docker-compose up
	```

## API Endpoints

#### Authentication
* `/user/create`
	* Create new user
	* fields: `username`, `password`
	* method: `POST`
	* response: `201`

* `/user/token/generate`
	* Generate jwt for user
	* fields: `username`, `password`
	* method: `POST`
	* reponse: `200`

#### Video
* `/api/movies/update`
  * Updates movie database
  * method: `GET`
  * response: `200`
  * authentication required

* `/api/movies/`
  * List of all movies
  * method: `GET`
  * response: `200`
  * authentication required


