package sms

import "fmt"

// SmsManager interface contains all behaviors for manager.
type SmsManager interface {
	//检查任务是否完成
	Send(data map[string]string) error
}

// SmsInstance is a function create a new SmsManager Instance
type SmsInstance func() SmsManager

var adapters = make(map[string]SmsInstance)

// Register makes a cp adapter available by the name.
func Register(name string, adapter SmsInstance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewSmsManager Create a new cache driver by flg.
func NewSmsManager(flg string) (task SmsManager, err error) {
	instanceFunc, ok := adapters[flg]
	if !ok {
		err = fmt.Errorf("sms manager: unknown newer name %q ", flg)
		return
	}
	task = instanceFunc()
	return
}
