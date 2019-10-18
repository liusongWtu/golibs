package ssdb

//往队列的首部添加一个或者多个元素
//
//  name  队列的名字
//  value  存贮的值，可以为多值.
//  返回 size，添加元素之后, 队列的长度
//  返回 err，执行的错误，操作成功返回 nil
func (sd *SSDB) QPushFront(name string, value interface{}) (int64, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return 0, err
	}
	defer c.Close()
	return c.QPushFront(name, value)
}

//往队列的尾部添加一个或者多个元素
//
//  name  队列的名字
//  value  存贮的值，可以为多值.
//  返回 size，添加元素之后, 队列的长度
//  返回 err，执行的错误，操作成功返回 nil
func (sd *SSDB) QPushBack(name string, value interface{}) (int64, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return 0, err
	}
	defer c.Close()
	return c.QPushBack(name, value)
}

//从队列首部弹出最后一个元素.
//
//  name 队列的名字
//  返回 v，返回一个元素，并在队列中删除 v；队列为空时返回空值
//  返回 err，执行的错误，操作成功返回 nil
func (sd *SSDB) QPopFront(name string) (string, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return "", err
	}
	defer c.Close()
	val, err := c.QPopFront(name)
	if err != nil {
		return "", err
	}
	return sd.getString(val), nil
}

//从队列尾部弹出最后一个元素.
//
//  name 队列的名字
//  返回 v，返回一个元素，并在队列中删除 v；队列为空时返回空值
//  返回 err，执行的错误，操作成功返回 nil
func (sd *SSDB) QPopBack(name string) (string, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return "", err
	}
	defer c.Close()
	val, err := c.QPopBack(name)
	if err != nil {
		return "", err
	}
	return sd.getString(val), nil
}

//返回队列的长度.
//
//  name  队列的名字
//  返回 size，队列的长度；
//  返回 err，执行的错误，操作成功返回 nil
func (sd *SSDB) QSize(name string) (int64, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return 0, err
	}
	defer c.Close()
	return c.QSize(name)
}

//返回指定位置的元素. 0 表示第一个元素, 1 是第二个 ... -1 是最后一个.
//
//  key  队列的名字
//  index 指定的位置，可传负数.
//  返回 val，返回的值.
//  返回 err，执行的错误，操作成功返回 nil
func (sd *SSDB) QGet(name string, index int64) (string, error) {
	c, err := sd.spool.NewClient()
	if err != nil {
		return "", err
	}
	defer c.Close()
	val, err := c.QGet(name, index)
	if err != nil {
		return "", err
	}
	return sd.getString(val), nil
}
