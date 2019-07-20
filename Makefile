VERSION := $(git describe --tags)

linux:
	go build -o ./target/refreturn -ldflags="-X main.version=${VERSION}" ./main.go
mac:
	go build -o ./target/refreturn -ldflags="-X main.version=${VERSION}" ./main.go
windows:
	go build -o ./target/refreturn.exe -ldflags="-X main.version=${VERSION}" ./main.go
clean:
	rm -rf ./target
all: linux mac windows