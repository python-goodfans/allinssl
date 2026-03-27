# ============================================
# Stage 1: 构建前端 / Build frontend
# ============================================
FROM node:18-alpine AS frontend-builder

WORKDIR /frontend

# 启用 pnpm / Enable pnpm
RUN corepack enable && corepack prepare pnpm@10.14.0 --activate

# 先复制依赖清单以利用 Docker 缓存
# Copy dependency manifests first for better cache utilization
COPY frontend/package.json frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml frontend/turbo.json ./

# 复制工作区各包的 package.json（monorepo 结构）
# Copy workspace package.json files (monorepo structure)
COPY frontend/apps/allin-ssl/package.json ./apps/allin-ssl/
COPY frontend/packages/gulp-build-tools/package.json ./packages/gulp-build-tools/
COPY frontend/packages/utils/package.json ./packages/utils/
COPY frontend/packages/vue/hooks/package.json ./packages/vue/hooks/
COPY frontend/packages/vue/i18n/package.json ./packages/vue/i18n/
COPY frontend/packages/vue/naive-ui/package.json ./packages/vue/naive-ui/
COPY frontend/packages/vue/pinia/package.json ./packages/vue/pinia/
COPY frontend/packages/vue/router/package.json ./packages/vue/router/
COPY frontend/packages/vue/vite/package.json ./packages/vue/vite/
COPY frontend/env/eslint/package.json ./env/eslint/
COPY frontend/env/prettier/package.json ./env/prettier/
COPY frontend/env/stylelint/package.json ./env/stylelint/
COPY frontend/env/typescript/package.json ./env/typescript/
COPY frontend/plugin/vite-plugin-ftp-sync/package.json ./plugin/vite-plugin-ftp-sync/
COPY frontend/plugin/vite-plugin-i18n/package.json ./plugin/vite-plugin-i18n/
COPY frontend/plugin/vite-plugin-path-random/package.json ./plugin/vite-plugin-path-random/
COPY frontend/plugin/vite-plugin-turborepo-deploy/package.json ./plugin/vite-plugin-turborepo-deploy/

# 安装依赖（使用 frozen-lockfile 保证可复现，lockfile 不同步时回退）
# Install dependencies (frozen-lockfile for reproducibility, fallback if lockfile is out of sync)
RUN pnpm install --frozen-lockfile || pnpm install

# 复制全部前端源码
# Copy all frontend source code
COPY frontend/ .

# 构建前端应用
# Build frontend application
RUN pnpm build --filter allin-ssl

# ============================================
# Stage 2: 构建 Go 后端 / Build Go binary
# ============================================
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

# ============================================
# Stage 3: 最终运行镜像 / Final runtime image
# ============================================
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
