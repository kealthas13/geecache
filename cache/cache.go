package cache

// 接口声明
type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
}
