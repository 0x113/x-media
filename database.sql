CREATE DATABASE IF NOT EXISTS xmedia;

USE xmedia;

CREATE TABLE IF NOT EXISTS movie (
    movie_id int NOT NULL AUTO_INCREMENT,
    title varchar(255) NOT NULL,
    description MEDIUMTEXT NOT NULL,
    director varchar(255) NOT NULL,
    genre varchar(80) NOT NULL,
    duration varchar(80) NOT NULL,
    rate FLOAT NOT NULL,
    release_date DATE NOT NULL,
    file_name varchar(255) NOT NULL UNIQUE,
    poster_path varchar(255) NOT NULL,
    PRIMARY KEY (movie_id)
);


CREATE TABLE IF NOT EXISTS series (
    series_id int NOT NULL AUTO_INCREMENT,
    title varchar(255) NOT NULL UNIQUE,
    description MEDIUMTEXT NOT NULL,
    director varchar(255) NOT NULL,
    genre varchar(80) NOT NULL,
    episode_duration varchar(80) NOT NULL,
    rate FLOAT NOT NULL,
    release_date DATE NOT NULL,
    dir_name varchar(255) NOT NULL UNIQUE,
    poster_path varchar(255),
    PRIMARY KEY (series_id)
);

