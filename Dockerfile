# Start from official Go image for building
FROM golang:alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install git and build dependencies
RUN apk update && apk add --no-cache git build-base

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app with optimizations
# -w -s flags strip debugging information to reduce binary size
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o wago ./src/main.go

# Start a new stage from a minimal alpine image
FROM alpine:latest

# Add ca-certificates for HTTPS (needed for OpenAI/OpenRouter APIs) and tzdata for timezones
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/wago .

# Copy static assets and templates
COPY --from=builder /app/src/static ./src/static

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./wago"]
