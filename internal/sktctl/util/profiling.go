// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package util

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/spf13/pflag"
)

var (
	profileName   string
	profileOutput string
)

func AddProfilingFlags(flags *pflag.FlagSet) {
	flags.StringVar(
		&profileName,
		"profile",
		"none",
		"Name of profile to capture. One of (none|cpu|heap|goroutine|threadcreate|block|mutex)",
	)
	flags.StringVar(
		&profileOutput,
		"profile-output",
		"profile.pprof",
		"Name of the file to write the profile to",
	)
}

func InitProfiling() error {
	switch profileName {
	case "none":
		return nil
	case "cpu":
		f, err := os.Create(profileOutput)
		if err != nil {
			return err
		}
		return pprof.StartCPUProfile(f)
	case "block":
		runtime.SetBlockProfileRate(1)
		return nil
	case "mutex":
		runtime.SetMutexProfileFraction(1)
		return nil
	default:
		if profile := pprof.Lookup(profileName); profile == nil {
			return fmt.Errorf("unknown profile '%s'", profileName)
		}
	}
	return nil
}

func FlushProfiling() error {
	switch profileName {
	case "none":
		return nil
	case "cpu":
		pprof.StopCPUProfile()
	case "heap":
		runtime.GC()
		fallthrough
	default:
		profile := pprof.Lookup(profileName)
		if profile == nil {
			return nil
		}

		f, err := os.Create(profileOutput)
		if err != nil {
			return err
		}
		_ = profile.WriteTo(f, 0)
	}
	return nil
}
