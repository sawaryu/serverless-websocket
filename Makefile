.PHONY: cretae
create:
	sls create -u https://github.com/serverless/serverless-golang/ -p serverless-websocket

.PHONY: build
build:
	GOOS=linux go build -o bin/handleRequest

.PHONY: deploy
deploy:
	sls deploy --verbose

.PHONY: invoke
invoke:
	sls invoke -f hello --log

.PHONY: remove
remove:
	sls remove