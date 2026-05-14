# Stage 1 — builder
# Start from the official Go image — has everything needed to compile Go code
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first
# Docker caches layers — if these files don't change, it won't re-download dependencies
COPY go.mod go.sum ./
RUN GOTOOLCHAIN=auto go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
# CGO_ENABLED=0 — disable C bindings, makes the binary fully static (needed for alpine)
# -o app — output file named "app"
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api/main.go

# ──────────────────────────────────────────
# Stage 2 — runner
# Start from a tiny image — no Go, no compiler, just a bare Linux
FROM alpine:latest

WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/app .

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./app"]