install:
	go install

format:
	gofmt -s -w .

clean:
	find . -maxdepth 1 -not -name 'README.md' -name '*.md' -delete
	find . -maxdepth 1 -name '*.epub' -delete
	find . -maxdepth 1 -name '*.mobi' -delete
	find . -maxdepth 1 -name 'papeer-v*' -delete
