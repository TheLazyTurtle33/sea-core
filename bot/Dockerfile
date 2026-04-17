FROM golang:1.26.1-alpine AS builder
WORKDIR /app
COPY bot/go.mod bot/go.sum* ./bot/
COPY shared/ ./shared/
RUN cd bot && go mod download
COPY bot/ ./bot/
RUN --mount=type=cache,target=/root/.cache/go-build \
cd bot && go build -gcflags="all=-N -l" -o /app/bin/bot ./cmd/bot

# RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest
    
    
FROM alpine:latest
RUN apk add --no-cache openssh-client ffmpeg
WORKDIR /app
COPY --from=builder /app/bin/bot .
# COPY --from=builder /go/bin/dlv /usr/local/bin/dlv
EXPOSE 8080
CMD ["./bot"]
# CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./bot"]