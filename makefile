test: test-unit test-integration

test-unit:
	go test -tags unit -v ./...

test-integration:
	go test -tags integration -v ./...

cover:
	go test -tags unit -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out