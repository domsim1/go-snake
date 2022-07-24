
.PHONY: snake
snake:
	go build -o snake -ldflags="-s -w" main.go
