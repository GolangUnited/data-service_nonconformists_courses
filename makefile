path:
	export PATH="$PATH:$(go env GOPATH)/bin"

protogen: 
	protoc -I api/proto --go_out=. --go-grpc_out=. api/proto/courses.proto

build:
	go build -v ./cmd/server