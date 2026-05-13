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
	&& fieldmapgen -model-dir internal/database -root User -out gen/database/fieldmap_generated.go -package database \
	&& patchfieldgen -file internal/feature/user/user.go -root CreateParams -out internal/feature/user/patch_generated.go -package user -fieldmap-import grpc-sandbox/gen/database -root-selector params.User -paths-selector params.Fields -bucket root:userFields:fieldmap.IsUserRootField -bucket profile:profileFields:fieldmap.IsUserProfileField -bucket profile.address:addressFields:fieldmap.IsUserAddressField -copy params.User.Profile:data.profile -copy params.User.Profile.Address:data.address:params.User.Profile

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
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
