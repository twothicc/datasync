
all: format lint compile

format:
	gofmt -s -w $$(find . -type f -name '*.go'| grep -v "/vendor/")

lint:
	golangci-lint run

compile:
	cd app; go build -o ../build/datasync

install:
	cd app; go install