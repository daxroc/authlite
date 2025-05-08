# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o authlite ./cmd/server


FROM gcr.io/distroless/static:nonroot

WORKDIR /app
# Copy the statically built binary
COPY --from=builder /app/authlite .
# Use nonroot user provided by distroless
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/authlite"]
