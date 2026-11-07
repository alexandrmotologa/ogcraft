# Stage 1: Build Web Visual Designer
FROM node:22-alpine AS web-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci --prefer-offline || npm install
COPY web/ ./
RUN npm run build

# Stage 2: Build Static Go Binary
FROM golang:1.24-alpine AS go-builder
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata git

COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .
# Copy compiled web assets into embed location
COPY --from=web-builder /app/web/dist ./web/dist

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o /app/ogcraft ./cmd/ogcraft

# Stage 3: Minimal Runtime Container
FROM alpine:3.21
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    mkdir -p /app/data/cache /app/data/templates

COPY --from=go-builder /app/ogcraft /app/ogcraft

EXPOSE 8080

VOLUME ["/app/data"]

ENTRYPOINT ["/app/ogcraft"]
CMD ["serve", "--port", "8080", "--cache-dir", "/app/data/cache", "--templates-dir", "/app/data/templates"]
