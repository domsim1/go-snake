
.PHONY: bin/snake
bin/snake:
	go build -o bin/snake -ldflags="-s -w" cmd/snake/main.go

.PHONY: bundle
bundle: bin/snake
	mkdir Snake.app
	mkdir Snake.app/Contents
	mkdir Snake.app/Contents/MacOS
	mkdir Snake.app/Contents/Resources
	cp bin/snake Snake.app/Contents/MacOS
	mkdir Snake.app/Contents/MacOS/resources
	cp resources/eat.ogg Snake.app/Contents/MacOS/resources

