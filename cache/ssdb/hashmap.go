package ssdb

// Get value from SSDB hashmap.
func (sd *SSDB) HGet(setName string, key string) (interface{}, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	res, err := c.HGet(setName, key)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//批量获取 hashmap 中多个 key 对应的权重值.
//  setName - hashmap 的名字.
//  keys - 包含 key 的数组 .
func (sd *SSDB) MultiHGet(setName string, key ...string) (map[string]interface{}, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	res, err := c.MultiHGet(setName, key...)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	value := make(map[string]interface{})
	for k, v := range res {
		value[k] = v.String()
	}
	return value, nil
}

// Get all value from SSDB hashmap.
func (sd *SSDB) HGetAll(setName string) (map[string]interface{}, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	res, err := c.HGetAll(setName)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	value := make(map[string]interface{})
	for k, v := range res {
		value[k] = v.String()
	}
	return value, nil
}

// Set value to SSDB hashmap.
func (sd *SSDB) HSet(setName string, key string, val interface{}) error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()
	return c.HSet(setName, key, val)
}

//批量设置 hashmap 中的 key-value.
//  setName - hashmap 的名字.
//  kvs - 包含 key-value 的关联数组 .
func (sd *SSDB) MultiHSet(setName string, kvs map[string]interface{}) error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()
	return c.MultiHSet(setName, kvs)
}

//判断指定的 key 是否存在于 hashmap 中.
func (sd *SSDB) HExists(setName string, key string) (bool, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return false, err
	}
	defer c.Close()
	return c.HExists(setName, key)
}

//删除 hashmap 中的所有 key
func (sd *SSDB) HClear(setName string) error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()
	return c.HClear(setName)
}

//删除 hashmap 中的指定 key，不能通过返回值来判断被删除的 key 是否存在.
func (sd *SSDB) HDel(setName string, key string) error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()
	return c.HDel(setName, key)
}

//批量获取 hashmap 中多个 key 对应的权重值.
//  setName - hashmap 的名字.
//  keys - 包含 key 的数组 .
func (sd *SSDB) MultiHDel(setName string, key ...string) error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()
	return c.MultiHDel(setName, key...)
}

//批量删除 hashmap 中的 key.（输入分片）
func (sd *SSDB) MultiHdelArray(setName string, key []string) (err error) {
	return sd.MultiHDel(setName, key...)
}

// 返回 hashmap 中的元素个数.
func (sd *SSDB) HSize(setName string) (int64, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return 0, err
	}
	defer c.Close()
	return c.HSize(setName)
}

//设置 hashmap 中指定 key 对应的值增加 num. 参数 num 可以为负数.
func (sd *SSDB) HIncr(setName string, key string, num int64) (interface{}, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	value, err := c.HIncr(setName, key, num)
	if err != nil {
		return nil, err
	}
	return value, nil
}
