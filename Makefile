APPNAME=hv

.PHONY: build

build:
	GOOS=linux GOARCH=amd64 go build -ldflags=" -s -w" -v -o $(APPNAME)

