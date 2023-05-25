// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package info

import (
	"fmt"
	"reflect"
	"strconv"

	hoststat "github.com/likexian/host-stat-go"
	"github.com/spf13/cobra"

	"github.com/changaolee/skeleton/internal/pkg/clioptions"

	"github.com/changaolee/skeleton/internal/sktctl/util"
	"github.com/changaolee/skeleton/internal/sktctl/util/templates"
	"github.com/changaolee/skeleton/pkg/util/iputil"
)

type Info struct {
	HostName  string
	IPAddress string
	OSRelease string
	CPUCore   uint64
	MemTotal  string
	MemFree   string
}

type InfoOptions struct {
	clioptions.IOStreams
}

var infoExample = templates.Examples(`
		# Print the host information
		sktctl info`)

func NewInfoOptions(ioStreams clioptions.IOStreams) *InfoOptions {
	return &InfoOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdInfo(f util.Factory, ioStreams clioptions.IOStreams) *cobra.Command {
	o := NewInfoOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "info",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Print the host information",
		Long:                  "Print the host information.",
		Example:               infoExample,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	return cmd
}

func (o *InfoOptions) Run(args []string) error {
	var info Info

	hostInfo, err := hoststat.GetHostInfo()
	if err != nil {
		return fmt.Errorf("get host info failed! error: %w", err)
	}

	info.HostName = hostInfo.HostName
	info.OSRelease = hostInfo.Release + " " + hostInfo.OSBit

	memStat, err := hoststat.GetMemStat()
	if err != nil {
		return fmt.Errorf("get mem stat failed! error: %w", err)
	}

	info.MemTotal = strconv.FormatUint(memStat.MemTotal, 10) + "M"
	info.MemFree = strconv.FormatUint(memStat.MemFree, 10) + "M"
	info.IPAddress = iputil.GetLocalIP()

	cpuStat, err := hoststat.GetCPUInfo()
	if err != nil {
		return fmt.Errorf("get cpu stat failed! error: %w", err)
	}

	info.CPUCore = cpuStat.CoreCount

	s := reflect.ValueOf(&info).Elem()
	typeOfInfo := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		v := fmt.Sprintf("%v", f.Interface())
		if v != "" {
			_, _ = fmt.Fprintf(o.Out, "%12s %v\n", typeOfInfo.Field(i).Name+":", f.Interface())
		}
	}

	return nil
}
