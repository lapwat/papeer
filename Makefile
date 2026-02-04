format:
	gofmt -s -w .

test:
	go test ./...

install:
	go install

update:
	go get -u ./...
	go mod tidy

clean:
	find . -maxdepth 1 -name '*.md' -not -name 'README.md' -not -name 'LICENSE.md' -delete
	find . -maxdepth 1 -name '*.epub' -delete
	find . -maxdepth 1 -name '*.mobi' -delete
