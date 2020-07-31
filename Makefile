
.PHONY: test
test: lint
	go test -race -short -cover -coverprofile coverage.txt ./... \
		&& go tool cover -func=coverage.txt

.PHONY: lint
lint:
	go vet ./...

.PHONY: clean
clean:
	rm -rf coverage*
