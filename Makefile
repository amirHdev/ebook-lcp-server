.PHONY: fmt lint test coverage vuln clean run acceptance e2e

COVERAGE_DIR := reports/coverage

fmt:
	gofumpt -w .

lint:
	golangci-lint run

test:
	go test ./... -cover -race -v

coverage:
	./scripts/coverage-report.sh $(COVERAGE_DIR)

vuln:
	govulncheck ./...

run:
	go run ./cmd/server

acceptance:
	sh scripts/acceptance-smoke.sh

e2e:
	sh scripts/e2e-readiness.sh

clean:
	go clean
