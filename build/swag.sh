#!/bin/sh
# the swagger lib we are using only works if the main is in the root dir,
# the deploy scripts only work if main is in the cmd dir, thus this script exists
mv ./cmd/hattrick/main.go .

swag init .

mv ./main.go ./cmd/hattrick/main.go