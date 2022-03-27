.PHONY:
.SILENT:

build:
	go build -o ./.bin/app cmd/server/main.go

run:
	./.bin/app

build-container:
	docker build -t recipe-app:v0.1 .

start-container:
	docker run --name recipe-app -p 8090:8090 --env-file .env recipe-app:v0.1