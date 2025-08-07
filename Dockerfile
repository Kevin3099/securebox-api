# -------- Stage 1: Build the Go binary --------
FROM golang:1.24-alpine AS builder

# Produce a small, statically-linked Linux binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

# Copy module definition (no go.sum needed at this point)
COPY go.mod ./
RUN go mod download

# Copy your source code
COPY main.go .

# Build the Go binary
RUN go build -o securebox .

# -------- Stage 2: Create minimal runtime image --------
FROM alpine:latest

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/securebox .

# Expose port the app listens on
EXPOSE 8080

# Run the binary
CMD ["./securebox"]
