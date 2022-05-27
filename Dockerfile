FROM golang:alpine

RUN apk update && apk add --no-cache git
RUN apk add openssh

WORKDIR /app

ARG SSH_RSA_KEY

ENV SSH_RSA_KEY=$SSH_RSA_KEY
ENV GIN_MODE=release

RUN mkdir /root/.ssh
RUN echo "$SSH_RSA_KEY" >> /root/.ssh/id_rsa
RUN chmod 400 /root/.ssh/id_rsa
RUN ssh-keyscan github.com > /root/.ssh/known_hosts
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"
RUN go env -w GOPRIVATE=github.com/toel-app/

COPY . .
RUN go mod tidy
RUN GOGC=off go build -o binary

ENTRYPOINT ["/app/binary"]
