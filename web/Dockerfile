# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o web ./cmd/web

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/web .
EXPOSE 8080
CMD ["./web"]
