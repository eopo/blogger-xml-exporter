# syntax=docker/dockerfile:1

# Frontend build stage
FROM node:22-alpine AS frontend-builder

WORKDIR /src

COPY package.json package-lock.json ./
COPY frontend/package.json ./frontend/

RUN npm ci && cd frontend && npm ci

COPY frontend ./frontend
COPY web ./web

RUN npm run build

# Go build stage
FROM golang:1.24-alpine AS go-builder

RUN go env -w GOTOOLCHAIN=auto
ENV GOTOOLCHAIN=auto
ENV CGO_ENABLED=0

WORKDIR /src

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend ./
RUN go build -o /out/blogger-xml-exporter .

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot@sha256:b7bb25d9f7c31d2bdd1982feb4dafcaf137703c7075dbe2febb41c24212b946f

WORKDIR /app

COPY --from=go-builder /out/blogger-xml-exporter ./blogger-xml-exporter
COPY --from=frontend-builder /src/web/static ./web/static

# config.yaml must be mounted at runtime
ENV PORT=8080
ENV CONFIG_PATH=/config/config.yaml
EXPOSE 8080

ENTRYPOINT ["./blogger-xml-exporter"]


