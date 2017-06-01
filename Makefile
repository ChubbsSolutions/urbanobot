all: build
build:
	go build
install:
	go install
buildall:
	    env GOARM=7 GOOS=linux GOARCH=arm go build -o urbanobot_linux_armv7
	    env GOOS=darwin GOARCH=amd64 go build -o urbanobot_darwin_amd64
	    env GOOS=windows GOARCH=amd64 go build -o urbanobot_windows_amd64.exe
	    env GOOS=linux GOARCH=amd64 go build  -o urbanobot_linux_amd64
	    env GOOS=linux GOARCH=arm64 go build  -o urbanobot_linux_arm64
clean:
	rm urbanobot
