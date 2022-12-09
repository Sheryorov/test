
FROM golang:1.18.1 AS builder

WORKDIR /usr/app

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd

# multistage build to copy only binary and config
FROM scratch

ENV PORT=:3000
COPY --from=builder /usr/app/main /

EXPOSE 3000

ENTRYPOINT ["/main"]