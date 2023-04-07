# ==============================================================================
# 定义全局 Makefile 变量方便后面引用

SHELL := /bin/bash

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

# 项目根目录
ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../../ && pwd -P))
endif

# 构建产物、临时文件存放目录
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
endif
ifeq ($(origin CERT_DIR),undefined)
CERT_DIR := $(OUTPUT_DIR)/cert
$(shell mkdir -p $(CERT_DIR))
endif
ifeq ($(origin TOOLS_DIR),undefined)
TOOLS_DIR := $(OUTPUT_DIR)/tools
$(shell mkdir -p $(TOOLS_DIR))
endif
ifeq ($(origin TMP_DIR),undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
$(shell mkdir -p $(TMP_DIR))
endif

# 定义 VERSION 语义化版本号
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif

# 检查代码仓库是否是 dirty（默认 dirty）
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

# 最小测试覆盖率
ifeq ($(origin COVERAGE),undefined)
COVERAGE := 60
endif

# 编译的操作系统可以是 linux/windows/darwin
PLATFORMS ?= darwin_amd64 windows_amd64 linux_amd64 linux_arm64

# 设置一个指定的操作系统
ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
endif

# Linux 命令设置
FIND := find . ! -path './third_party/*' ! -path './vendor/*'
XARGS := xargs --no-run-if-empty

# Makefile 设置
ifndef V
MAKEFLAGS += --no-print-directory
endif

# 在执行 makefile 时复制 githook 脚本
COPY_GITHOOK:=$(shell cp -f githooks/* .git/hooks/)

# 指定需要证书的组件
ifeq ($(origin CERTIFICATES),undefined)
CERTIFICATES=skt-apiserver skt-authz-server
endif

# 设置工具的严重级别: BLOCKER_TOOLS, CRITICAL_TOOLS, TRIVIAL_TOOLS.
# 缺失 BLOCKER_TOOLS 会导致 CI 执行失败, 如：`make all` 失败.
# 缺失 CRITICAL_TOOLS 会导致一些必要操作失败. 如：`make release` 失败.
# TRIVIAL_TOOLS 是可选工具，缺失不会有影响.
BLOCKER_TOOLS ?= gsemver golines go-junit-report golangci-lint addlicense goimports
CRITICAL_TOOLS ?= swagger mockgen gotests git-chglog github-release go-mod-outdated protoc-gen-go go-gitlint
TRIVIAL_TOOLS ?= depth go-callvis gothanks richgo rts kube-score