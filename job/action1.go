package job

func NewAction1() ActionManager {
	return &Action1{}
}

type Action1 struct {}
// 空方法
func (a *Action1) Check() {}

func (a *Action1) Hello(s1,s2 interface{}) string {
	return s1.(string)+" "+s2.(string)
}

func init() {
	Register("action1", NewAction1)
}
