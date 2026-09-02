FROM golang:1.27-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc libc6-dev libsqlite3-dev && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o /mathua ./cmd/mathua

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    libsqlite3-0 ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /mathua /mathua
COPY data/concepts/ /app/data/concepts/
COPY data/lessons/ /app/data/lessons/
COPY data/courses.json /app/data/courses.json

WORKDIR /app
ENV PORT=8080
ENV DATABASE_URL=/data/mathua.db

EXPOSE 8080

CMD ["/mathua", "-serve", "-port", "8080"]
