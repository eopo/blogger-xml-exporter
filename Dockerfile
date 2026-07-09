# syntax=docker/dockerfile:1

# Build arguments for versioning
ARG VERSION=main
ARG BUILD_DATE
ARG VCS_REF

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

# OCI Image Labels (org.opencontainers.image.*)
LABEL org.opencontainers.image.title="Blogger XML Exporter"
LABEL org.opencontainers.image.description="Export Blogger.com blog content to XML format"
LABEL org.opencontainers.image.vendor="eopo"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.revision="${VCS_REF}"
LABEL org.opencontainers.image.source="https://github.com/eopo/blogger-xml-exporter"
LABEL org.opencontainers.image.documentation="https://github.com/eopo/blogger-xml-exporter#readme"
LABEL org.opencontainers.image.licenses="ISC"

# config.yaml is intentionally not baked in — it must be mounted at runtime.
COPY --from=builder /src/bin/blogger-xml-exporter ./blogger-xml-exporter
COPY --from=builder /src/web ./web

EXPOSE 8080

ENTRYPOINT ["/app/blogger-xml-exporter"]
