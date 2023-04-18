package store

import "github.com/ory/ladon"

// PolicyStore defines the policy storage interface.
type PolicyStore interface {
	List() (map[string][]*ladon.DefaultPolicy, error)
}
