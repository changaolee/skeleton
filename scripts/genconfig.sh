#!/usr/bin/env bash

# 本脚本功能：根据 scripts/install/environment.sh 配置，生成 SKT 组件 YAML 配置文件。
# 示例：genconfig.sh scripts/install/environment.sh configs/skt-apiserver.yaml

env_file="$1"
template_file="$2"

SKT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${SKT_ROOT}/scripts/lib/init.sh"

if [ $# -ne 2 ]; then
  skt::log::error "Usage: genconfig.sh scripts/environment.sh configs/skt-apiserver.yaml"
  exit 1
fi

source "${env_file}"

declare -A envs

set +u
for env in $(sed -n 's/^[^#].*${\(.*\)}.*/\1/p' ${template_file}); do
  if [ -z "$(eval echo \$${env})" ]; then
    skt::log::error "environment variable '${env}' not set"
    missing=true
  fi
done

if [ "${missing}" ]; then
  skt::log::error 'You may run `source scripts/environment.sh` to set these environment'
  exit 1
fi

eval "cat << EOF
$(cat ${template_file})
EOF"
