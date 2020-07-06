GO_TEST ?= go test

.PHONY: test
test:
	$(GO_TEST) -v -race ./...

build/%:
	@go build -o bin/$* cmd/$*/*

.PHONY: clean
clean:
	@rm -rf bin
