# --- Stage 1: build the Vue frontend ---
FROM node:22-alpine AS frontend
WORKDIR /build
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ .
RUN npm run build

# --- Stage 2: build the Go backend ---
FROM golang:1.26-alpine AS backend
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

# --- Stage 3: runtime ---
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=backend /out/server /app/server
COPY --from=frontend /build/dist /app/web/dist
EXPOSE 8080
ENTRYPOINT ["/app/server"]
