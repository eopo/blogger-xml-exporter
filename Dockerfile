# syntax=docker/dockerfile:1

# --- Build-Stage ---------------------------------------------------------
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/blogger-xml-exporter .

# --- Runtime-Stage --------------------------------------------------------
# distroless/static enthält CA-Zertifikate, keine Shell/Package-Manager und
# läuft standardmäßig als nicht-root-User -> minimale Angriffsfläche.
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/blogger-xml-exporter ./blogger-xml-exporter
COPY --from=builder /src/config.yaml ./config.yaml
COPY --from=builder /src/templates ./templates
COPY --from=builder /src/web ./web

EXPOSE 8080

ENTRYPOINT ["/app/blogger-xml-exporter"]
