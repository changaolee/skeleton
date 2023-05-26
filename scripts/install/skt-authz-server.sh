#!/usr/bin/env bash

# The root of the build/dist directory
SKT_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${SKT_ROOT}/scripts/install/common.sh

# 安装后打印必要的信息
function skt::authzserver::info() {
  cat <<EOF
skt-authz-server insecure listen on: ${SKT_AUTHZ_SERVER_HOST}:${SKT_AUTHZ_SERVER_INSECURE_BIND_PORT}
skt-authz-server secure listen on: ${SKT_AUTHZ_SERVER_HOST}:${SKT_AUTHZ_SERVER_SECURE_BIND_PORT}
EOF
}

# 安装
function skt::authzserver::install() {
  pushd ${SKT_ROOT}

  # 1. 生成 CA 证书和私钥
  ./scripts/gencerts.sh generate-skt-cert ${LOCAL_OUTPUT_ROOT}/cert
  skt::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/cert/ca* ${SKT_CONFIG_DIR}/cert"

  ./scripts/gencerts.sh generate-skt-cert ${LOCAL_OUTPUT_ROOT}/cert skt-authz-server
  skt::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/cert/skt-authz-server*pem ${SKT_CONFIG_DIR}/cert"

  # 2. 构建 skt-authz-server
  make build BINS=skt-authz-server
  skt::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/platforms/linux/amd64/skt-authz-server ${SKT_INSTALL_DIR}/bin"

  # 3. 生成并安装 skt-authz-server 的配置文件（skt-authz-server.yaml）
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} configs/skt-authz-server.yaml > ${SKT_CONFIG_DIR}/skt-authz-server.yaml"

  # 4. 创建并安装 skt-authz-server systemd unit 文件
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} init/skt-authz-server.service > /etc/systemd/system/skt-authz-server.service"

  # 5. 启动 skt-authz-server 服务
  skt::common::sudo "systemctl daemon-reload"
  skt::common::sudo "systemctl restart skt-authz-server"
  skt::common::sudo "systemctl enable skt-authz-server"
  skt::authzserver::status || return 1
  skt::authzserver::info

  skt::log::info "install skt-authz-server successfully"
  popd
}

# 卸载
function skt::authzserver::uninstall() {
  set +o errexit
  skt::common::sudo "systemctl stop skt-authz-server"
  skt::common::sudo "systemctl disable skt-authz-server"
  skt::common::sudo "rm -f ${SKT_INSTALL_DIR}/bin/skt-authz-server"
  skt::common::sudo "rm -f ${SKT_CONFIG_DIR}/skt-authz-server.yaml"
  skt::common::sudo "rm -f /etc/systemd/system/skt-authz-server.service"
  skt::common::sudo "rm -f ${SKT_CONFIG_DIR}/cert/skt-authz-server*pem"
  set -o errexit
  skt::log::info "uninstall skt-authz-server successfully"
}

# 状态检查
function skt::authzserver::status() {
  # 查看 skt-authz-server 运行状态，如果输出中包含 active (running) 字样说明 skt-authz-server 成功启动。
  systemctl status skt-authz-server | grep -q 'active' || {
    skt::log::error "skt-authz-server failed to start, maybe not installed properly"
    return 1
  }

  if echo | telnet ${SKT_AUTHZSERVER_HOST} ${SKT_AUTHZSERVER_INSECURE_BIND_PORT} 2>&1 | grep refused &>/dev/null; then
    skt::log::error "cannot access insecure port, skt-authz-server maybe not startup"
    return 1
  fi
}

if [[ "$*" =~ skt::authzserver:: ]]; then
  eval $*
fi
