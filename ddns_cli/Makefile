mod:
	GO111MODULE=on go mod init ddnscli
tidy:
	GO111MODULE=on go mod tidy
build:
	GO111MODULE=on go build -ldflags '-s -w' -o ddns_cli64.exe
	GO111MODULE=on GOARCH=386 go build -ldflags '-s -w' -o ddns_cli32.exe
build-linux:
	GO111MODULE=on GOARCH=386 GOOS=linux go build -ldflags '-s -w' -o ddns_cli_linux64
	GO111MODULE=on GOARCH=amd64 GOOS=linux go build -ldflags '-s -w' -o ddns_cli_linux32
all:
	make build build-linux
