#!/usr/bin/env bash

# Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/changaolee/skeleton.


# shellcheck disable=SC2034 # Variables sourced in other scripts.

# The server platform we are building on.
readonly SKT_SUPPORTED_SERVER_PLATFORMS=(
  linux/amd64
  linux/arm64
)

# If we update this we should also update the set of platforms whose standard
# library is precompiled for in build/build-image/cross/Dockerfile
readonly SKT_SUPPORTED_CLIENT_PLATFORMS=(
  linux/amd64
  linux/arm64
)

# The set of server targets that we are only building for Linux
# If you update this list, please also update build/BUILD.
skt::golang::server_targets() {
  local targets=(
    skt-apiserver
    skt-authz-server
    skt-pump
    skt-watcher
  )
  echo "${targets[@]}"
}

IFS=" " read -ra SKT_SERVER_TARGETS <<<"$(skt::golang::server_targets)"
readonly SKT_SERVER_TARGETS
readonly SKT_SERVER_BINARIES=("${SKT_SERVER_TARGETS[@]##*/}")

# The set of server targets we build docker images for
skt::golang::server_image_targets() {
  # NOTE: this contains cmd targets for skt::build::get_docker_wrapped_binaries
  local targets=(
    cmd/skt-apiserver
    cmd/skt-authz-server
    cmd/skt-pump
    cmd/skt-watcher
  )
  echo "${targets[@]}"
}

IFS=" " read -ra SKT_SERVER_IMAGE_TARGETS <<<"$(skt::golang::server_image_targets)"
readonly SKT_SERVER_IMAGE_TARGETS
readonly SKT_SERVER_IMAGE_BINARIES=("${SKT_SERVER_IMAGE_TARGETS[@]##*/}")

# ------------
# NOTE: All functions that return lists should use newlines.
# bash functions can't return arrays, and spaces are tricky, so newline
# separators are the preferred pattern.
# To transform a string of newline-separated items to an array, use skt::util::read-array:
# skt::util::read-array FOO < <(skt::golang::dups a b c a)
#
# ALWAYS remember to quote your subshells. Not doing so will break in
# bash 4.3, and potentially cause other issues.
# ------------

# Returns a sorted newline-separated list containing only duplicated items.
skt::golang::dups() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort | uniq -d
}

# Returns a sorted newline-separated list with duplicated items removed.
skt::golang::dedup() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort -u
}

# Depends on values of user-facing SKT_BUILD_PLATFORMS, SKT_FASTBUILD,
# and SKT_BUILDER_OS.
# Configures SKT_SERVER_PLATFORMS and SKT_CLIENT_PLATFORMS, then sets them
# to readonly.
# The configured vars will only contain platforms allowed by the
# SKT_SUPPORTED* vars at the top of this file.
declare -a SKT_SERVER_PLATFORMS
declare -a SKT_CLIENT_PLATFORMS
skt::golang::setup_platforms() {
  if [[ -n "${SKT_BUILD_PLATFORMS:-}" ]]; then
    # SKT_BUILD_PLATFORMS needs to be read into an array before the next
    # step, or quoting treats it all as one element.
    local -a platforms
    IFS=" " read -ra platforms <<<"${SKT_BUILD_PLATFORMS}"

    # Deduplicate to ensure the intersection trick with skt::golang::dups
    # is not defeated by duplicates in user input.
    skt::util::read-array platforms < <(skt::golang::dedup "${platforms[@]}")

    # Use skt::golang::dups to restrict the builds to the platforms in
    # SKT_SUPPORTED_*_PLATFORMS. Items should only appear at most once in each
    # set, so if they appear twice after the merge they are in the intersection.
    skt::util::read-array SKT_SERVER_PLATFORMS < <(
      skt::golang::dups \
        "${platforms[@]}" \
        "${SKT_SUPPORTED_SERVER_PLATFORMS[@]}"
    )
    readonly SKT_SERVER_PLATFORMS

    skt::util::read-array SKT_CLIENT_PLATFORMS < <(
      skt::golang::dups \
        "${platforms[@]}" \
        "${SKT_SUPPORTED_CLIENT_PLATFORMS[@]}"
    )
    readonly SKT_CLIENT_PLATFORMS

  elif [[ "${SKT_FASTBUILD:-}" == "true" ]]; then
    SKT_SERVER_PLATFORMS=(linux/amd64)
    readonly SKT_SERVER_PLATFORMS
    SKT_CLIENT_PLATFORMS=(linux/amd64)
    readonly SKT_CLIENT_PLATFORMS
  else
    SKT_SERVER_PLATFORMS=("${SKT_SUPPORTED_SERVER_PLATFORMS[@]}")
    readonly SKT_SERVER_PLATFORMS

    SKT_CLIENT_PLATFORMS=("${SKT_SUPPORTED_CLIENT_PLATFORMS[@]}")
    readonly SKT_CLIENT_PLATFORMS
  fi
}

skt::golang::setup_platforms

# The set of client targets that we are building for all platforms
# If you update this list, please also update build/BUILD.
readonly SKT_CLIENT_TARGETS=(
  sktctl
)
readonly SKT_CLIENT_BINARIES=("${SKT_CLIENT_TARGETS[@]##*/}")

readonly SKT_ALL_TARGETS=(
  "${SKT_SERVER_TARGETS[@]}"
  "${SKT_CLIENT_TARGETS[@]}"
)
readonly SKT_ALL_BINARIES=("${SKT_ALL_TARGETS[@]##*/}")

# Asks golang what it thinks the host platform is. The go tool chain does some
# slightly different things when the target platform matches the host platform.
skt::golang::host_platform() {
  echo "$(go env GOHOSTOS)/$(go env GOHOSTARCH)"
}

# Ensure the go tool exists and is a viable version.
skt::golang::verify_go_version() {
  if [[ -z "$(command -v go)" ]]; then
    skt::log::usage_from_stdin <<EOF
Can't find 'go' in PATH, please fix and retry.
See http://golang.org/doc/install for installation instructions.
EOF
    return 2
  fi

  local go_version
  IFS=" " read -ra go_version <<<"$(go version)"
  local minimum_go_version
  minimum_go_version=go1.13.4
  if [[ "${minimum_go_version}" != $(echo -e "${minimum_go_version}\n${go_version[2]}" | sort -s -t. -k 1,1 -k 2,2n -k 3,3n | head -n1) && "${go_version[2]}" != "devel" ]]; then
    skt::log::usage_from_stdin <<EOF
Detected go version: ${go_version[*]}.
SKT requires ${minimum_go_version} or greater.
Please install ${minimum_go_version} or later.
EOF
    return 2
  fi
}

# skt::golang::setup_env will check that the `go` commands is available in
# ${PATH}. It will also check that the Go version is good enough for the
# SKT build.
#
# Outputs:
#   env-var GOBIN is unset (we want binaries in a predictable place)
#   env-var GO15VENDOREXPERIMENT=1
#   env-var GO111MODULE=on
skt::golang::setup_env() {
  skt::golang::verify_go_version

  # Unset GOBIN in case it already exists in the current session.
  unset GOBIN

  # This seems to matter to some tools
  export GO15VENDOREXPERIMENT=1

  # Open go module feature
  export GO111MODULE=on

  # This is for sanity.  Without it, user umasks leak through into release
  # artifacts.
  umask 0022
}
