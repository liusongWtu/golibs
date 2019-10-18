package utils

import (
	"bytes"
	"encoding/gob"
	"github.com/json-iterator/go"
	"sync"
)

// NewSafeMap return new safemap
func NewSafeMap() *SafeMap {
	return &SafeMap{
		lock: new(sync.RWMutex),
		sm:   make(map[interface{}]interface{}),
	}
}

// SafeMap is a map with lock
type SafeMap struct {
	lock *sync.RWMutex
	sm   map[interface{}]interface{}
}

// Get from maps return the k's value
func (m *SafeMap) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.sm[k]; ok {
		return val
	}
	return nil
}

// Get from maps return the k's clone value
func (m *SafeMap) GetItem(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.sm[k]; ok {
		var r interface{}
		m.jsonCopy(&r, m.sm[k])
		return r
	}
	return nil
}

// Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *SafeMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.sm[k] = v
	return true
}

// Returns true if k is exist in the map.
func (m *SafeMap) Check(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	_, ok := m.sm[k]
	return ok
}

// Delete the given key and value.
func (m *SafeMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.sm, k)
}

// Clear the given key and value.
func (m *SafeMap) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.sm = make(map[interface{}]interface{})
}

// Items returns all items in safemap.
func (m *SafeMap) Items() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make(map[interface{}]interface{})
	for k, v := range m.sm {
		r[k] = v
	}
	return r
}

// Items returns all items in safemap.
func (m *SafeMap) ItemsAndClear() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make(map[interface{}]interface{})
	for k, v := range m.sm {
		r[k] = v
	}
	m.sm = make(map[interface{}]interface{})
	return r
}

// GetItems returns all items in safemap with json copy.
func (m *SafeMap) GetItemsWithJsonCopy() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var r map[interface{}]interface{}
	m.jsonCopy(&r, m.sm)
	return r
}

// GetItems returns all items in safemap with gob copy.
func (m *SafeMap) GetItemsWithGobCopy() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var r map[interface{}]interface{}
	m.gobCopy(&r, m.sm)
	return r
}

// Count returns the number of items within the map.
func (m *SafeMap) Count() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.sm)
}

// 返回键列表
func (m *SafeMap) Keys() []interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	keys := make([]interface{}, 0)
	for key := range m.sm {
		keys = append(keys, key)
	}
	return keys
}

// deep copy with json
func (m *SafeMap) jsonCopy(dst, src interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}

// deep copy with gob
func (m *SafeMap) gobCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(src)
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
