package posixsignal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/changaolee/skeleton/pkg/shutdown"
)

const Name = "PosixSignalManager"

type Manager struct {
	signals []os.Signal
}

// NewPosixSignalManager 初始化 POSIX 信号 shutdown 监听实例.
func NewPosixSignalManager(signals ...os.Signal) *Manager {
	if len(signals) == 0 {
		signals = make([]os.Signal, 2)
		signals[0] = os.Interrupt
		signals[1] = syscall.SIGTERM
	}
	return &Manager{
		signals: signals,
	}
}

func (m *Manager) GetName() string {
	return Name
}

func (m *Manager) Start(gs shutdown.GSInterface) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, m.signals...)

		// 阻塞，等待 shutdown 信号
		<-c

		gs.StartShutdown(m)
	}()

	return nil
}

func (m *Manager) ShutdownStart() error {
	return nil
}

func (m *Manager) ShutdownFinish() error {
	os.Exit(0)

	return nil
}
