wire-build:
	wire ./...

build:
	@if ! command -v wire &> /dev/null; then \
        echo "Installing 'wire'..."; \
        go install github.com/google/wire/cmd/wire@latest; \
    else \
        echo "'wire' is already installed."; \
    fi
	wire ./...
	GOGC=off go build -o binary