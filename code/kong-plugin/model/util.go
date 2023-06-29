package model

import (
	conf "Learning_Record/code/kong-plugin/config"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/Kong/go-pdk"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// GetUserIdFromJwt 请求头Authorization值token中解析用户id信息
func GetUserIdFromJwt(kong *pdk.PDK) string {
	jwt, _ := kong.Request.GetHeader("Authorization")
	if len(jwt) == 0 {
		jwt, _ = kong.Request.GetHeader("token")
		if len(jwt) == 0 {
			return ""
		}
	}
	type kongUserInfo struct {
		UserId string `json:"sub"`
	}
	userInfos := strings.Split(jwt, ".")
	if len(userInfos) == 3 {
		des, err := base64.RawURLEncoding.DecodeString(userInfos[1])
		if err != nil {
			kong.Log.Err("[kong-gray] ", "Get UserId From Jwt DecodeString Error: ", userInfos[1])
			return ""
		}
		var kongUserInfo kongUserInfo
		err = json.Unmarshal(des, &kongUserInfo)
		if err != nil {
			kong.Log.Err("[kong-gray] ", "Get UserId From Jwt json Unmarshal error: ", err.Error())
			return ""
		}
		return kongUserInfo.UserId
	}
	return ""
}

// GetServiceList 获取服务列表 -- 读redis缓存
func GetServiceList(kong *pdk.PDK) (serviceHostList []conf.ServiceIpPortInfo, err error) {
	//获取请求的hostname，用于匹配对应的nacos服务实例列表（通过读取redis）
	serviceHost, err := kong.Request.GetHost()
	kong.Log.Notice("[kong-gray] Kong.GetHost result is: ", serviceHost)
	if err != nil {
		kong.Log.Err("[kong-gray] Kong.GetHost exec failed, ", err.Error())
		return nil, err
	}
	hostNames := strings.Split(serviceHost, ":")
	serviceHost = hostNames[0]

	//匹配hostname对应服务 ip & port
	redisContext := context.Background()
	serviceListStr, _ := Cache.Get(redisContext, "kong_gray_server_domain_"+serviceHost).Result()
	if serviceListStr != "" {
		err = json.Unmarshal([]byte(serviceListStr), &serviceHostList)
		if err != nil {
			kong.Log.Err("[kong-gray] ", "Unmarshal service list json err: ", err.Error())
			return nil, err
		}
	} else { //如域名未对应到服务，则解析到一级路由查询是否存在ip&port列表
		pathname, err := kong.Request.GetPath() // example pathname : /token/refresh
		kong.Log.Notice("[kong-gray] Kong.GetPath result is: ", pathname)
		if err != nil {
			kong.Log.Err("[kong-gray] Kong.GetPath exec failed, ", err.Error())
			return nil, err
		}
		pathNames := strings.Split(pathname, "/")
		if len(pathNames) < 2 {
			return nil, nil
		}
		pathname = pathNames[1]
		serviceListStr, _ = Cache.Get(redisContext, "kong_gray_server_domain_"+serviceHost+"/"+pathname).Result()
		if serviceListStr != "" {
			err = json.Unmarshal([]byte(serviceListStr), &serviceHostList)
			if err != nil {
				kong.Log.Err("[kong-gray] ", "Unmarshal service list json err: ", err.Error())
				return nil, err
			} else {
				return serviceHostList, nil
			}
		} else { // 如果不存在服务名配置，则直接返回
			return nil, nil
		}
	}
	return serviceHostList, err
}

// GetServiceHostListByEnv 获取某环境服务列表
func GetServiceHostListByEnv(env int, serviceHostList []conf.ServiceIpPortInfo) (targetHosts []conf.ServiceIpPortInfo) {
	for _, item := range serviceHostList {
		if env == conf.EnvGray {
			if item.MetaData == "grayTag" {
				targetHosts = append(targetHosts, item)
			}
		} else if env == conf.EnvLive {
			if item.MetaData == "" {
				targetHosts = append(targetHosts, item)
			}
		}
	}
	return targetHosts
}

// GetTargetMachine 随机分发到匹配ip&port
func GetTargetMachine(targetService []conf.ServiceIpPortInfo) (string, uint64) {
	targetHosts := ServiceMachineList(targetService)
	sort.Sort(&targetHosts)

	weights := []float64{}
	liveHosts := ServiceMachineList{}
	for _, target := range targetHosts {
		liveHosts = append(liveHosts, target)
		weights = append(weights, target.Weight)
	}
	if len(liveHosts) <= 0 {
		return "", 0
	}

	ipNo := GetWeight(weights)
	return liveHosts[ipNo].IPAddr, liveHosts[ipNo].Port
}

// ServiceMachineList 服务ip&port列表
type ServiceMachineList []conf.ServiceIpPortInfo

func (sl *ServiceMachineList) Len() int           { return len(*sl) }
func (sl *ServiceMachineList) Less(i, j int) bool { return (*sl)[i].Weight < (*sl)[j].Weight }
func (sl *ServiceMachineList) Swap(i, j int)      { (*sl)[i], (*sl)[j] = (*sl)[j], (*sl)[i] }

// GetWeight 获取权重信息
func GetWeight(x []float64) int {
	length := len(x)
	sum := 0.0
	for i := 0; i < length; i++ {
		sum += x[i]
	}
	randVal := randFloats(0.0, sum)
	idx := 0
	for i := 0; i < length; i++ {
		if randVal <= x[i] {
			idx = i
			break
		}
		randVal -= x[i]
	}

	return idx
}

// randFloats 随机分发
func randFloats(min, max float64) float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randVal := r.Float64()
	return min + randVal*(max-min)
}
