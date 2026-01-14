# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o boilerblade .

# Final stage
FROM alpine:latest

# Set dockerize version
ENV DOCKERIZE_VERSION=v0.7.0

# Install ca-certificates, openssl, and wget for dockerize
RUN apk --no-cache add ca-certificates tzdata openssl wget

# Install dockerize
RUN wget https://github.com/jwilder/dockerize/releases/download/${DOCKERIZE_VERSION}/dockerize-alpine-linux-amd64-${DOCKERIZE_VERSION}.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-${DOCKERIZE_VERSION}.tar.gz \
    && rm dockerize-alpine-linux-amd64-${DOCKERIZE_VERSION}.tar.gz

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/boilerblade .

# Copy environment example (for reference)
COPY --from=builder /app/env.example .

# Copy entrypoint script
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Expose port
EXPOSE 3000

# Use entrypoint script that dynamically waits for enabled services
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["./boilerblade"]
