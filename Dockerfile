FROM golang:1.17-alpine

COPY . /go/src/goland-united-courses/
WORKDIR  /go/src/goland-united-courses/

RUN go mod download
RUN go mod verify
RUN go build -o ./server ./cmd/main.go

EXPOSE 8080

RUN ls -la

CMD ["/go/src/golang-united-courses/server"]