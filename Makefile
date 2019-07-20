VERSION := $(shell git describe --tags)

linux:
	GOOS=darwin GOARCH=386 go build -o ./target/refreturn -ldflags="-X main.version=${VERSION}" ./*.go
mac:
	GOOS=darwin GOARCH=amd64 go build -o ./target/refreturn -ldflags="-X main.version=${VERSION}" ./*.go
windows:
	GOOS=windows GOARCH=386 go build -o ./target/refreturn.exe -ldflags="-X main.version=${VERSION}" ./*.go
clean:
	rm -rf ./target
all: linux mac windows