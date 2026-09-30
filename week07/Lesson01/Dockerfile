# Dockerfile - production-ready multi-stage build

# ---------- Stage 1: builder ----------
FROM golang:1.27 AS builder

WORKDIR /app

# Copy go.mod (and go.sum) first to leverage Docker layer caching:
# dependencies are re-downloaded only when these files change.
COPY go.mod ./
RUN go mod download

# Now copy the rest of the source code.
COPY . .

# Build a fully static binary so it runs on a minimal base image.
# CGO_ENABLED=0 removes the C dependency; -ldflags "-s -w" strips debug info.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server main.go

# ---------- Stage 2: runtime ----------
FROM alpine:3.20

WORKDIR /app

# Copy only the compiled binary from the builder stage.
COPY --from=builder /server /app/server

EXPOSE 8080
CMD ["/app/server"]
