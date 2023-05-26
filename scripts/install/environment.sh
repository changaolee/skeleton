#!/usr/bin/env bash

# Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/changaolee/skeleton.

# SKT 项目源码根目录
SKT_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

# 生成文件存放目录
LOCAL_OUTPUT_ROOT="${SKT_ROOT}/${OUT_DIR:-_output}"

# 设置统一的密码，方便记忆
readonly PASSWORD=${PASSWORD:-'SKT1826!z$'}

# Linux系统 goer 用户
readonly LINUX_USERNAME=${LINUX_USERNAME:-goer}
# Linux root & goer 用户密码
readonly LINUX_PASSWORD=${LINUX_PASSWORD:-${PASSWORD}}

# 设置安装目录
readonly INSTALL_DIR=${INSTALL_DIR:-/tmp/installation}
mkdir -p ${INSTALL_DIR}
readonly ENV_FILE=${SKT_ROOT}/scripts/install/environment.sh

# MariaDB 配置信息
readonly MARIADB_ADMIN_USERNAME=${MARIADB_ADMIN_USERNAME:-root}        # MariaDB root 用户
readonly MARIADB_ADMIN_PASSWORD=${MARIADB_ADMIN_PASSWORD:-${PASSWORD}} # MariaDB root 用户密码
readonly MARIADB_HOST=${MARIADB_HOST:-127.0.0.1:3306}                  # MariaDB 主机地址
readonly MARIADB_DATABASE=${MARIADB_DATABASE:-skeleton}                # MariaDB skeleton 应用使用的数据库名
readonly MARIADB_USERNAME=${MARIADB_USERNAME:-skt}                     # skeleton 数据库用户名
readonly MARIADB_PASSWORD=${MARIADB_PASSWORD:-${PASSWORD}}             # skeleton 数据库密码

# Redis 配置信息
readonly REDIS_HOST=${REDIS_HOST:-127.0.0.1}           # Redis 主机地址
readonly REDIS_PORT=${REDIS_PORT:-6379}                # Redis 监听端口
readonly REDIS_USERNAME=${REDIS_USERNAME:-''}          # Redis 用户名
readonly REDIS_PASSWORD=${REDIS_PASSWORD:-${PASSWORD}} # Redis 密码

# MongoDB 配置
readonly MONGO_ADMIN_USERNAME=${MONGO_ADMIN_USERNAME:-root}        # MongoDB root 用户
readonly MONGO_ADMIN_PASSWORD=${MONGO_ADMIN_PASSWORD:-${PASSWORD}} # MongoDB root 用户密码
readonly MONGO_HOST=${MONGO_HOST:-127.0.0.1}                       # MongoDB 地址
readonly MONGO_PORT=${MONGO_PORT:-27017}                           # MongoDB 端口
readonly MONGO_USERNAME=${MONGO_USERNAME:-skt}                     # MongoDB 用户名
readonly MONGO_PASSWORD=${MONGO_PASSWORD:-${PASSWORD}}             # MongoDB 密码

# skt 配置
readonly SKT_DATA_DIR=${SKT_DATA_DIR:-/data/skt}           # skt 各组件数据目录
readonly SKT_INSTALL_DIR=${SKT_INSTALL_DIR:-/opt/skt}      # skt 安装文件存放目录
readonly SKT_CONFIG_DIR=${SKT_CONFIG_DIR:-/etc/skt}        # skt 配置文件存放目录
readonly SKT_LOG_DIR=${SKT_LOG_DIR:-/var/log/skt}          # skt 日志文件存放目录
readonly CA_FILE=${CA_FILE:-${SKT_CONFIG_DIR}/cert/ca.pem} # CA

# skt-apiserver 配置
readonly SKT_APISERVER_HOST=${SKT_APISERVER_HOST:-127.0.0.1} # skt-apiserver 部署机器 IP 地址
readonly SKT_APISERVER_GRPC_BIND_ADDRESS=${SKT_APISERVER_GRPC_BIND_ADDRESS:-0.0.0.0}
readonly SKT_APISERVER_GRPC_BIND_PORT=${SKT_APISERVER_GRPC_BIND_PORT:-8081}
readonly SKT_APISERVER_INSECURE_BIND_ADDRESS=${SKT_APISERVER_INSECURE_BIND_ADDRESS:-127.0.0.1}
readonly SKT_APISERVER_INSECURE_BIND_PORT=${SKT_APISERVER_INSECURE_BIND_PORT:-8080}
readonly SKT_APISERVER_SECURE_BIND_ADDRESS=${SKT_APISERVER_SECURE_BIND_ADDRESS:-0.0.0.0}
readonly SKT_APISERVER_SECURE_BIND_PORT=${SKT_APISERVER_SECURE_BIND_PORT:-8443}
readonly SKT_APISERVER_SECURE_TLS_CERT_KEY_CERT_FILE=${SKT_APISERVER_SECURE_TLS_CERT_KEY_CERT_FILE:-${SKT_CONFIG_DIR}/cert/skt-apiserver.pem}
readonly SKT_APISERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE=${SKT_APISERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE:-${SKT_CONFIG_DIR}/cert/skt-apiserver-key.pem}

# skt-authz-server 配置
readonly SKT_AUTHZ_SERVER_HOST=${SKT_AUTHZ_SERVER_HOST:-127.0.0.1} # skt-authz-server 部署机器 IP 地址
readonly SKT_AUTHZ_SERVER_INSECURE_BIND_ADDRESS=${SKT_AUTHZ_SERVER_INSECURE_BIND_ADDRESS:-127.0.0.1}
readonly SKT_AUTHZ_SERVER_INSECURE_BIND_PORT=${SKT_AUTHZ_SERVER_INSECURE_BIND_PORT:-9090}
readonly SKT_AUTHZ_SERVER_SECURE_BIND_ADDRESS=${SKT_AUTHZ_SERVER_SECURE_BIND_ADDRESS:-0.0.0.0}
readonly SKT_AUTHZ_SERVER_SECURE_BIND_PORT=${SKT_AUTHZ_SERVER_SECURE_BIND_PORT:-9443}
readonly SKT_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_CERT_FILE=${SKT_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_CERT_FILE:-${SKT_CONFIG_DIR}/cert/skt-authz-server.pem}
readonly SKT_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE=${SKT_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE:-${SKT_CONFIG_DIR}/cert/skt-authz-server-key.pem}
readonly SKT_AUTHZ_SERVER_CLIENT_CA_FILE=${SKT_AUTHZ_SERVER_CLIENT_CA_FILE:-${CA_FILE}}
readonly SKT_AUTHZ_SERVER_RPCSERVER=${SKT_AUTHZ_SERVER_RPCSERVER:-${SKT_APISERVER_HOST}:${SKT_APISERVER_GRPC_BIND_PORT}}

# skt-pump 配置
readonly SKT_PUMP_HOST=${SKT_PUMP_HOST:-127.0.0.1} # skt-pump 部署机器 IP 地址
readonly SKT_PUMP_COLLECTION_NAME=${SKT_PUMP_COLLECTION_NAME:-skt_analytics}
readonly SKT_PUMP_MONGO_URL=${SKT_PUMP_MONGO_URL:-mongodb://${MONGO_USERNAME}:${MONGO_PASSWORD}@${MONGO_HOST}:${MONGO_PORT}/${SKT_PUMP_COLLECTION_NAME}?authSource=${SKT_PUMP_COLLECTION_NAME}}

# sktctl 配置
readonly CONFIG_USER_USERNAME=${CONFIG_USER_USERNAME:-admin}
readonly CONFIG_USER_PASSWORD=${CONFIG_USER_PASSWORD:-Admin@2023}
readonly CONFIG_USER_CLIENT_CERTIFICATE=${CONFIG_USER_CLIENT_CERTIFICATE:-${HOME}/.skt/cert/admin.pem}
readonly CONFIG_USER_CLIENT_KEY=${CONFIG_USER_CLIENT_KEY:-${HOME}/.skt/cert/admin-key.pem}
readonly CONFIG_SERVER_ADDRESS=${CONFIG_SERVER_ADDRESS:-${SKT_APISERVER_HOST}:${SKT_APISERVER_SECURE_BIND_PORT}}
readonly CONFIG_SERVER_CERTIFICATE_AUTHORITY=${CONFIG_SERVER_CERTIFICATE_AUTHORITY:-${CA_FILE}}
