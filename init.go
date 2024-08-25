package sdkbox

import (
	"github.com/textthree/provider/config"
	"github.com/textthree/provider/core"
	"sdkbox/alioss"
)

var services *core.ServicesContainer

func init() {
	services = core.NewContainer()
	services.Bind(&config.ConfigProvider{})
	services.Bind(&alioss.AliossProvider{})
}

func Svc() *core.ServicesContainer {
	return services
}
