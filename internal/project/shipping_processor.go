package project

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/spf13/viper"
	"os/exec"
)

type ShippingProcessor interface {
	Login(p *Project) []*exec.Cmd
	Push(p *Project, image string) []*exec.Cmd
}

func NewShippingProcessor(typ ShippingProcessorType) ShippingProcessor {
	switch typ {
	case SHIPPING_PROCESSOR_TYPE__ALIYUN_REGISTRY:
		return new(AliyunShippingProcessor)
	default:
		panic(fmt.Sprintf("unsupported shipping processor type: %s", typ))
	}
}

type AliyunLoginResponseData struct {
	AuthorizationToken string `json:"authorizationToken"`
	TempUserName       string `json:"tempUserName"`
	ExpireDate         uint64 `json:"expireDate"`
}

type AliyunLoginResponse struct {
	Data AliyunLoginResponseData `json:"data"`
}

type AliyunShippingProcessor struct{}

func (a *AliyunShippingProcessor) Login(p *Project) []*exec.Cmd {
	accessKey := viper.GetString("ALIYUN_ACCESS_KEY")
	accessSecret := viper.GetString("ALIYUN_ACCESS_SECRET")
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", accessKey, accessSecret)
	if err != nil {
		panic(err)
	}

	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Scheme = "https" // https | http
	request.Domain = "cr.cn-hangzhou.aliyuncs.com"
	request.Version = "2016-06-07"
	request.PathPattern = "/tokens"
	request.Headers["Content-Type"] = "application/json"

	body := `{}`
	request.Content = []byte(body)

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}

	result := response.GetHttpContentBytes()
	var resp = new(AliyunLoginResponse)
	err = json.Unmarshal(result, resp)
	if err != nil {
		panic(fmt.Sprintf("cannot unmarshal login response: %s", string(result)))
	}

	return []*exec.Cmd{p.Command("docker", "login",
		fmt.Sprintf("--username=%s", resp.Data.TempUserName),
		fmt.Sprintf("--password='%s'", resp.Data.AuthorizationToken),
		"registry.cn-hangzhou.aliyuncs.com")}
}

func (a *AliyunShippingProcessor) Push(p *Project, image string) []*exec.Cmd {
	return []*exec.Cmd{p.Command("docker", "push", image)}
}
