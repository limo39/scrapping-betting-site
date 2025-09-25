# Build stage
FROM golang:1.21-alpine AS builder

# Install dependencies for Chrome
RUN apk add --no-cache \
    chromium \
    ca-certificates \
    tzdata

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install Chrome and dependencies
RUN apk add --no-cache \
    chromium \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*

# Create app user
RUN addgroup -g 1001 -S app && \
    adduser -S app -u 1001

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web
COPY --from=builder /app/.env.example ./.env

# Change ownership
RUN chown -R app:app /app

# Switch to app user
USER app

# Expose port
EXPOSE 8080

# Set Chrome path for headless mode
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/bin/chromium-browser

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Run the application
CMD ["./main"]