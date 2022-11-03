# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR  /go/src/golang-united-courses/

COPY . ./

RUN go mod download
RUN go mod verify
RUN go build -o ./server ./cmd/main.go

EXPOSE 8080

CMD ["/go/src/golang-united-courses/server"]