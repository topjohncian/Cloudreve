package slave

import (
	"github.com/cloudreve/Cloudreve/v3/pkg/aria2"
	"github.com/cloudreve/Cloudreve/v3/pkg/serializer"
)

type Aria2InitService struct {
	options map[string]string
}

func (service Aria2InitService) Init() serializer.Response {
	aria2.SlaveInit(service.options)
	return serializer.Response{}
}
