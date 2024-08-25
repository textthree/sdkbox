package alioss

import (
	"github.com/textthree/provider/config"
	"github.com/textthree/provider/core"
	"sync"
)

const Name = "alioss"

func (self *AliossProvider) Name() string {
	return Name
}

var instance *AliossService

type AliossProvider struct {
	core.ServiceProvider
}

func (*AliossProvider) RegisterProviderInstance(c core.Container) core.NewInstanceFunc {
	return func(params ...interface{}) (interface{}, error) {
		instance = &AliossService{
			c:      c,
			lock:   sync.Mutex{},
			cfgSvc: c.NewSingle(config.Name).(config.Service),
		}
		return instance, nil
	}
}

func (*AliossProvider) InitOnBind() bool {
	return false
}

func (*AliossProvider) BeforeInit(c core.Container) error {
	return nil
}

func (*AliossProvider) Params(c core.Container) []interface{} {
	return []interface{}{c}
}

func (*AliossProvider) AfterInit(instance any) error {
	return nil
}
