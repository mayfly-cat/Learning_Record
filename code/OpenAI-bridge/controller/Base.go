package controller

import (
	"Learning_Record/code/OpenAI-bridge/config"
	"Learning_Record/code/OpenAI-bridge/internal"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strings"
)

type BaseController struct {
}

type UserCommon struct {
	BaseController
}

type ResData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// InitController Init Handler
func (b *BaseController) InitController(c *gin.Context) {
	ctx := internal.NewContext(c)
	//通过ctx.Language()获取语言
	//也可增加其他可用上下文存入context池
	lang := c.GetHeader("Accept-Language")
	ctx.SetLanguage(lang)
}

// InitVerifySignController 校验sign
func (b *BaseController) InitVerifySignController(c *gin.Context) {
	if strings.ToLower(c.Request.Method) == "post" {
		var signBool bool
		//获取加密的请求体
		postData, _ := io.ReadAll(c.Request.Body)
		bodyReceive := make(map[string]interface{})
		if err := json.Unmarshal(postData, &bodyReceive); err != nil {
			internal.PrintLogRus("error", "controller/Base.InitVerifySignController", "body json unmarshal failed, ", err.Error())
			ApiReturn(c, internal.CodeServerErr, internal.MsgServerErr, "body解析失败")
			return
		}
		var keys []string
		for k := range bodyReceive {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var stringReader, signString, signedKey string
		for _, k := range keys {
			var usedAppendString string
			//sign不进行校验
			if !signBool && "sign" == k {
				signString = fmt.Sprint(bodyReceive[k])
				continue
			}
			//时间戳科学技术法问题
			if k == "timestamp" && reflect.TypeOf(bodyReceive[k]).String() == reflect.Float64.String() {
				var newFloat float64
				_, err := fmt.Sscanf(fmt.Sprint(bodyReceive[k]), "%e", &newFloat)
				if err != nil {
					internal.PrintLogRus("error", "controller/Base.InitVerifySignController", "verify sign float64 scan err, ", bodyReceive[k], err.Error())
					ApiReturn(c, internal.CodeServerErr, internal.MsgServerErr, "参数类型转换错误")
					return
				}
				usedAppendString = fmt.Sprintf("%.f", newFloat)
			} else if k == "data" && (reflect.TypeOf(bodyReceive[k]).String()[:3] == reflect.Map.String() || reflect.TypeOf(bodyReceive[k]).String()[:11] == "[]interface") {
				printData, _ := json.Marshal(bodyReceive[k])
				usedAppendString = string(printData)
			} else {
				usedAppendString = fmt.Sprint(bodyReceive[k])
			}
			if len(usedAppendString) > 0 {
				stringReader = stringReader + k + "=" + usedAppendString + "&"
			}
		}
		appKey := c.GetHeader("appKey")
		stringReader += "key=" + config.SignKey[appKey]
		if signBool {
			signString = c.GetHeader("sign")
			postString := string(postData)
			postString = strings.ReplaceAll(postString, "\r\n", "")
			stringReader = postString + signedKey
		}
		internal.PrintLogRus("debug", "controller/Base.InitVerifySignController", "strings:", stringReader, "signString:", signString, "signedString:", strings.ToUpper(internal.Md5Sign(stringReader)))
		if signString != strings.ToUpper(internal.Md5Sign(stringReader)) {
			ApiReturn(c, internal.CodeParamsError, internal.MsgParamsError, "sign校验失败")
			return
		}
		//重新赋值，避免EOF错误
		c.Request.Body = io.NopCloser(bytes.NewBuffer(postData))
	}
}

// ApiReturn API统一结构返回接口
func ApiReturn(c *gin.Context, code int, message string, data interface{}) {
	internal.PrintLogRus("debug", "Base.ApiReturn", code, message, data) // , c.Request.URL
	c.JSON(http.StatusOK, &ResData{
		Code:    code,
		Message: message,
		Data:    data,
	})
	c.Abort()
}
