#!/usr/bin/env bash

# The root of the build/dist directory
SKT_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
source "${SKT_ROOT}/scripts/install/common.sh"

# 准备 Linux 环境
skt::install::prepare_linux() {
  # 1. 替换 Yum 源为阿里的 Yum 源
  skt::common::sudo "mv /etc/yum.repos.d /etc/yum.repos.d.$$.bak" # 先备份原有的 Yum 源
  skt::common::sudo "mkdir /etc/yum.repos.d"
  skt::common::sudo "wget -O /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-vault-8.5.2111.repo"
  skt::common::sudo "yum clean all"
  skt::common::sudo "yum makecache"

  if [[ -f $HOME/.bashrc ]]; then
    cp $HOME/.bashrc $HOME/bashrc.skt.backup
  fi

  # 2. 配置 $HOME/.bashrc
  cat <<'EOF' >$HOME/.bashrc
# .bashrc

# User specific aliases and functions

alias rm='rm -i'
alias cp='cp -i'
alias mv='mv -i'

# Source global definitions
if [ -f /etc/bashrc ]; then
    . /etc/bashrc
fi

if [ ! -d $HOME/workspace ]; then
    mkdir -p $HOME/workspace
fi

# User specific environment
# Basic envs
export LANG="en_US.UTF-8" # 设置系统语言为 en_US.UTF-8，避免终端出现中文乱码
export PS1='[\u@dev \W]\$ ' # 默认的 PS1 设置会展示全部的路径，为了防止过长，这里只展示："用户名@dev 最后的目录名"
export WORKSPACE="$HOME/workspace" # 设置工作目录
export PATH=$HOME/bin:$PATH # 将 $HOME/bin 目录加入到 PATH 变量中

# Default entry folder
cd $WORKSPACE # 登录系统，默认进入 workspace 目录

# User specific aliases and functions
EOF

  # 3. 安装依赖包
  skt::common::sudo "yum -y install make autoconf automake cmake perl-CPAN libcurl-devel libtool gcc gcc-c++ glibc-headers zlib-devel git-lfs telnet lrzsz jq expat-devel openssl-devel"

  # 4. 安装 Git
  rm -rf /tmp/git-2.36.1.tar.gz /tmp/git-2.36.1 # clean up
  cd /tmp || return 1
  wget --no-check-certificate https://mirrors.edge.kernel.org/pub/software/scm/git/git-2.36.1.tar.gz
  tar -xvzf git-2.36.1.tar.gz
  cd git-2.36.1/ || return 1
  ./configure
  make
  skt::common::sudo "make install"

  cat <<'EOF' >>$HOME/.bashrc
# Configure for git
export PATH=/usr/local/libexec/git-core:$PATH
EOF

  git --version | grep -q 'git version 2.36.1' || {
    skt::log::error "git version is not '2.36.1', maynot install git properly"
    return 1
  }

  # 5. 配置 Git
  git config --global user.name "changaolee"                   # 用户名改成自己的
  git config --global user.email "changao.li.work@outlook.com" # 邮箱改成自己的
  git config --global credential.helper store                  # 设置 Git，保存用户名和密码
  git config --global core.longpaths true                      # 解决 Git 中 'Filename too long' 的错误
  git config --global core.quotepath off
  git lfs install --skip-repo

  source $HOME/.bashrc
  skt::log::info "prepare linux basic environment successfully"
}

# 安装 Go 命令
function skt::install::go_command() {
  rm -rf /tmp/go1.18.3.linux-amd64.tar.gz $HOME/go/go1.18.3 # clean up

  # 1. 下载 go1.18.3 版本的 Go 安装包
  wget -P /tmp/ https://golang.google.cn/dl/go1.18.3.linux-amd64.tar.gz

  # 2. 安装 Go
  mkdir -p $HOME/go
  tar -xvzf /tmp/go1.18.3.linux-amd64.tar.gz -C $HOME/go
  mv $HOME/go/go $HOME/go/go1.18.3

  # 3. 配置 Go 环境变量
  cat <<'EOF' >>$HOME/.bashrc
# Go envs
export GOVERSION=go1.18.3 # Go 版本设置
export GO_INSTALL_DIR=$HOME/go # Go 安装目录
export GOROOT=$GO_INSTALL_DIR/$GOVERSION # GOROOT 设置
export GOPATH=$WORKSPACE/golang # GOPATH 设置
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH # 将 Go 语言自带的和通过 go install 安装的二进制文件加入到 PATH 路径中
export GO111MODULE="on" # 开启 Go moudles 特性
export GOPROXY=https://goproxy.cn,direct # 安装 Go 模块时，代理服务器设置
export GOPRIVATE=
export GOSUMDB=off # 关闭校验 Go 依赖包的哈希值
EOF
  source $HOME/.bashrc

  # 4. 初始化 Go 工作区
  mkdir -p $GOPATH && cd $GOPATH || return 1
  go work init

  skt::log::info "install go compile tool successfully"
}

# 安装 Protobuf
function skt::install::protobuf() {
  # 检查 protoc、protoc-gen-go 是否安装
  command -v protoc &>/dev/null && command -v protoc-gen-go &>/dev/null && return 0

  rm -rf /tmp/protobuf # clean up

  # 1. 安装 protobuf
  cd /tmp/ || return 1
  git clone -b v3.21.1 --depth=1 https://github.com/protocolbuffers/protobuf
  cd protobuf || return 1
  libtoolize --automake --copy --debug --force
  ./autogen.sh
  ./configure
  make
  sudo make install
  skt::common::sudo "make install"
  protoc --version | grep -q 'libprotoc 3.21.1' || {
    skt::log::error "protoc version is not '3.21.1', maynot install protobuf properly"
    return 1
  }

  skt::log::info "install protoc tool successfully"

  # 2. 安装 protoc-gen-go
  go install github.com/golang/protobuf/protoc-gen-go@v1.5.2

  skt::log::info "install protoc-gen-go plugin successfully"
}

# 安装 vim-go IDE
function skt::install::vim_ide() {
  rm -rf $HOME/.vim $HOME/.vimrc /tmp/gotools-for-vim.tgz # clean up

  mkdir -p ~/.vim/pack/plugins/start
  git clone --depth=1 https://github.com/fatih/vim-go.git $HOME/.vim/pack/plugins/start/vim-go
  cp "${SKT_ROOT}/scripts/install/vimrc" $HOME/.vimrc

  source $HOME/.bashrc
  skt::log::info "install vim ide successfully"
}

# 安装 Go 环境
function skt::install::go() {
  skt::install::go_command || return 1
  skt::install::protobuf || return 1

  skt::log::info "install go develop environment successfully"
}

# 初始化新申请的 Linux 服务器，使其成为一个友好的开发机
function skt::install::init_into_go_env() {
  # 1. Linux 服务器基本配置
  skt::install::prepare_linux || return 1

  # 2. Go 编译环境安装和配置
  skt::install::go || return 1

  # 3. Go 开发 IDE 安装和配置
  skt::install::vim_ide || return 1

  skt::log::info "initialize linux to go development machine  successfully"
}

# 自动配置环境并安装应用
function skt::install::install() {
  # 1. 配置 Linux 使其成为一个友好的 Go 开发机
  skt::install::init_into_go_env || return 1

  #  # 2. 安装 skeleton 应用
  #  skt::install::install_skt || return 1
  #
  #  # 3. 测试安装后的 skt 系统功能是否正常
  #  skt::test::test || return 1

  skt::log::info "$(echo -e '\033[32mcongratulations, install skeleton application successfully!\033[0m')"
}
