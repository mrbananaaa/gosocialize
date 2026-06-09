FROM golang:1.26.1 AS build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

FROM build-base AS dev

RUN go install github.com/air-verse/air@latest && \
  go install github.com/go-delve/delve/cmd/dlv@latest

COPY . .

CMD ["air", "-c", ".air.docker.toml"]

FROM build-base AS build-production

RUN useradd -u 1001 nonroot

COPY . .

RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o api-golang

FROM scratch

WORKDIR /

COPY --from=build-production /etc/passwd /etc/passwd

COPY --from=build-production /app/healthcheck/healthcheck healthcheck

COPY --from=build-production /app/api-golang api-golang

USER nonroot

EXPOSE 8080

CMD ["/api-golang"]