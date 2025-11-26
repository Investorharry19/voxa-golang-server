# ============================
# Stage 1: Build the Go binary
# ============================
FROM golang:1.22-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache gcc g++ make

WORKDIR /app

# Cache mod downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the server
RUN go build -o server .



# =============================================
# Stage 2: Runtime Image (small + FFmpeg ready)
# =============================================
FROM alpine:latest

# Install FFmpeg and required dependencies
RUN apk add --no-cache ffmpeg ca-certificates tzdata

# Copy compiled Go binary from builder stage
COPY --from=builder /app/server /server

# Copy environment file (optional — remove if using Render env vars)
# COPY --from=builder /app/.env /.env

# Expose the port your app runs on
EXPOSE 3000

# Run the app
CMD ["/server"]
