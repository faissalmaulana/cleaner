.PHONY: test-e2e clean-e2e

COVER_DIR=.coverdata

test-e2e:
	@echo "Running E2E tests..."
	@rm -rf $(COVER_DIR)
	@mkdir -p $(COVER_DIR)
	go test ./tests/e2e -v
	@echo "Coverage report:"
	@go tool covdata percent -i=$(COVER_DIR)

clean-e2e:
	@echo "Cleaning build artifacts..."
	@rm -rf bin $(COVER_DIR)
