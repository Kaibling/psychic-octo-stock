build:
	go build -o psychic-octo-stock
build-start:
	go build -o psychic-octo-stock && ./psychic-octo-stock
test:
	go test -v ./...
coverage:
	go test -v ./... -race -covermode=atomic -coverprofile=coverage.out