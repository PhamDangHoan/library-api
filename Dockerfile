# ==========================================
# Stage 1: Build
# ==========================================
FROM golang:1.27.1-alpine AS builder

WORKDIR /app

# Install CA certificates
RUN apk add --no-cache ca-certificates

# Copy dependency files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o library-api ./cmd/api


# ==========================================
# Stage 2: Runtime
# ==========================================
FROM alpine:latest

WORKDIR /app

# Install CA certificates
RUN apk add --no-cache ca-certificates

# Copy compiled binary
COPY --from=builder /app/library-api .

# Application port
EXPOSE 8080

# Start application
CMD ["./library-api"]