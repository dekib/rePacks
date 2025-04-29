FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy only dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o rePacks ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Create necessary directories
RUN mkdir -p ui config

# Copy only what's needed from builder
COPY --from=builder /app/rePacks .
COPY --from=builder /app/ui/ ./ui/
#COPY --from=builder /app/config/ ./config/  # Only if you add config files

# Copy static files directly
COPY ui/ ./ui/

EXPOSE 8086

CMD ["./rePacks"]