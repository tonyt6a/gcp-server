# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /src

# Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy all source (main.go + internal/)
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Runtime stage (minimal image)
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /src/server .

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/server"]