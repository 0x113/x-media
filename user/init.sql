CREATE DATABASE IF NOT EXISTS xmedia_users;
USE xmedia_users;

CREATE TABLE IF NOT EXISTS user (
  user_id int NOT NULL AUTO_INCREMENT,
  username varchar(255) NOT NULL UNIQUE,
  password varchar(60) NOT NULL,
  is_admin BOOLEAN NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY(user_id)
);

