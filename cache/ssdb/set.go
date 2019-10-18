package ssdb

//返回 zset 中的元素个数.
func (sd *SSDB) ZSize(setName string) (int64, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return 0, err
	}
	defer c.Close()

	return c.ZSize(setName)
}

// 设置 zset 中指定 key 对应的权重值.
func (sd *SSDB) ZSet(setName string, key string, score int64) error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	return c.ZSet(setName, key, score)
}

//批量设置 zset 中的 key-score.
func (sd *SSDB) MultiZSet(setName string, kvs map[string]int64) error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	return c.MultiZSet(setName, kvs)
}

//批量获取 zset 中的 key-score.
//
//  setName zset名称
//  key 要获取key的列表，支持多个key
//  返回 val 包含 key-score 的map
func (sd *SSDB) MultiZGet(setName string, keys []string) (map[string]int64, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	return c.MultiZGet(setName, keys...)
}

// 使 zset 中的 key 对应的值增加 num. 参数 num 可以为负数.
//  setName zset名称
//  key 要增加权重的key
//  num 要增加权重值
//  返回 int64 增加后的新权重值
//  返回 err，可能的错误，操作成功返回 nil
func (sd *SSDB) ZIncr(setName string, key string, num int64) (int64, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return 0, err
	}
	defer c.Close()

	return c.ZIncr(setName, key, num)
}

//获取 zset 中指定 key 对应的权重值.
func (sd *SSDB) ZGet(setName, key string) (int64, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return 0, err
	}
	defer c.Close()

	return c.ZGet(setName, key)
}

//删除 zset 中指定 key
func (sd *SSDB) ZDel(setName, key string) error {
	c, err := sd.spool.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	return c.ZDel(setName, key)
}

//判断指定的 key 是否存在于 zset 中.
func (sd *SSDB) ZExists(setName, key string) (bool, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return false, err
	}
	defer c.Close()

	return c.ZExists(setName, key)
}

//根据下标索引区间 [offset, offset + limit) 获取 key-score 对, 下标从 0 开始.注意! 本方法在 offset 越来越大时, 会越慢!
//
//  setName zset名称
//  offset 从此下标处开始返回. 从 0 开始.
//  limit  最多返回这么多个 key-score 对.
//  返回 val 排名
//  返回 err，可能的错误，操作成功返回 nil
func (sd *SSDB) ZRange(setName string, offset, limit int64) (val map[string]int64, err error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	return c.ZRange(setName, offset, limit)
}

//根据下标索引区间 [offset, offset + limit) 获取 key-score 对, 反向顺序获取.注意! 本方法在 offset 越来越大时, 会越慢!
//
//  setName zset名称
//  offset 从此下标处开始返回. 从 0 开始.
//  limit  最多返回这么多个 key-score 对.
//  返回 val 排名
//  返回 err，可能的错误，操作成功返回 nil
func (sd *SSDB) ZRRange(setName string, offset, limit int64) (val map[string]int64, err error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	return c.ZRRange(setName, offset, limit)
}
