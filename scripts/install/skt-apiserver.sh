#!/usr/bin/env bash

# Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/changaolee/skeleton.


# The root of the build/dist directory
SKT_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${SKT_ROOT}/scripts/install/common.sh

# 安装后打印必要的信息
function skt::apiserver::info() {
  cat <<EOF
skt-apserver insecure listen on: ${SKT_APISERVER_HOST}:${SKT_APISERVER_INSECURE_BIND_PORT}
skt-apserver secure listen on: ${SKT_APISERVER_HOST}:${SKT_APISERVER_SECURE_BIND_PORT}
EOF
}

# 安装
function skt::apiserver::install() {
  pushd ${SKT_ROOT}

  # 1. 生成 CA 证书和私钥
  ./scripts/gencerts.sh generate-skt-cert ${LOCAL_OUTPUT_ROOT}/cert
  skt::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/cert/ca* ${SKT_CONFIG_DIR}/cert"

  ./scripts/gencerts.sh generate-skt-cert ${LOCAL_OUTPUT_ROOT}/cert skt-apiserver
  skt::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/cert/skt-apiserver*pem ${SKT_CONFIG_DIR}/cert"

  # 2. 构建 skt-apiserver
  make build BINS=skt-apiserver
  skt::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/platforms/linux/amd64/skt-apiserver ${SKT_INSTALL_DIR}/bin"

  # 3.  生成并安装 skt-apiserver 的配置文件（skt-apiserver.yaml）
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} configs/skt-apiserver.yaml > ${SKT_CONFIG_DIR}/skt-apiserver.yaml"

  # 4. 创建并安装 skt-apiserver systemd unit 文件
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} init/skt-apiserver.service > /etc/systemd/system/skt-apiserver.service"

  # 5. 启动 skt-apiserver 服务
  skt::common::sudo "systemctl daemon-reload"
  skt::common::sudo "systemctl restart skt-apiserver"
  skt::common::sudo "systemctl enable skt-apiserver"
  skt::apiserver::status || return 1
  skt::apiserver::info

  skt::log::info "install skt-apiserver successfully"
  popd
}

# 卸载
function skt::apiserver::uninstall() {
  set +o errexit
  skt::common::sudo "systemctl stop skt-apiserver"
  skt::common::sudo "systemctl disable skt-apiserver"
  skt::common::sudo "rm -f ${SKT_INSTALL_DIR}/bin/skt-apiserver"
  skt::common::sudo "rm -f ${SKT_CONFIG_DIR}/skt-apiserver.yaml"
  skt::common::sudo "rm -f /etc/systemd/system/skt-apiserver.service"
  skt::common::sudo "rm -f ${SKT_CONFIG_DIR}/cert/skt-apiserver*pem"
  set -o errexit
  skt::log::info "uninstall skt-apiserver successfully"
}

# 状态检查
function skt::apiserver::status() {
  # 查看 apiserver 运行状态，如果输出中包含 active (running) 字样说明 apiserver 成功启动。
  systemctl status skt-apiserver | grep -q 'active' || {
    skt::log::error "skt-apiserver failed to start, maybe not installed properly"
    return 1
  }

  if echo | telnet ${SKT_APISERVER_HOST} ${SKT_APISERVER_INSECURE_BIND_PORT} 2>&1 | grep refused &>/dev/null; then
    skt::log::error "cannot access insecure port, skt-apiserver maybe not startup"
    return 1
  fi
}

if [[ "$*" =~ skt::apiserver:: ]]; then
  eval $*
fi
