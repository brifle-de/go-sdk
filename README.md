# Startup a Mock Server

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
TEST_TENANT=963e44de-c6b1-2eac-cbce-c5515031f9ea
ENDPOINT=http://localhost:8080
```

# Generate Api package

The api package is auto generated from the OpenAPI specification file `openapi.yaml` using the `openapi-generator-cli`.

```bash
make -f Makefile generate
```

