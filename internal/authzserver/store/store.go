package store

// IStore 定义了 Store 层接口.
type IStore interface {
	Policies() PolicyStore
	Secrets() SecretStore
}

var ins IStore

// Store 获取 Store 实例.
func Store() IStore {
	return ins
}

// SetStore 设置 Store 实例.
func SetStore(storeIns IStore) {
	ins = storeIns
}
