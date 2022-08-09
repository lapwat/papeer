format:
	gofmt -s -w .

test:
	go test ./...

install:
	go install

clean:
	find . -maxdepth 1 -name '*.md' -not -name 'README.md' -delete
	find . -maxdepth 1 -name '*.epub' -delete
	find . -maxdepth 1 -name '*.mobi' -delete
