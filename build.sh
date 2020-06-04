#!/bin/sh

rm aws-scraper.zip

set GOOS=linux
set GOARCH=amd64
go build .

build-lambda-zip -o aws-scraper.zip aws-scraper
rm aws-scraper
