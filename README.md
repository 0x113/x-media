<p align="center"><img src="https://imgur.com/jEZ0701.png"></p>

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
Work in progress :). Currenly working on some themes like 8bit or some modern one, so it is in a design phase.
If you'd like to create frontend for this project, feel free to contribute.
