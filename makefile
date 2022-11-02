all:
	@echo "test"
	
protogen: 	
	protoc -I api/proto --go_out=. --go-grpc_out=. api/proto/courses.proto

build:
	go build -v -o ./server ./cmd/main.go 

run:
	go run ./cmd/main.go

db-run: db-remove
	docker run --name postgres --network my-network -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=postgres -p 5432:5432 -d postgres:alpine || true

db-start:
	docker start postgres || true

db-stop:
	docker stop postgres || true

db-remove: db-stop
	docker rm postgres || true

app-start: app-stop
	docker build -t golang-united-courses .
	docker run --name golang-united-courses --network my-network -p 8080:8080 -d --rm -e COURSES_DB_HOST=postgres -e COURSES_DB_PORT -e COURSES_DB_USER -e COURSES_DB_PASSWORD -e COURSES_DB_NAME golang-united-courses

app-stop:
	docker stop golang-united-courses || true
	docker rm golang-united-courses || true
	docker rmi golang-united-courses || true	

app-shell:
	docker exec -it golang-united-courses /bin/sh