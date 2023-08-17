.PHONY: run build test

test:
	go test -race
run:
	go run cmd/vote/main.go