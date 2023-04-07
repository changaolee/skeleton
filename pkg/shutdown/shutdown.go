// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package shutdown

import "sync"

type Callback interface {
	OnShutdown(manager string) error
}

// CallbackFunc 是一个辅助函数，用于自定义回调函数.
type CallbackFunc func(string) error

// OnShutdown 是在 shutdown 触发时执行的方法.
func (f CallbackFunc) OnShutdown(manager string) error {
	return f(manager)
}

type GSInterface interface {
	AddManager(manager Manager)                // 添加 Shutdown Manager
	AddCallback(callback Callback)             // 添加 Shutdown Callback
	SetErrorHandler(errorHandler ErrorHandler) // 设置 ErrorHandler

	StartShutdown(manager Manager) // 执行指定 Manager 的 shutdown
	ReportError(err error)         // 向 ErrorHandler 报错
}

type Manager interface {
	GetName() string
	Start(gs GSInterface) error
	ShutdownStart() error
	ShutdownFinish() error
}

type ErrorHandler interface {
	OnError(err error)
}

type GracefulShutdown struct {
	callbacks    []Callback
	managers     []Manager
	errorHandler ErrorHandler
}

// New 初始化优雅关闭实例.
func New() *GracefulShutdown {
	return &GracefulShutdown{
		callbacks: make([]Callback, 0, 10),
		managers:  make([]Manager, 0, 3),
	}
}

// Start 在所有已添加的 Manager 上启动监听.
func (gs *GracefulShutdown) Start() error {
	for _, manager := range gs.managers {
		if err := manager.Start(gs); err != nil {
			return err
		}
	}
	return nil
}

// AddManager 添加 Manager 用于监听 shutdown 请求.
func (gs *GracefulShutdown) AddManager(manager Manager) {
	gs.managers = append(gs.managers, manager)
}

// AddCallback 添加 Callback 以便在 shutdown 时调用.
func (gs *GracefulShutdown) AddCallback(callback Callback) {
	gs.callbacks = append(gs.callbacks, callback)
}

// SetErrorHandler 设置 ErrorHandler 用于在 Manager 或 Callback 失败时调用.
func (gs *GracefulShutdown) SetErrorHandler(errorHandler ErrorHandler) {
	gs.errorHandler = errorHandler
}

// StartShutdown 用于执行指定 Manager 的 shutdown.
func (gs *GracefulShutdown) StartShutdown(manager Manager) {
	gs.ReportError(manager.ShutdownStart())

	var wg sync.WaitGroup
	for _, callback := range gs.callbacks {
		wg.Add(1)
		go func(callback Callback) {
			defer wg.Done()

			gs.ReportError(callback.OnShutdown(manager.GetName()))
		}(callback)
	}
	wg.Wait()

	gs.ReportError(manager.ShutdownFinish())
}

// ReportError 用于向 ErrorHandler 报错.
func (gs *GracefulShutdown) ReportError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.OnError(err)
	}
}
