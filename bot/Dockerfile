FROM golang:1.26.1-alpine AS builder
WORKDIR /app
COPY bot/go.mod bot/go.sum* ./bot/
COPY shared/ ./shared/
RUN cd bot && go mod download
COPY bot/ ./bot/
RUN --mount=type=cache,target=/root/.cache/go-build \
    cd bot && go build -o /app/bin/bot ./cmd/bot

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/bot .
EXPOSE 8080
CMD ["./bot"]