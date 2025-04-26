package limiter

type StoreStrategy interface {
	AllowRequest(key string, limit int, windowSec int) (bool, error)
	BlockDurationExceeded(key string) (bool, error)
	SetBlock(key string, durationSec int) error
}