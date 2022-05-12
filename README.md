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

```bash
serverless create -u https://github.com/serverless/serverless-golang/ -p myservice
```

2. Compile function

```bash

$ cd myservice
$ GOOS=linux go build -o bin/main

# --or--

$ GOOS=linux go build -o ./bin/handleRequest
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
# invoke
$ sls invoke -f hello

# socket connection
$ npm install -g wscat

# execute connection and send any messages
$ wscat -c wss://{YOUR-API-ID}.execute-api.{YOUR-REGION}.amazonaws.com/{STAGE}
```

`cleanup and destroy aws environment`
```bash
# simply below
$ sls remove
```
# Syntax

`${opt:stage, 'dev'}`
cli option

`${self:service.name}`
reference self variable value

`${file(../myCustomFile.yml)}`
file reference

```
aws dynamodb put-item \
    --table-name connectionsTable \
    --item \
        '{"connectionId": {"S": "sampleconnectionid"}}' \
    --return-consumed-capacity TOTAL  

```