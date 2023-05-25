package templates

import (
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/russross/blackfriday"
	"github.com/spf13/cobra"
)

const Indentation = `  `

// LongDesc 将命令的详细描述标准化.
func LongDesc(s string) string {
	if len(s) == 0 {
		return s
	}
	return normalizer{s}.heredoc().markdown().trim().string
}

// Examples 将命令的示例标准化.
func Examples(s string) string {
	if len(s) == 0 {
		return s
	}
	return normalizer{s}.trim().indent().string
}

// Normalize 对给定命令执行所有必需的标准化.
func Normalize(cmd *cobra.Command) *cobra.Command {
	if len(cmd.Long) > 0 {
		cmd.Long = LongDesc(cmd.Long)
	}
	if len(cmd.Example) > 0 {
		cmd.Example = Examples(cmd.Example)
	}
	return cmd
}

// NormalizeAll 在整个命令树中执行所有必需的标准化.
func NormalizeAll(cmd *cobra.Command) *cobra.Command {
	if cmd.HasSubCommands() {
		for _, subCmd := range cmd.Commands() {
			NormalizeAll(subCmd)
		}
	}
	Normalize(cmd)
	return cmd
}

type normalizer struct {
	string
}

func (s normalizer) markdown() normalizer {
	bytes := []byte(s.string)
	formatted := blackfriday.Markdown(
		bytes,
		&ASCIIRenderer{Indentation: Indentation},
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS,
	)
	s.string = string(formatted)
	return s
}

func (s normalizer) heredoc() normalizer {
	s.string = heredoc.Doc(s.string)
	return s
}

func (s normalizer) trim() normalizer {
	s.string = strings.TrimSpace(s.string)
	return s
}

func (s normalizer) indent() normalizer {
	indentedLines := []string{}
	for _, line := range strings.Split(s.string, "\n") {
		trimmed := strings.TrimSpace(line)
		indented := Indentation + trimmed
		indentedLines = append(indentedLines, indented)
	}
	s.string = strings.Join(indentedLines, "\n")
	return s
}
