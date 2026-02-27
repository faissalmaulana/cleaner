.PHONY: test-e2e clean-e2e tag

COVER_DIR=.coverdata

tag:
ifndef VERSION
	$(error Usage: make tag VERSION=0.0.0)
endif
	@echo "Updating VERSION in cmd/root.go to $(VERSION)"
	@sed -i 's/var Version = ".*"/var Version = "$(VERSION)"/' cmd/root.go
	@git add cmd/root.go
	@git commit -m "Bump version to $(VERSION)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "Done! Run 'git push && git push --tags' to publish"


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
