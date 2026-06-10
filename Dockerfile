FROM golang:alpine AS builder

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server main.go

FROM scratch

LABEL org.opencontainers.image.authors="Kacper Sokół" \
      org.opencontainers.image.title="Aplikacja Pogodowa - Zadanie 1" \
      org.opencontainers.image.description="Minimalistyczny kontener w Go"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/server /server

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/server", "-health"]

ENTRYPOINT ["/server"]
