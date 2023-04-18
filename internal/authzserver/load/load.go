package load

import (
	"context"
	"sync"
	"time"

	"github.com/changaolee/skeleton/internal/authzserver/cache"
	"github.com/changaolee/skeleton/pkg/log"
)

// Loader 定义了重载 secrets 和 policies 的方法.
type Loader interface {
	Reload() error
}

// Load 用于重载 secrets 和 policies.
type Load struct {
	ctx    context.Context
	lock   *sync.RWMutex
	loader Loader
}

// NewLoader 创建一个 Load.
func NewLoader(ctx context.Context, loader Loader) *Load {
	return &Load{
		ctx:    ctx,
		lock:   new(sync.RWMutex),
		loader: loader,
	}
}

// reloadQueue 用于通知需要重载 secrets 和 policies 的 channel.
var reloadQueue = make(chan func())

// 在下一次重载时需要执行的函数列表.
var requeue []func()

// 用于保护 requeue 的并发安全.
var requeueLock sync.Mutex

// Start 启动 Load 服务.
func (l *Load) Start() {
	go l.startPubSubLoop()
	go l.reloadQueueLoop()
	go l.reloadLoop()

	l.DoReload()
}

// DoReload 进行 secrets 和 policies 同步.
func (l *Load) DoReload() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if err := l.loader.Reload(); err != nil {
		log.Errorf("Fail to refresh target storage: %s", err.Error())
	}

	log.Debugw("Refresh target storage success")
}

// startPubSubLoop 监听订阅事件，触发时会通知 channel.
func (l *Load) startPubSubLoop() {
	cacheIns, err := cache.GetRedisInstance(nil)
	if err != nil {
		log.Errorf("Connection to redis failed")
		return
	}
	_ = cacheIns.StartPubSubHandler(l.ctx, RedisPubSubChannel, func(v interface{}) {
		handleRedisEvent(v, nil, nil)
	})
}

// reloadQueueLoop 将 channel 中的消息缓存在 requeue 中.
func (l *Load) reloadQueueLoop() {
	for {
		select {
		case <-l.ctx.Done():
			return
		case fn := <-reloadQueue:
			requeueLock.Lock()
			requeue = append(requeue, fn)
			requeueLock.Unlock()
			log.Infow("Reload queued")
		}
	}
}

// reloadLoop 每秒检查 requeue 是否为空，不空则重载数据.
func (l *Load) reloadLoop() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-l.ctx.Done():
			return
		case <-ticker.C:
			callbacks, ok := shouldReload()
			if !ok {
				continue
			}
			start := time.Now()
			l.DoReload()
			for _, callback := range callbacks {
				if callback != nil {
					callback()
				}
			}
			log.Infof("Reload: cycle completed in %v", time.Since(start))
		}
	}
}

// shouldReload 判断是否需要重载.
func shouldReload() ([]func(), bool) {
	requeueLock.Lock()
	defer requeueLock.Unlock()
	if len(requeue) == 0 {
		return nil, false
	}
	callbacks := requeue
	requeue = []func(){}

	return callbacks, true
}
