.PHONY: lint fmt help

MODULES := agent-identity agent-swarm channels-and-integrations/slack-bridge gateway

# lint requires golangci-lint v2 on PATH (go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest).
lint:
	@for module in $(MODULES); do (cd $$module && golangci-lint run ./...) || exit 1; done

# fmt runs golines via golangci-lint's built-in formatter; no standalone golines needed.
fmt:
	@for module in $(MODULES); do (cd $$module && golangci-lint fmt) || exit 1; done

help:
	@echo "lint - golangci-lint run ./... in every module (needs golangci-lint v2)"
	@echo "fmt  - golangci-lint fmt in every module, 120-char cap"
