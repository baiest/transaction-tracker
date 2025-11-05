# ---------- Stage 1: Build ----------
FROM golang:1.23.4 AS builder
WORKDIR /app

# Download Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and compile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/api-bin ./cmd/api


# ---------- Stage 2: Final image ----------
FROM python:3.12-alpine

# Create non-root user and add necessary tools
RUN addgroup -S appgroup && adduser -S appuser -G appgroup \
 && apk --no-cache add ca-certificates tzdata curl netcat-openbsd tar \
 && pip install --upgrade pip setuptools wheel

WORKDIR /app

# Install migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz \
    | tar xz -C /usr/bin/ && chmod +x /usr/bin/migrate

# Create directories for files with correct permissions
RUN mkdir -p /files && chown -R appuser:appgroup /files
RUN mkdir -p /logs && chown -R appuser:appgroup /logs

# Copy binary, environment file and scripts
COPY --from=builder --chown=appuser:appgroup /app/api-bin /app/api-bin
COPY --chown=appuser:appgroup .env /app/.env
COPY --chown=appuser:appgroup pkg/databases/postgres/migrations /app/migrations
COPY --chown=appuser:appgroup ./cmd/api/entrypoint.sh /app/entrypoint.sh

# Copy Python dependencies and code
COPY --chown=appuser:appgroup pkg/document-extractor /app/pkg/document-extractor

# Install Python dependencies
RUN pip install --no-cache-dir -r /app/pkg/document-extractor/requirements.txt

# Execution permissions
RUN chmod +x /app/entrypoint.sh

USER appuser

ENV GIN_MODE=release
ENV PATH="/usr/bin:$PATH"
ENV FILES_DIR=/app/files

ENTRYPOINT ["/app/entrypoint.sh"]
