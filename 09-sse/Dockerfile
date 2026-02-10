# --- Stage 1: The Build Environment ---
FROM golang:alpine AS builder

# Install git/ca-certificates if your private modules need them
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy the rest of the source code
COPY . .

# Build the binary
# CGO_ENABLED=0 creates a statically linked binary (ideal for scratch/alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o ssed ./cmd/main.go

# --- Stage 2: The Final Runtime Image ---
FROM alpine:latest

# Add a non-root user for security
RUN adduser -D agent
USER agent

WORKDIR /home/agent

# Copy ONLY the binary from the builder stage
COPY --from=builder /app/ssed .

# Set the entry point
ENTRYPOINT ["./ssed"]
