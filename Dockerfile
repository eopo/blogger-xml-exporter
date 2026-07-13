# syntax=docker/dockerfile:1

ARG VERSION=dev
ARG COMMIT_SHA=unknown
ARG BUILD_TIME=unknown

FROM node:22-alpine AS frontend-builder
WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend ./
RUN npm run build

FROM golang:1.26.5-alpine AS go-builder
ARG VERSION
ARG COMMIT_SHA
ARG BUILD_TIME
RUN go env -w GOTOOLCHAIN=auto
ENV GOTOOLCHAIN=auto
ENV CGO_ENABLED=0
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend ./
RUN go build \
    -ldflags "-X main.Version=${VERSION} -X main.CommitSHA=${COMMIT_SHA} -X main.BuildTime=${BUILD_TIME}" \
    -o /out/blogger-xml-exporter .

FROM gcr.io/distroless/static-debian12:nonroot@sha256:b7bb25d9f7c31d2bdd1982feb4dafcaf137703c7075dbe2febb41c24212b946f
WORKDIR /app
COPY --from=go-builder /out/blogger-xml-exporter ./blogger-xml-exporter
COPY --from=frontend-builder /src/web/static ./web/static
ENV PORT=8080
ENV CONFIG_PATH=/config/config.yaml
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["./blogger-xml-exporter"]