/*
	A "hello world" plugin in Go,
	which reads a request header and sets a response header.
*/

package main

import (
	conf "Learning_Record/code/kong-plugin/config"
	"Learning_Record/code/kong-plugin/model"
	"context"
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"log"
)

func main() {
	log.Println("enter main")
	//初始化redis连接
	model.InitRedis()
	//插件的入口，设置版本号，处理优先级
	server.StartServer(New, Version, Priority)
}

var Version = "1.0" //插件版本号
var Priority = 1    //执行优先级

type Config struct {
	Message string
}

func New() interface{} {
	log.Println("enter new")
	return &Config{}
}

func (config Config) Access(kong *pdk.PDK) {
	//set header 用于标识经过插件转发的请求
	err := kong.Response.SetHeader("x-hello-from-go-1", "Kong plugin kong-gray, version: 1.0")
	if err != nil {
		kong.Log.Notice("[kong-gray] ", "Set response header err: ", err.Error())
	}
	//获取全部ip&port列表
	srv, err := model.GetServiceList(kong)
	if len(srv) == 0 || err != nil { //如获取列表为空或有错误，则直接返回
		err = kong.Response.SetHeader("x-hello-from-go-2", "can not find the nacos service list")
		if err != nil {
			kong.Log.Notice("[kong-gray] ", "Set response header err: ", err.Error())
		}
		return
	} else {
		//判断是否转发到灰度节点，否则直接转发至live节点
		if parseTag(kong) { //解析请求方法返回 true，则该请求为灰度请求，转发到灰度节点
			if !processGray(srv, kong) {
				//如转发到灰度节点失败，则转发到普通节点
				processLive(srv, kong)
				return
			}
		} else { //解析请求方法返回 false，则该请求为正常请求，转发到正常节点
			processLive(srv, kong)
			return
		}
	}
	return
}

// parseTag 解析请求是否为灰度请求
func parseTag(kong *pdk.PDK) bool {
	//先判断请求是否带有 env 请求头, 值为 grayTag, 如果有则为灰度请求  返回true  转发至灰度节点
	env, _ := kong.Request.GetHeader("env")
	if env == "grayTag" {
		return true
	}
	//如无请求头env, 获取Authorization 解析出用户id信息
	userId := model.GetUserIdFromJwt(kong)
	//如未解析出用户信息，则直接返回false，转发请求到live
	if userId == "" {
		return false
	}
	//匹配用户是否存在于redis匹配规则里,  如存在，则返回true  转发至灰度节点
	redisContext := context.Background()
	isGray, _ := model.Cache.Get(redisContext, "kong_gray_user_"+userId).Result()
	if isGray == "true" {
		return true
	} else {
		kong.Log.Notice("[kong-gray] ", "this user is not gray user, user_id : ", userId)
		return false
	}
}

func processGray(serviceHostList []conf.ServiceIpPortInfo, kong *pdk.PDK) (processed bool) {
	//kong.Response.SetHeader("x-hello-from-go-3", "this is a gray request")
	//筛选灰度标识节点 ip&port列表
	srv := model.GetServiceHostListByEnv(conf.EnvGray, serviceHostList)
	if len(srv) == 0 {
		kong.Log.Notice("[kong-gray] ", "can not find the nacos service list by gray tag")
		return false
	}
	//随机分发到对应ip&port
	ip, port := model.GetTargetMachine(srv)
	kong.Log.Notice("[kong-gray] ", "target ip&port is: ", ip, ":", port)
	err := kong.Service.SetTarget(ip, int(port))
	if err != nil {
		kong.Log.Err("[kong-gray] ", "Set target err: ", err.Error())
	}
	//匹配灰度规则，则请求头增加env参数，值为 grayTag
	err = kong.ServiceRequest.SetHeader("env", "grayTag")
	if err != nil {
		kong.Log.Err("[kong-gray] ", "Set Request Header err: ", err.Error())
	}
	return true
}

func processLive(serviceHostList []conf.ServiceIpPortInfo, kong *pdk.PDK) {
	//筛选live标识节点 ip&port列表
	srv := model.GetServiceHostListByEnv(conf.EnvLive, serviceHostList)
	if len(srv) == 0 {
		srv = serviceHostList
	}
	//随机分发到对应ip&port
	ip, port := model.GetTargetMachine(srv)
	kong.Log.Notice("[kong-gray] ", "target ip&port is: ", ip, ":", port)
	if len(ip) == 0 || port == 0 {
		return
	}
	err := kong.Service.SetTarget(ip, int(port))
	if err != nil {
		kong.Log.Err("[kong-gray] ", "Set target err: ", err.Error())
	}
	return
}
