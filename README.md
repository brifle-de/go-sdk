# Brifle API Client

This go client allows to interact with the Brifle API via an sdk. The sdk is handling the renewing of the access token automatically. If the api only returns json strings it also converts those json string to go structs.

# Installation
To install the Brifle API Client, add the following line to your `go.mod` file:

```go
replace github.com/brifle-de/brifle-sdk => github.com/brifle-de/go-sdk v0.0.1

require github.com/brifle-de/brifle-sdk v0.0.1
```

Then run:

```bash
go get ./...
```


# Test

To run the test, we use a mock server. You either need mock data or valid credentials providing it via the .env.test file.

## Run Mock Server

```bash

docker run --rm \
sudo docker run -it --rm   -p 8080:8080   --name wiremock   -v $PWD/test/mock:/home/wiremock   wiremock/wiremock:3.13.1 --proxy-all="https://internaltest-api.brifle.de" --record-mappings --verbose
```

## .env.test
Create a `.env.test` file in the root directory of the project with the following content:

```dotenv
API_KEY=aaaa
API_SECRET=bbbb
TEST_RECEIVER_FIRST_NAME=Max
TEST_RECEIVER_LAST_NAME=Mustermann
TEST_RECEIVER_PLACE_OF_BIRTH=Berlin
TEST_RECEIVER_DATE_OF_BIRTH=1999-12-12
TEST_TENANT=567e44de-b6b6-4dac-cbce-c5515031f9ea
TEST_TENANT_ID=567e44de-b6b6-4dac-cbce-c5515031f9ea
TEST_USER_ID=65d626d2-9f59-4798-a26b-eec5d8439e48
TEST_DOC_ID_CERTIFICATE=53C9084932FA27B068424A5FCA81974873E54BC88AAB3B5CCB45C4E6E2C90BB1
TEST_DOC_ID_ACTIONS=14288A5F91EA1F2843A2EDEA542E26987F060314D1EF71EB6456CB88865DDA38
TEST_ACCOUNT_ID=2802510314782548
EXPORT_SIGNATURE_ID=1RTIJVJEAL9BlzdknOoPbgHlg5p0IKRizQaPACh0-O8=
ENDPOINT=http://localhost:8080
```

# Generate Api package

The api package is auto generated from the OpenAPI specification file `openapi.yaml` using the `openapi-generator-cli`.

```bash
make -f Makefile generate
```

