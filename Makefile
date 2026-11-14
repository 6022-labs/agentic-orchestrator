.PHONY: build test lint fmt hooks help

MODULES := agent-identity agent-swarm channels-and-integrations/slack-bridge gateway

build:
	@for module in $(MODULES); do (cd $$module && go build ./...) || exit 1; done

# test runs each module's unit tests; integration tests need Docker and stay in CI.
test:
	@for module in $(MODULES); do (cd $$module && go test $$(go list ./tests/... | grep -v '_integration_tests')) || exit 1; done

# lint requires golangci-lint v2 on PATH (go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest).
lint:
	@for module in $(MODULES); do (cd $$module && golangci-lint run ./...) || exit 1; done

# fmt runs golines via golangci-lint's built-in formatter; no standalone golines needed.
fmt:
	@for module in $(MODULES); do (cd $$module && golangci-lint fmt) || exit 1; done

# hooks installs pinned lefthook + gitleaks into GOPATH/bin and wires the git hooks.
hooks:
	go install github.com/evilmartians/lefthook/v2@v2.2.2
	go install github.com/zricethezav/gitleaks/v8@v8.30.1
	$$(go env GOPATH)/bin/lefthook install

help:
	@echo "build - go build ./... in every module"
	@echo "test  - unit tests in every module"
	@echo "lint  - golangci-lint run ./... in every module (needs golangci-lint v2)"
	@echo "fmt   - golangci-lint fmt in every module, 120-char cap"
	@echo "hooks - install git hooks (format + secret scan on commit, build/lint/test on push)"
