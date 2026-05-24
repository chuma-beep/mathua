FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o mathua ./cmd/mathua

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/mathua .
COPY --from=builder /app/web/next-app/out ./web/next-app/out
COPY --from=builder /app/data ./data

ENV DATABASE_URL=/data/mathua.db

EXPOSE 7860

CMD ["./mathua", "--serve", "--port", "7860"]