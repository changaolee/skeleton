// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package app

import (
	"bytes"
	goflag "flag"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"

	"github.com/changaolee/skeleton/pkg/log"
)

// WordSepNormalizeFunc 规范化标志中的 "_" 和 "-" 分隔符.
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	}
	return pflag.NormalizedName(name)
}

// WarnWordSepNormalizeFunc 替换 "_" 分隔符并打印警告.
func WarnWordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		nname := strings.ReplaceAll(name, "_", "-")
		log.Warnf("%s is DEPRECATED and will be removed in a future version. Use %s instead.", name, nname)

		return pflag.NormalizedName(nname)
	}
	return pflag.NormalizedName(name)
}

// InitFlags 规范化、解析并记录命令行标志.
func InitFlags(flags *pflag.FlagSet) {
	flags.SetNormalizeFunc(WordSepNormalizeFunc)
	flags.AddGoFlagSet(goflag.CommandLine) // 兼容 golang flag
}

// PrintFlags 打印指定 FlagSet 中的所有标记.
func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		log.Debugf("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}

// NamedFlagSets 按照调用 FlagSet 的顺序来存储.
type NamedFlagSets struct {
	// Order 是 FlagSet 名称的有序列表
	Order []string
	// FlagSets 按名称存储了对应的 FlagSet
	FlagSets map[string]*pflag.FlagSet
}

// FlagSet 返回指定名称的 FlagSet，并将其添加到已排序的名称列表中（如果它还不在列表中）.
func (nfs *NamedFlagSets) FlagSet(name string) *pflag.FlagSet {
	if nfs.FlagSets == nil {
		nfs.FlagSets = map[string]*pflag.FlagSet{}
	}
	if _, ok := nfs.FlagSets[name]; !ok {
		nfs.FlagSets[name] = pflag.NewFlagSet(name, pflag.ExitOnError)
		nfs.Order = append(nfs.Order, name)
	}
	return nfs.FlagSets[name]
}

// PrintSections 按 cols 作为最大值打印给定名称的 FlagSets
// 如果 cols 为 0，则不换行打印.
func PrintSections(w io.Writer, fss NamedFlagSets, cols int) {
	for _, name := range fss.Order {
		fs := fss.FlagSets[name]
		if !fs.HasFlags() {
			continue
		}

		wideFS := pflag.NewFlagSet("", pflag.ExitOnError)
		wideFS.AddFlagSet(fs)

		var zzz string
		if cols > 24 {
			zzz = strings.Repeat("z", cols-24)
			wideFS.Int(zzz, 0, strings.Repeat("z", cols-24))
		}

		var buf bytes.Buffer
		_, _ = fmt.Fprintf(
			&buf,
			"\n%s flags:\n\n%s",
			strings.ToUpper(name[:1])+name[1:],
			wideFS.FlagUsagesWrapped(cols),
		)

		if cols > 24 {
			i := strings.Index(buf.String(), zzz)
			lines := strings.Split(buf.String()[:i], "\n")
			_, _ = fmt.Fprint(w, strings.Join(lines[:len(lines)-1], "\n"))
			_, _ = fmt.Fprintln(w)
		} else {
			_, _ = fmt.Fprint(w, buf.String())
		}
	}
}
