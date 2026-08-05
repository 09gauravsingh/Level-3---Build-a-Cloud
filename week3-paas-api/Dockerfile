# Stage 1: Compile the Go API into a Linux executable.
FROM golang:1.26-alpine AS builder

WORKDIR /source

# Download dependencies first so Docker can cache this layer.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and build the API.
COPY . .

# CGO is disabled so the binary can run in a small runtime image.
RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /output/week3-paas-api \
    ./cmd/api

# Stage 2: Create the small runtime image.
FROM alpine:3.22

# CA certificates are needed for HTTPS communication with Kubernetes.
RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app

# Copy only the compiled executable, not the Go source code.
COPY --from=builder \
    /output/week3-paas-api \
    /app/week3-paas-api

# Run the API without root privileges.
USER app

EXPOSE 8080

ENTRYPOINT ["/app/week3-paas-api"]
