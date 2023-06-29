package config

import (
	"flag"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	Client             *openai.Client
	ApiKey             string
	ApiKeyV1           = "Key1" //openAI账户1的appKey
	ApiKeyV2           = "Key2" //openAI账户2的appKey
	Port               string
	Log                *logrus.Logger    //logRus
	Timezone           string            //Time Location
	TimeFormatDateTime string            //Time Format Datetime String
	SignKey            map[string]string //The key used to encrypt the sign
)

func init() {
	flag.StringVar(&Port, "port", ":8080", "http listen port")
	//加密sign使用的key值
	SignKey = map[string]string{
		"appKey1": "secretKey1", //业务渠道1
		"appKey2": "secretKey2", //业务渠道2
		//业务渠道n... 可拓展
	}
	//初始化请求openAI接口的Client
	ApiKey = ApiKeyV1
	//初始化请求openAI接口的Client
	Client = openai.NewClient(ApiKey)
	////本地请求经过代理----身份验证token---本地自测时可开启该配置
	//conf := openai.DefaultConfig(ApiKey)
	////开启本地代理
	//proxyUrl, err := url.Parse("http://localhost:4003")
	//if err != nil {
	//	panic(err)
	//}
	//transport := &http.Transport{
	//	Proxy: http.ProxyURL(proxyUrl),
	//}
	//conf.HTTPClient = &http.Client{
	//	Transport: transport,
	//}
	//Client = openai.NewClientWithConfig(conf)
}

// InitLog 初始化LogRus
func InitLog() {
	//初始化LogRus
	Log = logrus.New()
	//读取日志level
	Log.Level = logrus.DebugLevel
	Log.Formatter = &logrus.JSONFormatter{
		DisableHTMLEscape: true,
	}
	Log.Out = os.Stdout
}

// ChangeAppKey 变更请求openAI接口使用的appKey
func ChangeAppKey() {
	if ApiKey == ApiKeyV1 { //原使用appKey为账号1
		ApiKey = ApiKeyV2
	} else if ApiKey == ApiKeyV2 { //原使用appKey为账号2
		ApiKey = ApiKeyV1
	}
	Client = openai.NewClient(ApiKey)
}
