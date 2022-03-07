.PHONY:
.SILENT:

build:
	go build -o ./.bin/app cmd/server/main.go

run:
	./.bin/app

build-container:
	docker build -t recipe-app:v0.1 .