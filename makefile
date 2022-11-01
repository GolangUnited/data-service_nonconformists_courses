path:
	export PATH="$PATH:$(go env GOPATH)/bin"

protogen: 
	protoc -I courses --go_out=. --go-grpc_out=. courses/courses.proto

build:
	go build -v ./cmd/server