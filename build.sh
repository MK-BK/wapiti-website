#!/bin/bash
go build -o server main.go

docker build -t wapiti-server:$1 .

docker tag  wapiti-server:$1  wardknight/wapiti-server:$1

docker push wardknight/wapiti-server:$1

rm server

