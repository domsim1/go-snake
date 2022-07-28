
.PHONY: snake
bin/snake:
	go build -o bin/snake -ldflags="-s -w" cmd/snake/main.go
