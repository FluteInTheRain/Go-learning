# Multi-stage build for Go app with PostgreSQL support

# Stage 1: builder
FROM golang:1.27 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server main.go

# Stage 2: runtime
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /server /app/server

EXPOSE 8080

CMD ["/app/server"]
