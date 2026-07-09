# syntax=docker/dockerfile:1

# --- Build Stage ---------------------------------------------------------
FROM golang:1.26-alpine AS builder

# libstdc++ is required to run the prebuilt Tailwind CLI (musl variant).
RUN apk add --no-cache make curl libstdc++

# Static binary for the distroless/static runtime, which ships no libc.
ENV CGO_ENABLED=0

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

# Own layer so the Tailwind download is cached across source-code changes.
COPY Makefile ./
RUN make setup-css-tools

COPY . .
RUN make build

# --- Runtime Stage -------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# config.yaml is intentionally not baked in — it must be mounted at runtime.
COPY --from=builder /src/bin/blogger-xml-exporter ./blogger-xml-exporter
COPY --from=builder /src/web ./web

EXPOSE 8080

ENTRYPOINT ["/app/blogger-xml-exporter"]
