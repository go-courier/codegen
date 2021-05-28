test: download
	go test -v -race ./...

cover: download
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

download:
	go mod download -x