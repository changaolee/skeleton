#!/usr/bin/env bash

# Common utilities, variables and checks for all build scripts.
set -o errexit
set +o nounset
set -o pipefail

# Sourced flag
COMMON_SOURCED=true

# The root of the build/dist directory
SKT_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
source "${SKT_ROOT}/scripts/lib/init.sh"
source "${SKT_ROOT}/scripts/install/environment.sh"

# 不输入密码执行需要 root 权限的命令
function skt::common::sudo {
  echo ${LINUX_PASSWORD} | sudo -S $1
}
