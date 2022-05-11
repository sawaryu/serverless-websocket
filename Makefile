.PHONY: build
build:
	GOOS=linux go build -o bin/handleRequest

.PHONY: deploy
deploy:
	sls deploy --verbose

.PHONY: remove
remove:
	sls remove