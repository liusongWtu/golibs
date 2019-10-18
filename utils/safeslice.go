package utils

import (
	"bytes"
	"encoding/gob"
	"github.com/json-iterator/go"
	"sync"
)

// NewSafeSlice return new safemap
func NewSafeSlice() *SafeSlice {
	return &SafeSlice{
		lock: new(sync.RWMutex),
		sm:   make([]interface{},0),
	}
}

// SafeSlice is a map with lock
type SafeSlice struct {
	lock *sync.RWMutex
	sm   []interface{}
}

// Set slice
func (m *SafeSlice) Set(v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.sm = append(m.sm,v)
	return true
}

// Clear the given key and value.
func (m *SafeSlice) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.sm = make([]interface{},0)
}

// Items returns all items in safe slice.
func (m *SafeSlice) Items() []interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make([]interface{},0)
	for _, v := range m.sm {
		r = append(r,v)
	}
	return r
}

// Items returns all items in safe slice.
func (m *SafeSlice) ItemsAndClear() []interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make([]interface{},0)
	for _, v := range m.sm {
		r = append(r,v)
	}
	m.sm = make([]interface{},0)
	return r
}

// GetItems returns all items in safe slice with json copy.
func (m *SafeSlice) GetItemsWithJsonCopy() []interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var r []interface{}
	m.jsonCopy(&r, m.sm)
	return r
}

// GetItems returns all items in safe slice with gob copy.
func (m *SafeSlice) GetItemsWithGobCopy() []interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var r []interface{}
	m.gobCopy(&r, m.sm)
	return r
}

// Count returns the number of items within the map.
func (m *SafeSlice) Count() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.sm)
}

// 返回键列表
func (m *SafeSlice) Keys() []interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	keys := make([]interface{}, 0)
	for key := range m.sm {
		keys = append(keys, key)
	}
	return keys
}

// deep copy with json
func (m *SafeSlice) jsonCopy(dst, src interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}

// deep copy with gob
func (m *SafeSlice) gobCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(src)
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
