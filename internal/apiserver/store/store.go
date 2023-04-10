package store

// IStore 定义了 Store 层接口.
type IStore interface {
	Users() UserStore
	Close() error
}
