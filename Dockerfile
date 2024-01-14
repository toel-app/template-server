FROM golang:alpine

RUN mkdir /app
ADD . /app
WORKDIR /app

ENV GIN_MODE=release

RUN go mod tidy
RUN if ! command -v wire &> /dev/null; then \
        echo "Installing 'wire'..."; \
        go install github.com/google/wire/cmd/wire@latest; \
    else \
        echo "'wire' is already installed."; \
    fi
RUN wire ./...
RUN GOGC=off go build -o binary

CMD ["/app/binary"]
