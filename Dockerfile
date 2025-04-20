# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Create uploads directory with proper permissions
RUN mkdir -p ./uploads && chmod 755 ./uploads

# Install ca-certificates for HTTPS support
RUN apk --no-cache add ca-certificates tzdata

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o file-sharing .

# Stage 2: Create the runtime image from scratch (minimal image with no OS)
FROM scratch

WORKDIR /app

# Copy CA certificates from the builder stage for HTTPS support
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data from the builder stage
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary from the builder stage
COPY --from=builder /app/file-sharing .

# Copy necessary directories for the application
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/uploads ./uploads

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["/app/file-sharing"]
