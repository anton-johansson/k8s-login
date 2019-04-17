#!/bin/bash

k8s-login::utility::git-tag() {
    local fallback=${1}
    local value=""
    if ! value=$(git describe --tags --exact-match 2> /dev/null); then
        value="${fallback}"
    fi
    echo "${value}"
}

k8s-login::version::ldflag() {
    local key=${1}
    local value=${2}

    echo "-X github.com/anton-johansson/k8s-login/version.${key}=${value}"
}

k8s-login::version::ldflags() {
    local gitTag=$(k8s-login::utility::git-tag "dev")
    local ldflags=$(k8s-login::version::ldflag "gitTag" "${gitTag}")

    echo "${ldflags[*]-}"
}

k8s-login::build::native-target() {
    local os=${1}
    local arch=${2}
    local ext=""
    if [ ${os} = "windows" ]; then
        ext=".exe"
    fi

    echo "Building ${os}/${arch}"
    GOOS=${os}
    GOARCH=${arch}
    go build -o build/k8s-login-${os}-${arch}${ext} -ldflags "$(k8s-login::version::ldflags)" .
}

k8s-login::build::native() {
    k8s-login::build::native-target "darwin" "amd64"
    k8s-login::build::native-target "linux" "amd64"
    k8s-login::build::native-target "windows" "amd64"
}

k8s-login::build::native
