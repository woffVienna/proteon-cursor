# -------- Build Stage --------
FROM golang:1.25.5-alpine AS builder

WORKDIR /repo

COPY go.work go.work.sum* ./
COPY contracts ./contracts
COPY libs ./libs
COPY services ./services

WORKDIR /repo/services/auth

# Download deps (cached)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Build static binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/service ./cmd/auth

# -------- Runtime Stage --------
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /out/service /app/service
COPY --from=builder /repo/services/auth/.build/generated/openapi.bundle.yml /app/.build/generated/openapi.bundle.yml

USER nonroot:nonroot
EXPOSE 8083
ENTRYPOINT ["/app/service"]
