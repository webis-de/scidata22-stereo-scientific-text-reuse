clean:
	if [ -a coverage.cov ]; then rm coverage.cov; fi;
	go clean -testcache
	go mod tidy

dependencies: go.mod
	go mod download && go mod verify && go mod tidy

test: clean dependencies
	go test -covermode=count -coverprofile "coverage.cov" ./... && go tool cover -func=coverage.cov && go tool cover -html=coverage.cov
	rm coverage.cov

build: clean dependencies
	CGO_ENABLED=0 go build -o ../../jobs/alignment main.go
