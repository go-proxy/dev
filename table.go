package proxy

import "sync"

//路由表
type table struct {
	mu   sync.RWMutex
	data map[string]string
}

func newTable() *table {
	return &table{
		data: make(map[string]string),
	}
}

//设置数据
func (t *table) Set(key string, val string) {
	t.mu.Lock()
	t.data[key] = val
	t.mu.Unlock()
}

//获取数据
func (t *table) Get(key string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data[key]
}

//获取数据
func (t *table) GetAll() map[string]string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data
}

//删除
func (t *table) DelAll() {
	t.mu.Lock()
	t.data = make(map[string]string)
	t.mu.Unlock()
}

//长度
func (t *table) Len() int {
	return len(t.data)
}

//Key是否存在
func (t *table) Exists(key string) bool {
	_, ok := t.data[key]
	return ok
}
