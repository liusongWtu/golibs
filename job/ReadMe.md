1:动作名称，服务端标识(manager动作注册名称)，客户端标识，可配置参数列表----p1:"string123",  p2:2,   p3:[2,4], p5:[[1,2],[3,4]]， p5:{"n1":1,"n2":"2"}

2:服务端对应处理定义
 2.1   manager定义

    type ActionManager interface {
        // 空方法
         Check()
    }

 2.2  action定义

    func NewAction2() ActionManager {
        return &Action2{}
    }
    
    type Action2 struct {}
    // 空方法
    func (a *Action2) Check() {}
    
    func (a *Action2) Run(s1,s2 interface{}) string {
        return s1.(string) + " " + s2.(string)
    }
    
    func init() {
        Register("action2", NewAction2)
    }
3:服务端对应处理逻辑

  3.1 反射方式获取Action类名
	manager,_ := NewActionManager(flg)
	t := reflect.TypeOf(manager)
	name := strings.Split( t.String(),".")
  3.2 根据Action类名，调用Run方法
    1--根据类名，获取后台管理系统中配置的参数列表
	2--其他参数获取
	3--初始化类
	4--调用Run方法，按顺序传入参数
    5--返回处理结果(客户端动作标识，客户端所需函数列表)

