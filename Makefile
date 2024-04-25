.PHONY: linter
linter:
	golangci-lint run -v --color=always $$GO_LINT_FLAGS $$GO_PACKAGES --timeout 4m