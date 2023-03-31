#!/usr/bin/env bash

# The root of the build/dist directory
SKT_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
source "${SKT_ROOT}/scripts/install/common.sh"

source "${SKT_ROOT}/scripts/install/mariadb.sh"
source "${SKT_ROOT}/scripts/install/redis.sh"
source "${SKT_ROOT}/scripts/install/mongodb.sh"

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

# 安装数据库
function skt::install::install_storage() {
  skt::mariadb::install || return 1
  skt::redis::install || return 1
  skt::mongodb::install || return 1

  skt::log::info "install storage successfully"
}

# 如果是通过脚本安装，需要先尝试获取安装脚本指定的 Tag，Tag 记录在 version 文件中
function skt::install::obtain_branch_flag() {
  if [ -f "${SKT_ROOT}"/version ]; then
    echo $(cat "${SKT_ROOT}"/version)
  fi
}

# 安装 CFSSL
function skt::install::install_cfssl() {
  mkdir -p $HOME/bin/
  wget https://github.com/cloudflare/cfssl/releases/download/v1.6.1/cfssl_1.6.1_linux_amd64 -O $HOME/bin/cfssl
  wget https://github.com/cloudflare/cfssl/releases/download/v1.6.1/cfssljson_1.6.1_linux_amd64 -O $HOME/bin/cfssljson
  wget https://github.com/cloudflare/cfssl/releases/download/v1.6.1/cfssl-certinfo_1.6.1_linux_amd64 -O $HOME/bin/cfssl-certinfo
  #wget https://pkg.cfssl.org/R1.2/cfssl_linux-amd64 -O $HOME/bin/cfssl
  #wget https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64 -O $HOME/bin/cfssljson
  #wget https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64 -O $HOME/bin/cfssl-certinfo
  chmod +x $HOME/bin/{cfssl,cfssljson,cfssl-certinfo}

  skt::log::info "install cfssl tools successfully"
}

# 准备 skeleton 安装环境
function skt::install::prepare_skeleton() {
  rm -rf $WORKSPACE/golang/src/github.com/changaolee/skeleton # clean up

  # 1. 下载 skeleton 项目代码，先强制删除 skeleton 目录，确保 skeleton 源码都是最新的指定版本
  mkdir -p $WORKSPACE/golang/src/github.com/changaolee && cd $WORKSPACE/golang/src/github.com/changaolee
  git clone -b $(skt::install::obtain_branch_flag) --depth=1 https://github.com/changaolee/skeleton
  go work use ./skeleton

  # 注意：因为切换编译路径，所以这里要重新赋值 SKT_ROOT 和 LOCAL_OUTPUT_ROOT
  SKT_ROOT=$WORKSPACE/golang/src/github.com/changaolee/skeleton
  LOCAL_OUTPUT_ROOT="${SKT_ROOT}/${OUT_DIR:-_output}"

  pushd ${SKT_ROOT}

  # 2. 配置 $HOME/.bashrc 添加一些便捷入口
  if ! grep -q 'Alias for quick access' $HOME/.bashrc; then
    cat <<'EOF' >>$HOME/.bashrc
# Alias for quick access
export GOSRC="$WORKSPACE/golang/src"
export SKT_ROOT="$GOSRC/github.com/changaolee/skeleton"
alias ca="cd $GOSRC/github.com/changaolee"
alias skt="cd $GOSRC/github.com/changaolee/skeleton"
EOF
  fi

  # 3. 初始化 MariaDB 数据库，创建 skt 数据库

  # 3.1 登录数据库并创建 skt 用户
  mysql -h127.0.0.1 -P3306 -u"${MARIADB_ADMIN_USERNAME}" -p"${MARIADB_ADMIN_PASSWORD}" <<EOF
grant all on ${MARIADB_DATABASE}.* TO ${MARIADB_USERNAME}@127.0.0.1 identified by "${MARIADB_PASSWORD}";
flush privileges;
EOF

  # 3.2 用 skt 用户登录 mysql，执行 skeleton.sql 文件，创建 skt 数据库
  mysql -h127.0.0.1 -P3306 -u${MARIADB_USERNAME} -p"${MARIADB_PASSWORD}" <<EOF
source configs/skeleton.sql;
show databases;
EOF

  # 4. 创建必要的目录
  echo ${LINUX_PASSWORD} | sudo -S mkdir -p ${SKT_DATA_DIR}/{skt-apiserver,skt-authz-server,skt-pump,skt-watcher}
  skt::common::sudo "mkdir -p ${SKT_INSTALL_DIR}/bin"
  skt::common::sudo "mkdir -p ${SKT_CONFIG_DIR}/cert"
  skt::common::sudo "mkdir -p ${SKT_LOG_DIR}"

  # 5. 安装 cfssl 工具集
  ! command -v cfssl &>/dev/null || ! command -v cfssl-certinfo &>/dev/null || ! command -v cfssljson &>/dev/null && {
    skt::install::install_cfssl || return 1
  }

  # 6. 配置 hosts
  if ! egrep -q 'skt.*lichangao.com' /etc/hosts; then
    echo ${LINUX_PASSWORD} | sudo -S bash -c "cat << 'EOF' >> /etc/hosts
    127.0.0.1 skt.api.lichangao.com
    127.0.0.1 skt.authz.lichangao.com
    EOF"
  fi

  skt::log::info "prepare for skeleton installation successfully"
  popd
}

# 安装 Skeleton 应用
function skt::install::install_skeleton() {
  # 1. 安装并初始化数据库
  skt::install::install_storage || return 1

  # 2. 先准备安装环境
  skt::install::prepare_skeleton || return 1

  #  # 3. 安装 skt-apiserver 服务
  #  skt::apiserver::install || return 1
  #
  #  # 4. 安装 sktctl 客户端工具
  #  skt::sktctl::install || return 1
  #
  #  # 5. 安装 skt-authz-server 服务
  #  skt::authzserver::install || return 1
  #
  #  # 6. 安装 skt-pump 服务
  #  skt::pump::install || return 1
  #
  #  # 7. 安装 skt-watcher 服务
  #  skt::watcher::install || return 1
  #
  #  # 8. 安装 man page
  #  skt::man::install || return 1

  skt::log::info "install skt application successfully"
}

# 自动配置环境并安装应用
function skt::install::install() {
  # 1. 配置 Linux 使其成为一个友好的 Go 开发机
  skt::install::init_into_go_env || return 1

  # 2. 安装 skeleton 应用
  skt::install::install_skeleton || return 1

  #  # 3. 测试安装后的 skt 系统功能是否正常
  #  skt::test::test || return 1

  skt::log::info "$(echo -e '\033[32mcongratulations, install skeleton application successfully!\033[0m')"
}

eval $*
