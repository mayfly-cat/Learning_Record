package conf

var (
	RedisCacheConfig RedisCacheConfigData // RedisCacheConfig 缓存配置信息
)

const (
	EnvLive = 0 //获取live标识节点
	EnvGray = 1 //获取gray标识节点
	EnvAll  = 2 //获取全部节点
)

// RedisCacheConfigData redis缓存配置
type RedisCacheConfigData struct {
	RedisAddr string `json:"redisAddr"` //redis连接ip地址
	RedisPwd  string `json:"redisPwd"`  //redis连接password
}

// ServiceIpPortInfo 服务ip端口信息结构
type ServiceIpPortInfo struct {
	IPAddr   string  `json:"ipAddr"`
	Port     uint64  `json:"port"`
	MetaData string  `json:"metaData"`
	Weight   float64 `json:"weight"`
}
