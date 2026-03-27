# Stage 1: Build frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /frontend

RUN corepack enable && corepack prepare pnpm@10.14.0 --activate

COPY frontend/ .

RUN pnpm install

RUN pnpm build --filter allin-ssl

# Stage 2: Build Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Replace stale pre-committed frontend build with fresh build
COPY --from=frontend-builder /frontend/apps/allin-ssl/dist/ ./static/build/

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o allinssl ./cmd/main.go

# Final stage
FROM frolvlad/alpine-glibc

WORKDIR /www/allinssl/

# Install runtime dependencies
RUN apk add --no-cache tzdata

# Copy binary and script from builder
COPY --from=builder /build/allinssl /www/allinssl/allinssl
COPY --from=builder /build/script/allinssl.sh /www/allinssl/allinssl.sh

RUN chmod +x /www/allinssl/allinssl.sh

ENV TZ=Asia/Shanghai
RUN cat > /entrypoint.sh <<'EOF'
#!/bin/sh
if [ ! -f /www/allinssl/data/.initialized ]; then
    echo ${ALLINSSL_USER:-allinssl} | /www/allinssl/allinssl 5
    echo ${ALLINSSL_URL:-/} | /www/allinssl/allinssl 4
    echo ${ALLINSSL_PWD:-allinssldocker} | /www/allinssl/allinssl 6
    echo 8888 | /www/allinssl/allinssl 7
    touch /www/allinssl/data/.initialized
fi
/www/allinssl/allinssl 2
exec /www/allinssl/allinssl start
EOF
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 8888
