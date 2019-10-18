package job

import (
	"fmt"
)

// ActionManager interface contains all behaviors for manager.
type ActionManager interface {
	// 空方法
	 Check()
}

// ActionInstance is a function create a new ActionManager Instance
type ActionInstance func() ActionManager

var adapters = make(map[string]ActionInstance)

// Register makes a cp adapter available by the name.
func Register(name string, adapter ActionInstance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewActionManager Create a new cache driver by flg.
func NewActionManager(flg string) (task ActionManager, err error) {
	instanceFunc, ok := adapters[flg]
	if !ok {
		err = fmt.Errorf("action manager: unknown newer name %q ", flg)
		return
	}
	task = instanceFunc()
	return
}