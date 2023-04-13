package code

//go:generate codegen -type=int

// skt-apiserver: user errors.
const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 110001

	// ErrUserAlreadyExist - 400: User already exist.
	ErrUserAlreadyExist
)

// skt-apiserver: secret errors.
const (
	// ErrReachMaxCount - 400: Secret reach the max count.
	ErrReachMaxCount int = iota + 110101

	// ErrSecretNotFound - 404: Secret not found.
	ErrSecretNotFound
)

// skt-apiserver: policy errors.
const (
	// ErrPolicyNotFound - 404: Policy not found.
	ErrPolicyNotFound int = iota + 110201
)
