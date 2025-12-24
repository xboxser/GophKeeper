client:
	go run cmd/client/main.go

test:
	go test -v ./...

testPC:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep "total:"
	rm coverage.out