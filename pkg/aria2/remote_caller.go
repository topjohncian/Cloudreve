package aria2

import (
	"encoding/json"
	model "github.com/cloudreve/Cloudreve/v3/models"
	"github.com/cloudreve/Cloudreve/v3/pkg/aria2/rpc"
	"github.com/cloudreve/Cloudreve/v3/pkg/auth"
	"github.com/cloudreve/Cloudreve/v3/pkg/request"
	"github.com/cloudreve/Cloudreve/v3/pkg/serializer"
	"github.com/cloudreve/Cloudreve/v3/pkg/util"
	"net/url"
	"path"
	"strings"
)

// RemoteService 通过RPC服务的Aria2任务管理器
type RemoteService struct {
	Policy       model.Policy
	AuthInstance auth.Auth
	Caller       request.Client
}

func (client *RemoteService) CreateTask(task *model.Download, options map[string]interface{}) error {
	return nil
}

func (client *RemoteService) Status(task *model.Download) (rpc.StatusInfo, error) {
	return rpc.StatusInfo{}, nil
}

func (client *RemoteService) Cancel(task *model.Download) error {
	return nil
}

func (client *RemoteService) Select(task *model.Download, files []int) error {
	return nil
}

// Init 初始化
func (client *RemoteService) Init(policy model.Policy, options map[string]string) error {
	client.Policy = policy
	client.Caller = request.HTTPClient{}
	client.AuthInstance = auth.HMACAuth{SecretKey: []byte(policy.SecretKey)}

	reqBody := serializer.Aria2InitRequest{
		Options: options,
	}
	reqBodyEncoded, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// 发送列表请求
	bodyReader := strings.NewReader(string(reqBodyEncoded))
	signTTL := model.GetIntSetting("slave_api_timeout", 60)
	resp, err := client.Caller.Request(
		"POST",
		client.getAPIUrl("init"),
		bodyReader,
		request.WithCredential(client.AuthInstance, int64(signTTL)),
	).CheckHTTPResponse(200).DecodeResponse()
	if err != nil {
		return err
	}

	if resStr, ok := resp.Data.(string); ok {
		util.Log().Info(resStr)
	}
	return nil
}

func (client RemoteService) getAPIUrl(scope string, routes ...string) string {
	serverURL, err := url.Parse(client.Policy.Server)
	if err != nil {
		return ""
	}
	var controller *url.URL

	switch scope {
	case "init":
		controller, _ = url.Parse("/api/v3/slave/aria2/init")
	default:
		controller = serverURL
	}

	for _, r := range routes {
		controller.Path = path.Join(controller.Path, r)
	}

	return serverURL.ResolveReference(controller).String()
}
