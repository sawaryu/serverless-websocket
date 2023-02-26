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

### AWS setting

```bash
$ export AWS_ACCESS_KEY_ID=<your_access_key_id>
$ export AWS_SECRET_ACCESS_KEY=<your_secret_access_key>

# confirm
$ aws sts get-caller-identity --query Account --output text
```

### Build operation

```bash
# create
$ make create

# build
$ make build
```

# Test

```bash
$ npm install -g wscat

# execute connection and send any messages
$ wscat -c wss://{YOUR-API-ID}.execute-api.{YOUR-REGION}.amazonaws.com/{STAGE}
```