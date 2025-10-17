ARG GO_VERSION=1.24.1
FROM golang:${GO_VERSION}-alpine AS builder


ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

# metadata
LABEL org.opencontainers.image.created="${BUILD_DATE:-unknown}" \
    org.opencontainers.image.source="${VCS_REF:-unknown}"

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /app/main ./...

FROM alpine:3.19 AS runtime

RUN apk add --no-cache ca-certificates curl

RUN addgroup -S app && adduser -S -G app -h /app app

WORKDIR /app

COPY --from=builder /app/main /app/main

RUN chown app:app /app/main && chmod 755 /app/main

USER app

EXPOSE 8080 

ENTRYPOINT ["/app/main"]
