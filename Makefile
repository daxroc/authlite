.PHONY: test cover html

test:
	go test ./...

cover:
	go test -coverprofile=coverage.out ./...

html: cover
	go tool cover -html=coverage.out