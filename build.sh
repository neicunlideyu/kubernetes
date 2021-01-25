#!/bin/bash
make clean

export HTTPS_PROXY=http://bj-rd-proxy.byted.org:3128
HOSTARCH=$(go env GOHOSTARCH)
export GOARCH=${HOSTARCH}

case ${ARCH} in
  aarch64)
    export GOARCH=arm64
    ;;
  *)
    export GOARCH=amd64
    ;;
esac

GOOS=$(go env GOOS)
if [[ "${HOSTARCH}" == "${GOARCH}" ]]; then
    TARGET="_output/local/go/bin/*"
else
    TARGET="_output/local/go/bin/${GOOS}_${GOARCH}/*"
fi


make all KUBE_BUILD_PLATFORMS=linux/${GOARCH} && \
  mkdir -p output && \
  mv $TARGET output/
