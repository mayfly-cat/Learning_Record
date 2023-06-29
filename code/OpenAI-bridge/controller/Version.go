package controller

import (
	"Learning_Record/code/OpenAI-bridge/config"
	"Learning_Record/code/OpenAI-bridge/internal"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommonController struct {
	BaseController
}

// HealthCheck 运维规定心跳方法
func (v *CommonController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		// 正常范围：2xx - 3xx
		// 异常范围：404 - 5xx
		"status":    http.StatusOK,
		"msg":       "Up",
		"version":   "1.0.0",
		"timestamp": internal.NowTimeReturn(config.TimeFormatDateTime),
	})
}
