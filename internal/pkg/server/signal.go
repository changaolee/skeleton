package server

import (
	"os"
	"os/signal"
)

var onlyOneSignalHandler = make(chan struct{})

var shutdownHandler chan os.Signal

// SetupSignalHandler 返回一个 stop 通道，应用程序在监听到该通道关闭后，执行资源清理.
func SetupSignalHandler() <-chan struct{} {
	// 避免保证该函数只能执行一次，再次执行会 panic
	close(onlyOneSignalHandler)

	shutdownHandler = make(chan os.Signal, 2)

	stop := make(chan struct{})

	// shutdownHandler 通道会监听 shutdownSignals
	signal.Notify(shutdownHandler, shutdownSignals...)

	// 收到一次 shutdownSignal，程序优雅关闭
	// 收到两次 shutdownSignal，程序强制关闭
	go func() {
		<-shutdownHandler
		close(stop)
		<-shutdownHandler
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
