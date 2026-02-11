.PHONY: test tinygo-build ci

test:
	GOWORK=off go test ./...

tinygo-build:
	@if command -v tinygo >/dev/null 2>&1; then \
		GOWORK=off tinygo build -target wasm -o /tmp/helianthus-tinyebus-placeholder.wasm ./firmware; \
		rm -f /tmp/helianthus-tinyebus-placeholder.wasm; \
		echo "tinygo build ok"; \
	else \
		echo "tinygo not installed; skipping tinygo build"; \
	fi

ci: test
