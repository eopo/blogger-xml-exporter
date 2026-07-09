# syntax=docker/dockerfile:1

# --- Go build stage ------------------------------------------------------
# Only Go sources are copied here, so front-end changes never rebuild the binary.
FROM golang:1.26-alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS go-builder

# Static binary for the distroless/static runtime, which ships no libc.
ENV CGO_ENABLED=0

WORKDIR /src

# Dependencies (cached while go.mod/go.sum unchanged).
COPY go.mod go.sum* ./
RUN go mod download

COPY main.go ./
COPY internal/ ./internal/
RUN go build -o /out/blogger-xml-exporter .

# --- CSS build stage -----------------------------------------------------
# Reuses the go-builder base image (already pulled) and downloads the official
# Tailwind CLI — same source as local `make`, so no version drift.
FROM golang:1.26-alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS css-builder

# libstdc++ is required to run the prebuilt Tailwind CLI (musl variant).
RUN apk add --no-cache make curl libstdc++

WORKDIR /src

# Tailwind CLI download stays cached until the Makefile changes.
COPY Makefile ./
RUN make setup-css-tools

# One dir copy of the whole web/ tree — no fragile file lists to keep in sync.
# Vendored assets rarely change, so letting them share this layer keeps the
# Dockerfile simple; a web/ change re-runs only the fast Tailwind build.
COPY web ./web
RUN make build-css

# --- Runtime stage -------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot@sha256:b7bb25d9f7c31d2bdd1982feb4dafcaf137703c7075dbe2febb41c24212b946f

WORKDIR /app

# config.yaml is intentionally not baked in — it must be mounted at runtime.
COPY --from=go-builder /out/blogger-xml-exporter ./blogger-xml-exporter

# The css-builder stage already holds the full web/ tree plus the compiled CSS.
COPY --from=css-builder /src/web ./web

EXPOSE 8080

ENTRYPOINT ["/app/blogger-xml-exporter"]
