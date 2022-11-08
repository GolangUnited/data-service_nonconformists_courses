# syntax=docker/dockerfile:1

FROM golang:1.17-alpine as build
WORKDIR  /go/src/golang-united-courses/
COPY . ./
RUN go mod download
RUN go build -o ./server ./cmd/main.go

FROM alpine
COPY --from=build /go/src/golang-united-courses/server ./server

CMD ["./server"]