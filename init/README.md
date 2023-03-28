# Systemd 配置、安装和启动

- [Systemd 配置、安装和启动](#systemd-配置安装和启动)
    - [1. 前置操作](#前置操作)
    - [2. 创建 skeleton systemd unit 模板文件](#创建-skeleton-systemd-unit-模板文件)
    - [3. 复制 systemd unit 模板文件到 sysmted 配置目录](#复制-systemd-unit-模板文件到-sysmted-配置目录)
    - [4. 启动 systemd 服务](#启动-systemd-服务)

## 1. 前置操作

1. 创建需要的目录

```bash
sudo mkdir -p /data/skeleton /opt/skeleton/bin /etc/skeleton /var/log/skeleton
```

2. 编译构建 `skeleton` 二进制文件

```bash
make build # 编译源码生成 skeleton 二进制文件
```

3. 将 `skeleton` 可执行文件安装在 `bin` 目录下

```bash
sudo cp _output/platforms/linux/amd64/skeleton /opt/skeleton/bin # 安装二进制文件
```

4. 安装 `skeleton` 配置文件

```bash
sed 's/.\/_output/\/etc\/skeleton/g' configs/skeleton.yaml > skeleton.sed.yaml # 替换 CA 文件路径
sudo mv skeleton.sed.yaml /etc/skeleton/ # 安装配置文件
```

5. 安装 CA 文件

```bash
make ca # 创建 CA 文件
sudo cp -a _output/cert/ /etc/skeleton/ # 将 CA 文件复制到 skeleton 配置文件目录
```

## 2. 创建 skeleton systemd unit 模板文件

执行如下 shell 脚本生成 `skeleton.service.template`

```bash
cat > skeleton.service.template <<EOF
[Unit]
Description=APIServer skeleton.
Documentation=https://github.com/changaolee/skeleton/blob/master/init/README.md
[Service]
WorkingDirectory=/data/skeleton
ExecStartPre=/usr/bin/mkdir -p /data/skeleton
ExecStartPre=/usr/bin/mkdir -p /var/log/skeleton
ExecStart=/opt/skeleton/bin/skeleton --config=/etc/skeleton/skeleton.yaml
Restart=always
RestartSec=5
StartLimitInterval=0
[Install]
WantedBy=multi-user.target
EOF
```

## 3. 复制 systemd unit 模板文件到 sysmted 配置目录

```bash
sudo cp skeleton.service.template /etc/systemd/system/skeleton.service
```

## 4. 启动 systemd 服务

```bash
sudo systemctl daemon-reload && systemctl enable skeleton && systemctl restart skeleton
```