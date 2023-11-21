NAME=go-fileshare

windows:
	env GOOS=windows GOARCH=amd64 go build -o ${NAME}.exe

linux:
	env GOOS=linux GOARCH=amd64 go build -o ${NAME}-linux

mac:
	env GOOS=darwin GOARCH=amd64 go build -o ${NAME}-darwin

all: windows linux mac