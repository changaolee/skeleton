## 目录规范

本项目遵循 [project-layout](https://github.com/golang-standards/project-layout) 目录规范。

跟 project-layout 目录规范唯一不一样的地方是，skeleton 将具体实现目录 `apiserver` 放在 `internal/` 目录下，而非 `internal/app/` 目录下，例如：

```bash
$ ls internal/         
apiserver  pkg
```

这样做既可以保证 `internal` 目录下的文件功能清晰、整齐，又能缩短引用路径。
