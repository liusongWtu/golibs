// Package redis for cache provider
//
// depend on github.com/gomodule/redigo/redis
//
// go install github.com/gomodule/redigo/redis
//
// Usage:
// import(
//     "libs/redis"
// )
//
//  c, err := cache.NewCache("redis", `{"conn":"127.0.0.1:6379"}`)
//
package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"strings"
	"time"
)

var (
	// DefaultKey the collection name of redis for cache adapter.
	DefaultKey = ""
)

// Cache is Redis cache adapter.
type RedisCache struct {
	p         *redis.Pool // redis connection pool
	conninfo  string
	dbNum     int
	key       string
	password  string
	maxIdle   int
	maxActive int
}

// NewRedisCache create new redis cache with default collection name.
func NewRedisCache() *RedisCache {
	return &RedisCache{key: DefaultKey}
}

func (rc *RedisCache) GetConn() redis.Conn {
	return rc.p.Get()
}

//actually do the redis cmds by conn
func (rc *RedisCache) DoByConn(conn redis.Conn, commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	return conn.Do(commandName, args...)
}

// actually do the redis cmds
func (rc *RedisCache) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return rc.do(commandName, args...)
}

// actually do the redis cmds, args[0] must be the key name.
func (rc *RedisCache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	args[0] = rc.associate(args[0])
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// associate with config key.
func (rc *RedisCache) associate(originKey interface{}) string {
	if rc.key == "" {
		return fmt.Sprintf("%s", originKey)
	}
	return fmt.Sprintf("%s:%s", rc.key, originKey)
}

func (rc *RedisCache) SetStruct(key string, val interface{}) error {
	value, err := json.Marshal(val)
	if err != nil {
		return err
	}
	_, err = rc.Do("SET", key, value)
	return err
}

func (rc *RedisCache) SetStructWithExpire(key string, val interface{}, timeout time.Duration) error {
	value, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return rc.Put(key, value, timeout)
}

// Set cache to redis.
func (rc *RedisCache) Set(key string, val interface{}) error {
	_, err := rc.do("SET", key, val)
	return err
}

// Get cache from redis.
func (rc *RedisCache) Get(key string) interface{} {
	if v, err := rc.do("GET", key); err == nil {
		return v
	}
	return nil
}

// Put put cache to redis.
func (rc *RedisCache) Expire(key string, timeout time.Duration) error {
	_, err := rc.do("EXPIRE", key, int64(timeout/time.Second))
	return err
}

// GetMulti get cache from redis.
func (rc *RedisCache) GetMulti(keys []string) []interface{} {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	values, err := redis.Values(c.Do("MGET", args...))
	if err != nil {
		return nil
	}
	return values
}

// Put put cache to redis.
func (rc *RedisCache) Put(key string, val interface{}, timeout time.Duration) error {
	_, err := rc.do("SETEX", key, int64(timeout/time.Second), val)
	return err
}

// Delete delete cache in redis.
func (rc *RedisCache) Delete(key string) error {
	_, err := rc.do("DEL", key)
	return err
}

// IsExist check cache's existence in redis.
func (rc *RedisCache) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

// Incr increase counter in redis.
func (rc *RedisCache) Incr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, 1))
	return err
}

// Decr decrease counter in redis.
func (rc *RedisCache) Decr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, -1))
	return err
}

// Hash Multi Set in redis.
func (rc *RedisCache) Hmset(key string, data map[string]string) error {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	args = append(args, rc.associate(key))
	for k, v := range data {
		args = append(args, k, v)
	}
	_, err := c.Do("HMSET", args...)
	return err
}

// Hash Multi Get in redis.
func (rc *RedisCache) Hmget(key string, keys []string) (interface{}, error) {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	args = append(args, rc.associate(key))
	for _, v := range keys {
		args = append(args, v)
	}
	res, err := c.Do("HMGET", args...)
	return res, err
}

// ClearAll clean all cache in redis. delete this redis collection.
func (rc *RedisCache) ClearAll() error {
	c := rc.p.Get()
	defer c.Close()
	cachedKeys, err := redis.Strings(c.Do("KEYS", rc.key+":*"))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = c.Do("DEL", str); err != nil {
			return err
		}
	}
	return err
}

// StartAndGC start redis cache adapter.
// config is like {"key":"collection key","conn":"127.0.0.1:6379","dbnum":"0","password":"","maxidle":3,"maxactive":2000}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *RedisCache) StartAndGC(conf map[string]string) error {
	if _, ok := conf["key"]; !ok {
		conf["key"] = DefaultKey
	}
	if _, ok := conf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	// Format redis://<password>@<host>:<port>
	conf["conn"] = strings.Replace(conf["conn"], "redis://", "", 1)
	if i := strings.Index(conf["conn"], "@"); i > -1 {
		conf["password"] = conf["conn"][0:i]
		conf["conn"] = conf["conn"][i+1:]
	}
	if _, ok := conf["dbnum"]; !ok {
		conf["dbnum"] = "0"
	}
	if _, ok := conf["password"]; !ok {
		conf["password"] = ""
	}
	if _, ok := conf["maxidle"]; !ok {
		conf["maxidle"] = "3"
	}
	if _, ok := conf["maxactive"]; !ok {
		conf["maxactive"] = "0"
	}

	rc.key = conf["key"]
	rc.conninfo = conf["conn"]
	rc.dbNum, _ = strconv.Atoi(conf["dbnum"])
	rc.password = conf["password"]
	rc.maxIdle, _ = strconv.Atoi(conf["maxidle"])
	rc.maxActive, _ = strconv.Atoi(conf["maxactive"])

	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()

	return c.Err()
}

// connect to redis.
func (rc *RedisCache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     rc.maxIdle,
		MaxActive:   rc.maxActive,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}
