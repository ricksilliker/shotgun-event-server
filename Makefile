.PHONY: build-client
build-client:
	docker run --rm -v "$PWD/client":/usr/src/event-client -w /usr/src/event-client golang:1.14 go build -v

.PHONY: build-server
build-server:
	docker run --rm -v "$PWD/server":/usr/src/event-server -w /usr/src/event-server golang:1.14 go build -v
