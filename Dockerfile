# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build scheduler and worker binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler ./cmd/scheduler
RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker

# Runtime image
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/scheduler /app/worker ./

ENV PORT=8080
# Default entrypoint: scheduler (worker overrides this in k8s)
CMD ["/app/scheduler"]
