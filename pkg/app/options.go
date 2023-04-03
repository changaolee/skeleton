package app

// CliOptions 定义了从命令行读取参数的配置选项接口.
type CliOptions interface {
	Flags() (fss NamedFlagSets) // 通过命令行参数解析出的 FlagSets
	Validate() error            // 用于校验参数是否合法
}

// CompletableOptions 定义了完整选项需要实现的接口.
type CompletableOptions interface {
	Complete() error
}

// PrintableOptions 定义了可被打印的选项需要实现的接口.
type PrintableOptions interface {
	String() string
}
