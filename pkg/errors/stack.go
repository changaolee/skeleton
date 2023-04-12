package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// Frame 表示一个栈帧中的程序计数器.
type Frame uintptr

// pc 返回当前帧的程序计数器.
func (f Frame) pc() uintptr { return uintptr(f) - 1 }

// file 返回包含此帧的 pc 函数的文件的完整路径.
func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line 返回此帧的 pc 函数的源代码行号。
func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// name 返回函数名.
func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// Format 根据 fmt.Formatter 接口对帧进行格式化.
//
//    %s    源文件.
//    %d    源代码行号.
//    %n    函数名.
//    %v    相当于 %s:%d.
//    %+s   源文件相对于编译时 GOPATH 的函数名和路径，由 \n\t 分隔（<funcname>\n\t<path>）.
//    %+v   相当于 %+s:%d.
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			_, _ = io.WriteString(s, f.name())
			_, _ = io.WriteString(s, "\n\t")
			_, _ = io.WriteString(s, f.file())
		default:
			_, _ = io.WriteString(s, path.Base(f.file()))
		}
	case 'd':
		_, _ = io.WriteString(s, strconv.Itoa(f.line()))
	case 'n':
		_, _ = io.WriteString(s, funcname(f.name()))
	case 'v':
		f.Format(s, 's')
		_, _ = io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

// MarshalText 将堆栈格式化为文本字符串.
// 输出与 fmt.Sprintf("%+v", f) 相同，但没有换行符或制表符.
func (f Frame) MarshalText() ([]byte, error) {
	name := f.name()
	if name == "unknown" {
		return []byte(name), nil
	}
	return []byte(fmt.Sprintf("%s %s:%d", name, f.file(), f.line())), nil
}

// StackTrace 是一个从内层到外层的栈帧列表.
type StackTrace []Frame

// Format 根据 fmt.Formatter 接口对栈帧进行格式化.
//
//    %s	列出堆栈中每个帧的源文件.
//    %v	列出堆栈中每个帧的源文件和行号.
//    %+v   打印堆栈中每个帧的文件名、函数和行号.
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, f := range st {
				_, _ = io.WriteString(s, "\n")
				f.Format(s, verb)
			}
		case s.Flag('#'):
			_, _ = fmt.Fprintf(s, "%#v", []Frame(st))
		default:
			st.formatSlice(s, verb)
		}
	case 's':
		st.formatSlice(s, verb)
	}
}

// formatSlice 会将此 StackTrace 格式化为帧的切片到给定的缓冲区中.
// 仅在使用 "%s" 或 "%v" 调用时有效.
func (st StackTrace) formatSlice(s fmt.State, verb rune) {
	_, _ = io.WriteString(s, "[")
	for i, f := range st {
		if i > 0 {
			_, _ = io.WriteString(s, " ")
		}
		f.Format(s, verb)
	}
	_, _ = io.WriteString(s, "]")
}

// stack 表示包含程序计数器的栈.
type stack []uintptr

func (s *stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				f := Frame(pc)
				_, _ = fmt.Fprintf(st, "\n%+v", f)
			}
		}
	}
}

func (s *stack) StackTrace() StackTrace {
	f := make([]Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = Frame((*s)[i])
	}
	return f
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

// funcname 删除文件路径前缀.
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
