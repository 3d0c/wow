package cache

type Cache interface {
	Add(key int, expiration int64) error
	Get(key int) (bool, error)
	Delete(key int)
}
