FROM cgr.dev/chainguard/go AS builder
ARG TARGETOS
ARG TARGETARCH
ARG GH_ACTOR=nduyphuong
ARG GH_REPO=gorya
LABEL org.opencontainers.image.source https://github.com/${GH_ACTOR}/${GH_REPO}
LABEL org.opencontainers.image.licenses MIT
ARG VERSION_PACKAGE=github.com/nduyphuong/gorya/internal/version

ARG VERSION
ARG GIT_COMMIT
ARG GIT_TREE_STATE
COPY . /app

RUN cd /app && GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags "-w -X ${VERSION_PACKAGE}.version=${VERSION} -X ${VERSION_PACKAGE}.buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ') -X ${VERSION_PACKAGE}.gitCommit=${GIT_COMMIT} -X ${VERSION_PACKAGE}.gitTreeState=${GIT_TREE_STATE}" \
    -o gorya ./cmd \
    && ./gorya version

FROM cgr.dev/chainguard/glibc-dynamic
COPY --from=builder /app/gorya /usr/bin/
CMD ["/usr/bin/gorya", "api"]
