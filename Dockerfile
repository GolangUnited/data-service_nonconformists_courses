# syntax=docker/dockerfile:1

FROM golang:1.17-alpine as build
WORKDIR  /go/src/golang-united-courses/
COPY . ./
RUN go mod download
RUN go build -o ./server ./cmd/main.go

FROM alpine
COPY --from=build /go/src/golang-united-courses/server ./server
COPY --from=build /go/src/golang-united-courses/.env ./.env

EXPOSE 8080

CMD ["./server"]