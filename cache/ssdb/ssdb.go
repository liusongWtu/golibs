package ssdb

import (
	"errors"
	"fmt"
	"github.com/seefan/gossdb"
	gconf "github.com/seefan/gossdb/conf"
	"github.com/seefan/gossdb/pool"
	"time"
)

//NewSSDB create new ssdb adapter.
func NewSSDB() *SSDB {
	return &SSDB{}
}

// SSDB adapter
type SSDB struct {
	spool    *pool.Connectors
	conninfo map[string]interface{}
}

// Get value from SSDB.
func (sd *SSDB) Get(key string) (interface{}, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	value, err := c.Get(key)
	if err == nil {
		return value, nil
	}
	return nil, nil
}

// GetMulti get value from SSDB.
func (sd *SSDB) GetMulti(keys []string) ([]interface{}, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	size := len(keys)
	var values []interface{}
	res, err := c.Do("multi_get", keys)
	resSize := len(res)
	if err == nil {
		for i := 1; i < resSize; i += 2 {
			values = append(values, res[i+1])
		}
		return values, nil
	}
	for i := 0; i < size; i++ {
		values = append(values, err)
	}
	return values, nil
}

// DelMulti get value from memSSDB.
func (sd *SSDB) DelMulti(keys []string) error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Do("multi_del", keys)
	return err
}

// Put put value to memSSDB. only support string.
func (sd *SSDB) Put(key string, value interface{}, timeout time.Duration) error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	v, ok := value.(string)
	if !ok {
		return errors.New("value must string")
	}
	var resp []string
	ttl := int(timeout / time.Second)
	if ttl < 0 {
		resp, err = c.Do("set", key, v)
	} else {
		resp, err = c.Do("setx", key, v, ttl)
	}
	if err != nil {
		return err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return nil
	}
	return errors.New("bad response")
}

// Delete delete value in memSSDB.
func (sd *SSDB) Delete(key string) error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	err = c.Del(key)
	return err
}

// Incr increase counter.
func (sd *SSDB) Incr(key string) error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Do("incr", key, 1)
	return err
}

// Decr decrease counter.
func (sd *SSDB) Decr(key string) error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Do("incr", key, -1)
	return err
}

// IsExist check value exists in memSSDB.
func (sd *SSDB) IsExist(key string) (bool, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return false, err
	}
	defer c.Close()

	resp, err := c.Do("exists", key)
	if err != nil {
		return false, err
	}
	if len(resp) == 2 && resp[1] == "1" {
		return true, nil
	}
	return false, nil

}

// ClearAll clear all SSDBd in memSSDB.
func (sd *SSDB) ClearAll() error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	keyStart, keyEnd, limit := "", "", 50
	resp, err := sd.Scan(keyStart, keyEnd, limit)
	for err == nil {
		size := len(resp)
		if size == 1 {
			return nil
		}
		keys := []string{}
		for i := 1; i < size; i += 2 {
			keys = append(keys, resp[i])
		}
		_, e := c.Do("multi_del", keys)
		if e != nil {
			return e
		}
		keyStart = resp[size-2]
		resp, err = sd.Scan(keyStart, keyEnd, limit)
	}
	return err
}

// Scan key all SSDBd in ssdb.
func (sd *SSDB) Scan(keyStart string, keyEnd string, limit int) ([]string, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	resp, err := c.Do("scan", keyStart, keyEnd, limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// actually do the ssdb cmds
func (sd *SSDB) Do(commandName string, args ...interface{}) ([]string, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	return sd.do(c, commandName, args...)
}

// actually do the ssdb cmds, args[0] must be the key name.
func (sd *SSDB) do(c *pool.Client, commandName string, args ...interface{}) ([]string, error) {
	if commandName == "" {
		return nil, errors.New("missing command argument")
	}
	if len(args) == 0 {
		return nil, errors.New("missing required arguments")
	}
	param := []interface{}{commandName}
	for _, v := range args {
		param = append(param, v)
	}
	res, err := c.Do(param...)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	if fmt.Sprintf("%v", res[0]) == "ok" {
		return res[1:], nil
	}
	return nil, nil
}

// StartAndGC start memSSDB adapter.
func (sd *SSDB) StartAndGC(config map[string]interface{}) error {
	sd.conninfo = config
	if err := sd.connectInit(); err != nil {
		return err
	}
	return nil
}

// connect to ssdb and keep the connection.
func (sd *SSDB) connectInit() error {
	var err error
	sd.spool, err = gossdb.NewPool(&gconf.Config{
		Host:             sd.conninfo["host"].(string),
		Port:             sd.conninfo["port"].(int),
		HealthSecond:     sd.conninfo["health_second"].(int),
		MaxWaitSize:      sd.conninfo["max_wait_size"].(int),
		MinPoolSize:      sd.conninfo["min_pool_size"].(int),
		MaxPoolSize:      sd.conninfo["max_pool_size"].(int),
		GetClientTimeout: sd.conninfo["get_client_timeout"].(int),
	})
	if err != nil {
		return err
	}
	return nil
}
