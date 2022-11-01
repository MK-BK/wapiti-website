#!/bin/bash
go build -o server main.go

docker build -t wapiti-server:v1 .

rm -f server
