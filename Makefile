.PHONY: test
test:
	ginkgo run -v ./...


.PHONY: test-report
test-report:
	ginkgo run --cover -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out