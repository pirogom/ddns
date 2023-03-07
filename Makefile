mod:
	GO111MODULE=on go mod init ddns
tidy:
	GO111MODULE=on go mod tidy
build:
	GO111MODULE=on go build -ldflags '-s -w' -o ddns_server64.exe
	GO111MODULE=on GOARCH=386 go build -ldflags '-s -w' -o ddns_server32.exe
build-linux:
	GO111MODULE=on GOARCH=386 GOOS=linux go build -ldflags '-s -w' -o ddns_server_linux64
	GO111MODULE=on GOARCH=amd64 GOOS=linux go build -ldflags '-s -w' -o ddns_server_linux32
all:
	make build build-linux
