package main

import (
	"Learning_Record/code/OpenAI-bridge/config"
	"Learning_Record/code/OpenAI-bridge/route"
)

func main() {
	//初始化日志打印
	config.InitLog()

	r := route.InitRoute()
	r.Run(config.Port)
}
