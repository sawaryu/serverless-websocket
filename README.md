# UPDATE

Starting from [version 1.26](https://github.com/serverless/serverless/releases/tag/v1.26.0) Serverless Framework includes two Golang templates:

* `aws-go` - basic template with two functions
* `aws-go-dep` - **recommended** template using [`dep`](https://github.com/golang/dep) package manager

You can use them with `create` command:

```
serverless create -t aws-go-dep
```

Original README below.

---

# Serverless Template for Golang

This repository contains template for creating serverless services written in Golang.

## Quick Start

1. Create a new service based on this template

```
serverless create -u https://github.com/serverless/serverless-golang/ -p myservice
```

2. Compile function

```
cd myservice
GOOS=linux go build -o bin/main
```

3. Deploy!

```
serverless deploy
```

# Setup

`operation to build`
```bash
# start
$ sls create -u https://github.com/serverless/serverless-golang/ -p micro-socket

# init go mod and get github.com/aws/aws-lambda-go
$ go mod init github.com/sawaryu/micro-socket
$ go mod tidy

# build
$ GOOS=linux go build -o bin/main
```

`serverless.yml`
```yml
frameworkVersion: ^3.15.2

service: micro-socket

provider:
  stage: dev
  region: ap-northeast-1
  name: aws
  runtime: go1.x
  environment:
    TZ: Asia/Tokyo

package:
  exclude:
    - ./**
  include:
    - ./bin/**

functions:
  hello:
    handler: bin/main
    name: micro-socket-example-function
```

`deploy`
```bash
# deploy
$ sls deploy --verbose
```

`test`
```bash
$ sls invoke -f hello
```

`cleanup and destroy aws environment`
```bash
# simply below
$ sls remove
```
