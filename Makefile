PROJ_PATH = $(GOPATH)/src/github.com/aaronflower/shippy-service-user
build:
	protoc -I. --go_out=plugins=micro:$(PROJ_PATH) proto/auth/auth.proto
	GOOS=linux GOARCH=amd64 go build -o service.user
	docker build --rm -t service.user .

run:
	docker run --net="host" \
		-p 50051  \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns \
		-e DB_HOST=localhost \
		-e DB_PASS=password \
		-e DB_USER=postgres \
		service.user
	
clean:
	go clean
	rm service.user
