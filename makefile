build:
	go mod tidy && go build -o psychic-octo-stock
start:
	./psychic-octo-stock
test:
	go test -v ./... -race -covermode=atomic