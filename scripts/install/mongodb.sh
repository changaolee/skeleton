#!/usr/bin/env bash

# Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/changaolee/skeleton.

# The root of the build/dist directory
SKT_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source "${SKT_ROOT}/scripts/install/common.sh"

# 安装后打印必要的信息
function skt::mongodb::info() {
  cat <<EOF
MongoDB Login: mongo mongodb://${MONGO_USERNAME}:'${MONGO_PASSWORD}'@${MONGO_HOST}:${MONGO_PORT}/skt_analytics?authSource=skt_analytics
EOF
}

# 安装
function skt::mongodb::install() {
  # 1. 配置 MongoDB Yum 源
  echo ${LINUX_PASSWORD} | sudo -S bash -c "cat << 'EOF' > /etc/yum.repos.d/mongodb-org-5.0.repo
[mongodb-org-5.0]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/\$releasever/mongodb-org/5.0/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-5.0.asc
EOF"

  # 2. 安装 MongoDB 和 MongoDB 客户端
  skt::common::sudo "yum install -y mongodb-org"

  # 3. 禁用 SELinux
  echo ${LINUX_PASSWORD} | sudo -S setenforce 0 || true
  echo ${LINUX_PASSWORD} | sudo -S sed -i 's/^SELINUX=.*$/SELINUX=disabled/' /etc/selinux/config

  # 4. 开启外网访问权限和登录验证
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/bindIp/{s/127.0.0.1/0.0.0.0/}' /etc/mongod.conf
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^#security/a\security:\n  authorization: enabled' /etc/mongod.conf

  # 5. 启动 MongoDB，并设置开机启动
  skt::common::sudo "systemctl enable mongod"
  skt::common::sudo "systemctl start mongod"

  # 6. 创建管理员账号，设置管理员密码
  mongosh --quiet "mongodb://${MONGO_HOST}:${MONGO_PORT}" <<EOF
use admin
db.createUser({user:"${MONGO_ADMIN_USERNAME}",pwd:"${MONGO_ADMIN_PASSWORD}",roles:["root"]})
db.auth("${MONGO_ADMIN_USERNAME}", "${MONGO_ADMIN_PASSWORD}")
exit
EOF

  # 7. 创建 ${MONGO_USERNAME} 用户
  mongosh --quiet mongodb://${MONGO_ADMIN_USERNAME}:${MONGO_ADMIN_PASSWORD}@${MONGO_HOST}:${MONGO_PORT}/skt_analytics?authSource=admin <<EOF
use skt_analytics
db.createUser({user:"${MONGO_USERNAME}",pwd:"${MONGO_PASSWORD}",roles:["dbOwner"]})
db.auth("${MONGO_USERNAME}", "${MONGO_PASSWORD}")
exit
EOF

  skt::mongodb::status || return 1
  skt::mongodb::info
  skt::log::info "install MongoDB successfully"
}

# 卸载
function skt::mongodb::uninstall() {
  set +o errexit
  skt::common::sudo "systemctl stop mongodb"
  skt::common::sudo "systemctl disable mongodb"
  skt::common::sudo "yum -y remove mongodb-org"
  skt::common::sudo "rm -rf /var/lib/mongo"
  skt::common::sudo "rm -f /etc/yum.repos.d/mongodb-10.5.repo"
  skt::common::sudo "rm -f /etc/mongod.conf"
  skt::common::sudo "rm -f /lib/systemd/system/mongod.service"
  skt::common::sudo "rm -f /tmp/mongodb-*.sock"
  set -o errexit
  skt::log::info "uninstall MongoDB successfully"
}

# 状态检查
function skt::mongodb::status() {
  # 查看 mongodb 运行状态，如果输出中包含 active (running) 字样说明 mongodb 成功启动。
  systemctl status mongod | grep -q 'active' || {
    skt::log::error "mongodb failed to start, maybe not installed properly"
    return 1
  }

  echo "show dbs" | mongosh --quiet "mongodb://${MONGO_HOST}:${MONGO_PORT}" &>/dev/null || {
    skt::log::error "cannot connect to mongodb, mongo maybe not installed properly"
    return 1
  }

  echo "show dbs" |
    mongosh --quiet mongodb://${MONGO_ADMIN_USERNAME}:${MONGO_ADMIN_PASSWORD}@${MONGO_HOST}:${MONGO_PORT}/skt_analytics?authSource=admin &>/dev/null || {
    skt::log::error "can not login with ${MONGO_ADMIN_USERNAME}, mongo maybe not initialized properly"
    return 1
  }

  echo "show dbs" |
    mongosh --quiet mongodb://${MONGO_USERNAME}:${MONGO_PASSWORD}@${MONGO_HOST}:${MONGO_PORT}/skt_analytics?authSource=skt_analytics &>/dev/null || {
    skt::log::error "can not login with ${MONGO_USERNAME}, mongo maybe not initialized properly"
    return 1
  }
}

if [[ "$*" =~ skt::mongodb:: ]]; then
  eval $*
fi
