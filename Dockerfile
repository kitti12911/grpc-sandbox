# check=skip=InvalidDefaultArgInFrom
ARG TOOLCHAIN_IMAGE
FROM ${TOOLCHAIN_IMAGE} AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY buf.gen.yaml ./
COPY cmd ./cmd
COPY internal ./internal

RUN rm -rf gen/grpc gen/database \
	&& buf generate \
	&& mapgen fields \
	&& mapgen patch \
	&& mapgen filter \
	&& mapgen map

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
	go build -trimpath -ldflags="-s -w" -o /out/grpc-sandbox ./cmd/server

FROM alpine:3.23@sha256:fd791d74b68913cbb027c6546007b3f0d3bc45125f797758156952bc2d6daf40

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S app \
	&& adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/grpc-sandbox /app/grpc-sandbox
COPY --chown=app:app config.example.yml /app/config.yml

USER app

EXPOSE 50051

ENTRYPOINT ["/app/grpc-sandbox"]
