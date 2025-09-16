test:
	go test -race ./...

test-cookbook:
	@echo "Running tests for cookbook examples..."
	@for dir in cookbook/*/; do \
		if [ -f "$$dir/server_test.go" ]; then \
			echo "Testing $$dir..."; \
			(cd "$$dir" && go test -v ./...) || exit 1; \
		fi; \
	done

test-verbose:
	go test -race -v ./...

test-cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

benchmark:
	@echo "Running benchmarks for cookbook examples..."
	@for dir in cookbook/*/; do \
		if [ -f "$$dir/server_test.go" ]; then \
			echo "Benchmarking $$dir..."; \
			(cd "$$dir" && go test -bench=. -benchmem) || true; \
		fi; \
	done

serve:
	docker run --rm -it --name echo-docs -v ${CURDIR}/website:/home/app -w /home/app -p 3000:3000 -u node node:lts /bin/bash -c "npm install && npm start -- --host=0.0.0.0"

deps:
	@echo "Installing test dependencies for cookbook examples..."
	@for dir in cookbook/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "Installing deps for $$dir..."; \
			(cd "$$dir" && go mod tidy && go mod download) || exit 1; \
		fi; \
	done

clean:
	@echo "Cleaning test artifacts..."
	find . -name "coverage.out" -delete
	find . -name "coverage.html" -delete
	find . -name "*.prof" -delete

.PHONY: test test-cookbook test-verbose test-cover benchmark serve deps clean
