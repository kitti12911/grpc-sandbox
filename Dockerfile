FROM golang:1.26.2-alpine@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1 AS builder

WORKDIR /src

ARG BUILDOS
ARG BUILDARCH
ARG TARGETOS
ARG TARGETARCH

COPY go.mod go.sum ./
RUN --mount=type=cache,id=gomodcache,target=/go/pkg/mod \
	--mount=type=cache,id=gobuild-${BUILDOS}-${BUILDARCH},target=/root/.cache/go-build \
	go mod download

RUN apk add --no-cache git make

RUN --mount=type=cache,id=gomodcache,target=/go/pkg/mod \
	--mount=type=cache,id=gobuild-${BUILDOS}-${BUILDARCH},target=/root/.cache/go-build \
	go install github.com/bufbuild/buf/cmd/buf@v1.69.0 \
	&& go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11 \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.1

COPY Makefile buf.gen.yaml ./
COPY cmd ./cmd
COPY internal ./internal

RUN --mount=type=cache,id=gomodcache,target=/go/pkg/mod \
	--mount=type=cache,id=gobuild-${BUILDOS}-${BUILDARCH},target=/root/.cache/go-build \
	make gen

RUN --mount=type=cache,id=gomodcache,target=/go/pkg/mod \
	--mount=type=cache,id=gobuild-${TARGETOS}-${TARGETARCH},target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
	go build -trimpath -ldflags="-s -w" -o /out/grpc-sandbox ./cmd/server

FROM alpine:3.22@sha256:310c62b5e7ca5b08167e4384c68db0fd2905dd9c7493756d356e893909057601

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S app \
	&& adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/grpc-sandbox /app/grpc-sandbox
COPY --chown=app:app config.example.yml /app/config.yml

USER app

EXPOSE 50051

ENTRYPOINT ["/app/grpc-sandbox"]
