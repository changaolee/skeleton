## Skeleton - Go 语言开发脚手架

Skeleton 是一个基于 Go 语言的开发脚手架，供开发者克隆后二次开发，快速构建自己的应用。

## Features

- 使用了简洁架构；
- 使用众多常用的 Go 包：gorm, casbin, govalidator, jwt-go, gin, cobra, viper, pflag, zap, pprof, grpc, protobuf 等；
- 规范的目录结构，使用 [project-layout](https://github.com/golang-standards/project-layout) 目录规范；
- 具备认证（JWT）和授权功能（casbin）；
- 独立设计的 log 包、error 包；
- 使用高质量的 Makefile 管理项目；
- 静态代码检查；
- 带有单元测试、性能测试、模糊测试、Mock 测试测试案例；
- 丰富的 Web 功能（调用链、优雅关停、中间件、跨域、异常恢复等）；
    - HTTP、HTTPS、gRPC 服务器实现；
    - JSON、Protobuf 数据交换格式实现；
- 项目遵循众多开发规范：代码规范、版本规范、接口规范、日志规范、错误规范、提交规范等；
- 访问 MySQL 编程实现；
- 实现的业务功能：用户管理、博客管理；
- RESTful API 设计规范；
- OpenAPI 3.0/Swagger 2.0 API 文档；

## Installation

```bash
$ cd /tmp && wget https://github.com/changaolee/skeleton/archive/refs/heads/main.zip -O skeleton.zip
$ unzip skeleton.zip && cd skeleton-main
$ bash ./scripts/install/install.sh skt::install::install
```

## Test

```bash
# 创建用户
$ sktctl user create foo Foo@2023 foo@test.com
```

## Documentation

- [开发手册](./docs/devel/zh-CN/README.md)

## Feedback

如果您有任何反馈，请通过 `changao.li.work@outlook.com` 与我联系。

### 开发规范

本项目遵循以下开发规范：[skeleton 项目开发规范](./docs/devel/zh-CN/conversions/README.md)。

## License

[MIT](./LICENSE)
