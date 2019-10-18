package orm

import (
	"sync"
)

const (
	odCascade             = "cascade"
	odSetNULL             = "set_null"
	odSetDefault          = "set_default"
	odDoNothing           = "do_nothing"
	defaultStructTagName  = "orm"
	defaultStructTagDelim = ";"
)

var (
	modelCache = &_modelCache{
		cache:           make(map[string]*modelInfo),
		cacheByFullName: make(map[string]*modelInfo),
	}
)

// model info collection
type _modelCache struct {
	sync.RWMutex    // only used outsite for bootStrap
	orders          []string
	cache           map[string]*modelInfo
	cacheByFullName map[string]*modelInfo
	done            bool
}

// get all model info
func (mc *_modelCache) all() map[string]*modelInfo {
	m := make(map[string]*modelInfo, len(mc.cache))
	for k, v := range mc.cache {
		m[k] = v
	}
	return m
}

// get ordered model info
func (mc *_modelCache) allOrdered() []*modelInfo {
	m := make([]*modelInfo, 0, len(mc.orders))
	for _, table := range mc.orders {
		m = append(m, mc.cache[table])
	}
	return m
}

// get model info by table name
func (mc *_modelCache) get(table string) (mi *modelInfo, ok bool) {
	mi, ok = mc.cache[table]
	return
}

// get model info by full name
func (mc *_modelCache) getByFullName(name string) (mi *modelInfo, ok bool) {
	mi, ok = mc.cacheByFullName[name]
	return
}

// set model info to collection
func (mc *_modelCache) set(table string, mi *modelInfo) *modelInfo {
	mii := mc.cache[table]
	mc.cache[table] = mi
	mc.cacheByFullName[mi.fullName] = mi
	if mii == nil {
		mc.orders = append(mc.orders, table)
	}
	return mii
}

// clean all model info.
func (mc *_modelCache) clean() {
	mc.orders = make([]string, 0)
	mc.cache = make(map[string]*modelInfo)
	mc.cacheByFullName = make(map[string]*modelInfo)
	mc.done = false
}

// ResetModelCache Clean model cache. Then you can re-RegisterModel.
// Common use this api for test case.
func ResetModelCache() {
	modelCache.clean()
}
