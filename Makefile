generate:
	go generate ./...
dev:
	go run main.go
start_mock_server:
	docker run -it --rm   -p 8080:8080   --name wiremock   -v ./test/mock:/home/wiremock   wiremock/wiremock:3.13.1 --proxy-all="https://internaltest-api.brifle.de" --record-mappings --verbose