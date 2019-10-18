package job

func NewAction2() ActionManager {
	return &Action2{}
}

type Action2 struct {}
// 空方法
func (a *Action2) Check() {}

func (a *Action2) Hello2(s1,s2 interface{}) string {
	return s1.(string)+" "+s2.(string)
}

func init() {
	Register("action2", NewAction2)
}
