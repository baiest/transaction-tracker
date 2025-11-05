# ---------- Stage 1: Build (Optimized for Cache) ----------
FROM golang:1.23.4 AS builder
WORKDIR /app

# 1. Copy and download dependencies (cached unless go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy the rest of the source code (this step invalidates cache if app code changes)
#    This ensures the build only runs when necessary.
COPY . .

# 3. Build the binary. Rebuilds only if the previous step changes.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/tracker-bin ./cmd/tracker

# ----------------------------------------------------------------

# ---------- Stage 2: Final image ----------
FROM alpine:3.20

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install SSL certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary and .env file with correct permissions
COPY --from=builder --chown=appuser:appgroup /app/tracker-bin /app/tracker-bin
COPY --chown=appuser:appgroup .env /app/.env

# Copy service account JSON (as fallback)
# Copy credentials file from the project root.
COPY --chown=appuser:appgroup key.json /app/sa-key.json

# Switch to non-root user
USER appuser

# Set runtime env (useful for logging)
ENV GIN_MODE=release

# Run the binary
CMD ["/app/tracker-bin"]
