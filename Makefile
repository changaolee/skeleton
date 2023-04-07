# ==============================================================================
# 定义 Makefile all 伪目标，执行 `make` 时，会默认会执行 all 伪目标.

.DEFAULT_GOAL := all

.PHONY: all
all: tidy gen add-copyright format lint cover build

# ==============================================================================
# 定义包名

ROOT_PACKAGE=github.com/changaolee/skeleton
VERSION_PACKAGE=github.com/changaolee/skeleton/pkg/version

# ==============================================================================
# Includes

# 确保 `include common.mk` 位于第一行，common.mk 中定义了一些变量，后面的子 makefile 有依赖.
include scripts/make-rules/common.mk
include scripts/make-rules/copyright.mk
include scripts/make-rules/generate.mk
include scripts/make-rules/golang.mk
include scripts/make-rules/tools.mk

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  BINS             The binaries to build. Default is all of cmd.
                   This option is available when using: `make build` or `make build.multiarch`
                   Example: `make build BINS="skt-apiserver skt-authz-server"`
  PLATFORMS        The multiple platforms to build. Default is linux_amd64 and linux_arm64.
                   This option is available when using: `make build.multiarch` or `make image.multiarch` or `make push.multiarch`
                   Example: `make image.multiarch IMAGES="iam-apiserver" PLATFORMS="linux_amd64 linux_arm64"`
  VERSION          The version information compiled into binaries.
                   The default is obtained from gsemver or git.
  V                Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

# ==============================================================================
# Targets

## build: 为主机所在平台构建源代码.
.PHONY: build
build:
	@$(MAKE) go.build

## build.multiarch: 为多个平台构建源代码，参考选项 PLATFORMS.
.PHONY: build.multiarch
build.multiarch:
	@$(MAKE) go.build.multiarch

## image: 为主机所在平台构建 docker 镜像.
.PHONY: image
image:
	@$(MAKE) image.build

## image.multiarch: 为多个平台构建 docker 镜像，参考选项 PLATFORMS.
.PHONY: image.multiarch
image.multiarch:
	@$(MAKE) image.build.multiarch

## push: 为主机所在平台构建 docker 镜像并推送到仓库.
.PHONY: push
push:
	@$(MAKE) image.push

## push.multiarch: 为多个平台构建 docker 镜像并推送到仓库，参考选项 PLATFORMS.
.PHONY: push.multiarch
push.multiarch:
	@$(MAKE) image.push.multiarch

## deploy: 将更新的组件部署到开发环境.
.PHONY: deploy
deploy:
	@$(MAKE) deploy.run

## clean: 移除所有构建过程产生的文件.
.PHONY: clean
clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

## lint: 检查 GO 源代码的语法和风格.
.PHONY: lint
lint:
	@$(MAKE) go.lint

## test: 运行单元测试.
.PHONY: test
test:
	@$(MAKE) go.test

## cover: 运行单元测试并获取测试覆盖率.
.PHONY: cover
cover:
	@$(MAKE) go.test.cover

## release.build: 构建项目.
.PHONY: release.build
release.build:
	@$(MAKE) push.multiarch

## release: 发布项目.
.PHONY: release
release:
	@$(MAKE) release.run

## format: 格式化源代码.
.PHONY: format
format:
	@$(MAKE) go.format

## verify-copyright: 验证所有文件的版权声明.
.PHONY: verify-copyright
verify-copyright:
	@$(MAKE) copyright.verify

## add-copyright: 确保所有文件都有版权声明.
.PHONY: add-copyright
add-copyright:
	@$(MAKE) copyright.add

## gen: 生成所有必要的文件.
.PHONY: gen
gen:
	@$(MAKE) gen.run

## ca: 为所有组件生成 CA 证书.
.PHONY: ca
ca:
	@$(MAKE) gen.ca

## swagger: 生成 Swagger 文档.
.PHONY: swagger
swagger:
	@$(MAKE) swagger.run

## serve-swagger: 开启 Swagger 文档服务.
.PHONY: swagger.serve
serve-swagger:
	@$(MAKE) swagger.serve

## dependencies: 安装必要的依赖.
.PHONY: dependencies
dependencies:
	@$(MAKE) dependencies.run

## tools: 安装依赖工具.
.PHONY: tools
tools:
	@$(MAKE) tools.install

## check-updates: 检查项目的过时依赖.
.PHONY: check-updates
check-updates:
	@$(MAKE) go.updates

## tidy: 整理 Go 模块依赖.
.PHONY: tidy
tidy:
	@$(MAKE) go.tidy

## help: 展示帮助信息.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
