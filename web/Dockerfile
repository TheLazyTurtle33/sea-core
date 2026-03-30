FROM golang:1.26.1-alpine AS builder
WORKDIR /app
# COPY web/go.mod web/go.sum* ./web/
COPY web/ ./web/
COPY shared/ ./shared/
RUN cd web && go mod download
RUN --mount=type=cache,target=/root/.cache/go-build \
    cd web && go build -o /app/bin/web ./cmd/web

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/web .
EXPOSE 8080
CMD ["./web"]