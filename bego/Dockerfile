# Stage 1: Build the Go binary
FROM golang:1.23 AS builder

WORKDIR /app

# Copy the source code
COPY . .

# Build the application
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ElkLogMonitor && chmod +x ElkLogMonitor

# Stage 2: Create a minimal runtime image
FROM alpine:latest

RUN apk add --no-cache tzdata

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/ElkLogMonitor .
# Copy the .env file
COPY --from=builder /app/.env .

# copy static file
COPY --from=builder /app/static/ /app/static/

# Set executable permissions
RUN chmod +x ElkLogMonitor

# Command to run the application
CMD ["./ElkLogMonitor", "run-monitor"]