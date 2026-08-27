# ==========================================
# STAGE 1: Build Frontend (Vue 3 + Vite)
# ==========================================
FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ ./
RUN touch .env && echo "VITE_API_URL=$VITE_API_URL" >> .env && echo "VITE_BACKEND_URL=$VITE_BACKEND_URL" >> .env
RUN npm run build

# ==========================================
# STAGE 2: Build Backend (Go)
# ==========================================
FROM golang:alpine AS backend-builder

WORKDIR /app/backend

RUN apk add --no-cache ca-certificates git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/backend/server ./cmd/server

# ==========================================
# STAGE 3: Final Production Image (Alpine)
# ==========================================
FROM alpine:3.21

WORKDIR /app

# Install SSL certificates and timezone data
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy compiled Go binary and frontend static build
COPY --from=backend-builder /app/backend/server /app/server
COPY --from=frontend-builder /app/frontend/dist /app/dist

USER appuser

EXPOSE 8080

CMD ["/app/server"]
