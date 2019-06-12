#!/bin/bash

env db_user="xmedia_user" db_pass="password" db_host="127.0.0.1" db_port="3306" db_name="xmedia" jwt_secret="secret" video_dir="/home/xa0s/Downloads/Movies/" movies_sub_dir="/home/xa0s/Downloads/Movies/sub" frontend_dir="./dist" ./x-media

