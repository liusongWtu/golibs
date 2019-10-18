//
// Usage:
// import(
//   "libs/session"
// )
//
//	func init() {
//		globalSessions, _ = session.NewManager("memcache", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"127.0.0.1:11211"}``)
//		go globalSessions.GC()
//	}
//
package session

import (
	"github.com/bradfitz/gomemcache/memcache"
	"net/http"
	"strings"
	"sync"
)

var memcachepder = &MemcacheProvider{}
var client *memcache.Client

// MemcacheSessionStore memcache session store
type MemcacheSessionStore struct {
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime int64
}

// Set value in memcache session
func (rs *MemcacheSessionStore) Set(key, value interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values[key] = value
	return nil
}

// Get value in memcache session
func (rs *MemcacheSessionStore) Get(key interface{}) interface{} {
	rs.lock.RLock()
	defer rs.lock.RUnlock()
	if v, ok := rs.values[key]; ok {
		return v
	}
	return nil
}

// Delete value in memcache session
func (rs *MemcacheSessionStore) Delete(key interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.values, key)
	return nil
}

// Flush clear all values in memcache session
func (rs *MemcacheSessionStore) Flush() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values = make(map[interface{}]interface{})
	return nil
}

// SessionID get memcache session id
func (rs *MemcacheSessionStore) SessionID() string {
	return rs.sid
}

// SessionRelease save session values to memcache
func (rs *MemcacheSessionStore) SessionRelease(w http.ResponseWriter) {
	b, err := EncodeGob(rs.values)
	if err != nil {
		return
	}
	item := memcache.Item{Key: rs.sid, Value: b, Expiration: int32(rs.maxlifetime)}
	client.Set(&item)
}

// MemcacheProvider memcache session provider
type MemcacheProvider struct {
	maxlifetime int64
	conninfo    []string
	poolsize    int
	password    string
}

// SessionInit init memcache session
// savepath like
// e.g. 127.0.0.1:9090
func (rp *MemcacheProvider) SessionInit(maxlifetime int64, savePath string) error {
	rp.maxlifetime = maxlifetime
	rp.conninfo = strings.Split(savePath, ";")
	client = memcache.New(rp.conninfo...)
	return nil
}

// SessionRead read memcache session by sid
func (rp *MemcacheProvider) SessionRead(sid string) (Store, error) {
	if client == nil {
		if err := rp.connectInit(); err != nil {
			return nil, err
		}
	}
	item, err := client.Get(sid)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			rs := &MemcacheSessionStore{sid: sid, values: make(map[interface{}]interface{}), maxlifetime: rp.maxlifetime}
			return rs, nil
		}
		return nil, err
	}
	var kv map[interface{}]interface{}
	if len(item.Value) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = DecodeGob(item.Value)
		if err != nil {
			return nil, err
		}
	}
	rs := &MemcacheSessionStore{sid: sid, values: kv, maxlifetime: rp.maxlifetime}
	return rs, nil
}

// SessionExist check memcache session exist by sid
func (rp *MemcacheProvider) SessionExist(sid string) bool {
	if client == nil {
		if err := rp.connectInit(); err != nil {
			return false
		}
	}
	if item, err := client.Get(sid); err != nil || len(item.Value) == 0 {
		return false
	}
	return true
}

// SessionRegenerate generate new sid for memcache session
func (rp *MemcacheProvider) SessionRegenerate(oldsid, sid string) (Store, error) {
	if client == nil {
		if err := rp.connectInit(); err != nil {
			return nil, err
		}
	}
	var contain []byte
	if item, err := client.Get(sid); err != nil || len(item.Value) == 0 {
		// oldsid doesn't exists, set the new sid directly
		// ignore error here, since if it return error
		// the existed value will be 0
		item.Key = sid
		item.Value = []byte("")
		item.Expiration = int32(rp.maxlifetime)
		client.Set(item)
	} else {
		client.Delete(oldsid)
		item.Key = sid
		item.Expiration = int32(rp.maxlifetime)
		client.Set(item)
		contain = item.Value
	}

	var kv map[interface{}]interface{}
	if len(contain) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		var err error
		kv, err = DecodeGob(contain)
		if err != nil {
			return nil, err
		}
	}

	rs := &MemcacheSessionStore{sid: sid, values: kv, maxlifetime: rp.maxlifetime}
	return rs, nil
}

// SessionDestroy delete memcache session by id
func (rp *MemcacheProvider) SessionDestroy(sid string) error {
	if client == nil {
		if err := rp.connectInit(); err != nil {
			return err
		}
	}

	return client.Delete(sid)
}

func (rp *MemcacheProvider) connectInit() error {
	client = memcache.New(rp.conninfo...)
	return nil
}

// SessionGC Impelment method, no used.
func (rp *MemcacheProvider) SessionGC() {
}

// SessionAll return all activeSession
func (rp *MemcacheProvider) SessionAll() int {
	return 0
}

func init() {
	Register("memcache", memcachepder)
}
