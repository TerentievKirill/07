FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o parcel-tracker .

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y sqlite3 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/parcel-tracker .
COPY tracker.db .

CMD ["./parcel-tracker"]
