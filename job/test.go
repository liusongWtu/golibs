package job

import (
	"reflect"
	"strings"
)

func main() {
	flg := "action2"
	manager,_ := NewActionManager(flg)
	t := reflect.TypeOf(manager)
	name := strings.Split( t.String(),".")
	Run(name[1])
}

func Run(className string) map[string]interface{} {
	result := make(map[string]interface{})
	switch className {
	case "Action1":
		act := &Action1{}
		res := act.Hello("hello1","world1")
		result["name"] = "Action1"
		result["data"] = res
	case "Action2":
		act := &Action2{}
		res := act.Hello2("hello2","world2")
		result["name"] = "Action1"
		result["data"] = res
	}
	return result
}