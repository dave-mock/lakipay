# Build stage
FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Copy .env file into the build context and container
COPY .env .env

RUN go build -o main ./src/cmd/v1/main.go

# Run stage with Ubuntu as the base image
FROM ubuntu:22.04

RUN useradd -m appuser

WORKDIR /home/appuser/

# Copy application and environment files
COPY --from=builder /app/main .
COPY --from=builder /app/.env /home/appuser/.env

# Ensure the appuser has read permissions to the .env file
RUN chmod 644 /home/appuser/.env && chown appuser:appuser /home/appuser/.env

EXPOSE 3000

# Run the app as a non-root user
USER appuser
CMD ["./main"]
