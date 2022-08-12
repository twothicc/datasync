
all: format lint compile

fieldAlignment:
	fieldalignment -fix github.com/twothicc/datasync

format:
	gofmt -s -w $$(find . -type f -name '*.go'| grep -v "/vendor/")

lint:
	golangci-lint run

compile:
	cd app; go build -o ../build/datasync

clearLog:
	> server.log

start:
	./build/datasync