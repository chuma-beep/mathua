FROM node:20-bookworm AS webbuilder

WORKDIR /app
COPY web/next-app/package.json web/next-app/package-lock.json ./web/next-app/
RUN cd web/next-app && npm ci

COPY web/next-app ./web/next-app
# Single-binary: same origin => relative fetches. Empty API base bakes fetch('/api/...') 
# which works on any Fly hostname (mathua.fly.dev, mathua-xxx.fly.dev). For split Vercel+Fly,
# set build arg NEXT_PUBLIC_API_URL=https://mathua.fly.dev
ARG NEXT_PUBLIC_API_URL=""
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
RUN cd web/next-app && npm run build

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
COPY --from=webbuilder /app/web/next-app/out /app/web/next-app/out

WORKDIR /app
ENV PORT=8080
ENV DATABASE_URL=/data/mathua.db

EXPOSE 8080

CMD ["/mathua", "-serve", "-port", "8080"]
