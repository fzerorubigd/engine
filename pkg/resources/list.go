package resources

import "sync"

var (
	all  map[string]string
	lock sync.RWMutex
)

// RegisterResource try to register a new method for new resource
func RegisterResource(method, resource string) {
	lock.Lock()
	defer lock.Unlock()

	all[method] = resource
}

// QueryResource get a resource
func QueryResource(method string) (string, bool) {
	lock.RLock()
	defer lock.RUnlock()
	r, ok := all[method]
	return r, ok
}
