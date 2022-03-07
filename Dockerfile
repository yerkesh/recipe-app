FROM golang:1.17.8-alpine3.15 AS builder
COPY . /recipe-app/
WORKDIR /recipe-app/

RUN go mod download
RUN go build -o /bin/app cmd/server/main.go


FROM alpine:latest

WORKDIR /root/

COPY --from=0 /recipe-app/bin/app .
COPY --from=0 /recipe-app/configs configs/

EXPOSE 8090

CMD ["./app"]


