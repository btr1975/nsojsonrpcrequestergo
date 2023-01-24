# GO-Lang Helpers
# Version: 2023.1.14.001
# Author: Ben Trachtenberg
#
# Description:
#     This is a Makefile to help in building and testing GO projects
#
# Notes:
#     To see all OS/Architecture's you can build use "go tool dist list"
#
#     For compiling for ARM based architecture's you may require the GOARM variable
#     see docs for more info
#
.PHONY: all test coverage coverage-html

all:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test:
	go test ./...

coverage:
	go test ./... -cover

coverage-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
