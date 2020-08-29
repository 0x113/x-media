<p align="center"><img src="https://imgur.com/jEZ0701.png"></p>
<p align="center">
<img src="https://travis-ci.com/0x113/x-media.svg?branch=master">
</p>

## How To Use
To clone and run this application, you'll need [git](https://git-scm.com/), [docker](https://docker.com/) and [docker-compose](https://docs.docker.com/compose/install/). Follow this steps to make this work:

1. Clone this repository.
	```
	$ git clone https://github.com/0x113/x-media
	$ cd x-media
	```
2. Run it via docker.
	```
	docker-compose up
	```

## Services

#### TV shows service
* Runs at port `8001`

#### User service
* Runs at port `8002`

#### Authentication service
* Runs at port `8003`

#### Movie service
* Runs at port `8004`

## Config
Config files are avaiable at `<service-name>/config/config.json`

### List of required config changes

#### TV shows service
* `tv_show_directories` - currently only one directory is supported, must be same like in the `docker-compose.yml`. Have idea how to fix it, but don't have time :smile:

#### Authentication service
`access_secret` - secret key to generate authentication token
`refresh_secret` - secret key to generate refresh token

#### Movie service
* `tmdb_api_key` - API key for the [TMDb](https://www.themoviedb.org/)
* `movie_directories` - currently only one directory is supported, must be same like in the `docker-compose.yml`. Have idea how to fix it, but don't have time :smile:

## API docs
Generated using [swagger](https://swagger.io/).<br>
You can read the docs for each service at the `localhost:<service-port>/docs`. <br>
Example: `localhost:8001/docs` will give you docs for the tv show service.
<p align="center"><img src="https://imgur.com/5oJIjjD.png"></p>

## Reverse proxy
Using [treafik](https://github.com/containous/traefik/).<br>
Visit `localhost:8080` for the dashboard.
<p align="center"><img src="https://imgur.com/4e46r50.png"></p>
<br>

#### Calling services
You can call each service separatly with `curl --header "Host: <hostname>" localhost/<api-endpoint>`
Example: `curl --header "Host: tvshowsvc" localhost/api/v1/tvshows/get/all`

#### Hosts
* User service -> `usersvc`
* Authentication service -> `authsvc`
* Movie service -> `moviesvc`
* TV shows service -> `tvshowsvc`

## Frontend
Work in progress :). <br>
Currenly working on some themes like 8bit or some modern one, so it is in a design phase.
If you'd like to create frontend for this project, feel free to contribute.

## Third party APIs
[TMDb](https://www.themoviedb.org) - getting data about movies<br><br>
<img src="https://www.themoviedb.org/assets/2/v4/logos/v2/blue_square_2-d537fb228cf3ded904ef09b136fe3fec72548ebc1fea3fbbd1ad9e36364db38b.svg" alt="tmdb-logo" width="20%"/>

<hr>

[TVmaze](https://www.tvmaze.com) - getting data about tv shows<br><br>
<img src="https://static.tvmaze.com/images/tvm-header-logo.png" alt="tmdb-logo" width="20%"/>
