build-dev:
	docker build . --target development --tag openslides-search-dev

run-tests:
	docker build . --target testing --tag openslides-search-test
	docker run openslides-search-test

all: gofmt gotest golinter

gotest:
	go test ./...

golinter:
	golint -set_exit_status ./...

gofmt:
	gofmt -l -s -w .
