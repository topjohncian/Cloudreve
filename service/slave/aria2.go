package slave

import (
	"github.com/cloudreve/Cloudreve/v3/pkg/aria2"
	"github.com/cloudreve/Cloudreve/v3/pkg/serializer"
)

type Aria2InitService struct {
	Options map[string]string `json:"options"`
}

func (service Aria2InitService) Init() serializer.Response {
	aria2.SlaveInit(service.Options)
	return serializer.Response{}
}
