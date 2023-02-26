# Serverless Websocket

serverless simple websocket with `lambda` and `dynamodb`

<br>

# Setup

### AWS setting

```bash
$ export AWS_ACCESS_KEY_ID=<your_access_key_id>
$ export AWS_SECRET_ACCESS_KEY=<your_secret_access_key>

# confirm
$ aws sts get-caller-identity --query Account --output text

# create .env
$ echo "AWS_ACCOUNT_ID=<your_account_id>" > .env
```

<br>

### Deploy

```bash
# build
$ make build

# deploy
$ make deploy

# remove
$ make remove
```

<br>

# Test

test with `wscat` module

```bash
# install ws cat
$ npm install -g wscat

# execute connection and send any messages
$ wscat -c wss://<your_apigateway_id>.execute-api.ap-northeast-1.amazonaws.com/dev
```
