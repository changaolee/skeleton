#!/usr/bin/env bash

# 根目录
SKT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${SKT_ROOT}/scripts/lib/init.sh"

# 证书的主机名
readonly CERT_HOSTNAME="${CERT_HOSTNAME:-skt.api.changaolee.com,skt.authz.changaolee.com},127.0.0.1,localhost"

# 运行 cfssl 命令为 skt 服务生成证书文件.
#
# 参数:
#   $1 (证书文件存放目录)
#   $2 (证书文件名前缀)
function generate-skt-cert() {
  local cert_dir=${1}
  local prefix=${2:-}

  mkdir -p "${cert_dir}"
  pushd "${cert_dir}" || return 1

  skt::util::ensure-cfssl

  if [ ! -r "ca-config.json" ]; then
    cat >ca-config.json <<EOF
{
  "signing": {
    "default": {
      "expiry": "87600h"
    },
    "profiles": {
      "skt": {
        "usages": [
          "signing",
          "key encipherment",
          "server auth",
          "client auth"
        ],
        "expiry": "876000h"
      }
  }
}
}
EOF
  fi

  if [ ! -r "ca-csr.json" ]; then
    cat >ca-csr.json <<EOF
{
  "CN": "skt-ca",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "changaolee",
      "OU": "skt"
    }
  ],
  "ca": {
    "expiry": "876000h"
  }
}
EOF
  fi

  if [[ ! -r "ca.pem" || ! -r "ca-key.pem" ]]; then
    ${CFSSL_BIN} gencert -initca ca-csr.json | ${CFSSLJSON_BIN} -bare ca -
  fi

  if [[ -z "${prefix}" ]]; then
    return 0
  fi

  echo "Generate "${prefix}" certificates..."
  echo '{"CN":"'"${prefix}"'","hosts":[],"key":{"algo":"rsa","size":2048},"names":[{"C":"CN","ST":"BeiJing","L":"BeiJing","O":"changaolee","OU":"'"${prefix}"'"}]}' |
    ${CFSSL_BIN} gencert -hostname="${CERT_HOSTNAME},${prefix}" -ca=ca.pem -ca-key=ca-key.pem \
      -config=ca-config.json -profile=skt - | ${CFSSLJSON_BIN} -bare "${prefix}"

  # the popd will access `directory stack`, no `real` parameters is actually needed
  popd || return 1
}
