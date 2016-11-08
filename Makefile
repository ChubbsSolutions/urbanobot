all: build
build:
	go build
install:
	go install
buildall:
	    env GOOS=darwin GOARCH=amd64 go build -o urbanobot_darwin_amd64
	    env GOOS=windows GOARCH=amd64 go build -o urbanobot_windows_amd64.exe
	    env GOOS=linux GOARCH=amd64 go build  -o urbanobot_linux_amd64
clean:
	rm urbanobot
