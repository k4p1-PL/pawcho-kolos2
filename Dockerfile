# ETAP 1: Budowanie (Builder) - używamy lekkiego obrazu golang na bazie alpine
FROM golang:alpine AS builder

# Instalacja certyfikatów SSL (potrzebne do zapytań HTTPS do API pogodowego w pustym obrazie scratch)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Optymalizacja cache: Kopiujemy i kompilujemy w jednej warstwie
COPY main.go .

# Kompilacja statyczna 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server main.go

# ETAP 2: Obraz docelowy
FROM scratch

# Standard OCI: Informacje o autorze 
LABEL org.opencontainers.image.authors="Kacper Sokół" \
      org.opencontainers.image.title="Aplikacja Pogodowa - Zadanie 1" \
      org.opencontainers.image.description="Minimalistyczny kontener w Go"

# Kopiowanie certyfikatów z etapu budowania 
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Kopiowanie tylko gotowego pliku binarnego
COPY --from=builder /app/server /server

# Deklaracja portu
EXPOSE 8080

# Healthcheck wywołujący naszą aplikację z flagą -health 
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/server", "-health"]

# Uruchomienie aplikacji
ENTRYPOINT ["/server"]