#!/usr/bin/env bash

# Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/changaolee/skeleton.


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
