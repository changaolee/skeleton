#!/usr/bin/env bash

# The root of the build/dist directory
SKT_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${SKT_ROOT}/scripts/install/common.sh

# 安装后打印必要的信息
function skt::sktctl::info() {
  cat <<EOF
sktctl test command: sktctl user list
EOF
}

# 安装
function skt::sktctl::install() {
  pushd ${SKT_ROOT}

  # 1. 生成并安装 CA 证书和私钥
  ./scripts/gencerts.sh generate-skt-cert ${LOCAL_OUTPUT_ROOT}/cert
  skt::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/cert/ca* ${SKT_CONFIG_DIR}/cert"

  ./scripts/gencerts.sh generate-skt-cert ${LOCAL_OUTPUT_ROOT}/cert admin
  #skt::common::sudo "cp ${LOCAL_OUTPUT_ROOT}/cert/admin*pem ${SKT_CONFIG_DIR}/cert"
  cert_dir=$(dirname ${CONFIG_USER_CLIENT_CERTIFICATE})
  key_dir=$(dirname ${CONFIG_USER_CLIENT_KEY})
  mkdir -p ${cert_dir} ${key_dir}
  cp ${LOCAL_OUTPUT_ROOT}/cert/admin.pem ${CONFIG_USER_CLIENT_CERTIFICATE}
  cp ${LOCAL_OUTPUT_ROOT}/cert/admin-key.pem ${CONFIG_USER_CLIENT_KEY}

  # 2. 构建 sktctl
  make build BINS=sktctl
  cp ${LOCAL_OUTPUT_ROOT}/platforms/linux/amd64/sktctl $HOME/bin/

  # 3. 生成并安装 sktctl 的配置文件（sktctl.yaml）
  mkdir -p $HOME/.skt
  ./scripts/genconfig.sh ${ENV_FILE} configs/sktctl.yaml >$HOME/.skt/sktctl.yaml
  skt::sktctl::status || return 1
  skt::sktctl::info

  skt::log::info "install sktctl successfully"
  popd
}

# 卸载
function skt::sktctl::uninstall() {
  set +o errexit
  rm -f $HOME/bin/sktctl
  rm -f $HOME/.skt/sktctl.yaml
  #skt::common::sudo "rm -f ${SKT_CONFIG_DIR}/cert/admin*pem"
  rm -f ${CONFIG_USER_CLIENT_CERTIFICATE}
  rm -f ${CONFIG_USER_CLIENT_KEY}
  set -o errexit

  skt::log::info "uninstall sktctl successfully"
}

# 状态检查
function skt::sktctl::status() {
  sktctl user list | grep -q admin || {
    skt::log::error "cannot list user, sktctl maybe not installed properly"
    return 1
  }

  if echo | telnet ${SKT_APISERVER_HOST} ${SKT_APISERVER_INSECURE_BIND_PORT} 2>&1 | grep refused &>/dev/null; then
    skt::log::error "cannot access insecure port, sktctl maybe not startup"
    return 1
  fi
}

if [[ "$*" =~ skt::sktctl:: ]]; then
  eval $*
fi
