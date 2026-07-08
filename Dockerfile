# ── Stage 1: builder ─────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Descargar dependencias en capa separada para aprovechar cache de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente y compilar el binario
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/garra-api ./cmd/main.go

# ── Stage 2: runner ───────────────────────────────────────────────────────────
FROM alpine:3.19 AS runner

# Certificados para HTTPS y zona horaria
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Solo copiamos el binario compilado — imagen mínima sin toolchain de Go
COPY --from=builder /app/garra-api .

EXPOSE 8080

ENTRYPOINT ["/app/garra-api"]
