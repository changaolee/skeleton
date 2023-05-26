#!/usr/bin/env bash

# Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/changaolee/skeleton.

set -o errexit
set +o nounset
set -o pipefail

# Unset CDPATH so that path interpolation can work correctly
unset CDPATH

# Default use go modules
export GO111MODULE=on

# The root of the build/dist directory
SKT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

source "${SKT_ROOT}/scripts/lib/util.sh"
source "${SKT_ROOT}/scripts/lib/logging.sh"
source "${SKT_ROOT}/scripts/lib/color.sh"

skt::log::install_errexit

source "${SKT_ROOT}/scripts/lib/version.sh"
source "${SKT_ROOT}/scripts/lib/golang.sh"
