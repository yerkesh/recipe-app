FROM golang:1.17.8-alpine3.15 AS builder
COPY . /recipe-app/
WORKDIR /recipe-app/

RUN go mod download
RUN go build -o ./bin/app cmd/server/main.go

ENV POSTGRES_PASSWORD=pass \
    POSTGRES_USER=user \
    POSTGRES_DB=recipe_app_db

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /recipe-app/bin/app .
COPY --from=0 /recipe-app/resources/configs resources/configs/

EXPOSE 8090

CMD ["./app"]


