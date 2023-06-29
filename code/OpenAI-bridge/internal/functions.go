package internal

import (
	"Learning_Record/code/OpenAI-bridge/config"
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	rand2 "math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// SHA256Sign signed as sha256 - sha256加密字符串
func SHA256Sign(signString string) (r string) {
	h := sha256.New()
	h.Write([]byte(signString))
	cipherStr := h.Sum(nil)
	r = hex.EncodeToString(cipherStr)
	return
}

// RandNumberString return a random numeric string - 生成随机字符串 纯数字
func RandNumberString(n uint8) string {
	var numbers = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	var container string
	length := bytes.NewReader(numbers).Len()

	var i uint8
	for i = 1; i <= n; i++ {
		random, _ := rand.Int(rand.Reader, big.NewInt(int64(length)))
		container += fmt.Sprintf("%d", numbers[random.Int64()])
	}
	return container
}

// RandString Generate Rand String - 生成随机字符串
func RandString(n uint8) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytesData := []byte(str)
	var result []byte
	rand2.Seed(time.Now().UnixNano() + int64(rand2.Intn(100)))
	var i uint8
	for i = 0; i < n; i++ {
		result = append(result, bytesData[rand2.Intn(len(bytesData))])
	}
	return string(result)
}

// Md5Sign Md5 user/relate signature - MD5加密字符串
func Md5Sign(s string) (r string) {
	h := md5.New()
	h.Write([]byte(s))
	cipherStr := h.Sum(nil)
	r = hex.EncodeToString(cipherStr)
	return
}

// UnixTimestampS 10 unix timestamp - 获得秒级时间戳
func UnixTimestampS() int64 {
	return time.Now().Unix()
}

// UnixTimestampMs 13 unix timestamp - 获得毫秒级时间戳
func UnixTimestampMs() int64 {
	return time.Now().UnixNano() / 1e6
}

// NowTimeReturn UnixTime.Now - 获得 2006-01-02 15:04:05 格式时间
func NowTimeReturn(timeFormat string) (s string) {
	var timezone, _ = time.LoadLocation(config.Timezone)
	return time.Now().In(timezone).Format(timeFormat)
}

// FormatPage format the page for sql limit - 格式化传入分页值
func FormatPage(pageNum, pageSize string) (index, num int) {
	index, _ = strconv.Atoi(pageNum)
	num, _ = strconv.Atoi(pageSize)
	if index == 0 {
		index = 1
	}
	if num == 0 {
		num = 10
	}
	index = (index - 1) * num
	return index, num
}

// TokenEnvString 返回环境信息
func TokenEnvString() string {
	return os.Getenv("PROFILES")
}

// CreateSign create the param sign - (json)生成sign-内部server校验
func CreateSign(createData interface{}, key string) (sign, signStr string) {
	stringReader := ""
	t := reflect.TypeOf(createData)
	v := reflect.ValueOf(createData)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() && len(fmt.Sprint(v.Field(i).Interface())) > 0 {
			stringReader = stringReader + strings.ToLower(t.Field(i).Name[:1]) + t.Field(i).Name[1:] + "=" + fmt.Sprint(v.Field(i).Interface()) + "&"
		}
	}
	signStr = stringReader + "key=" + key
	sign = Md5Sign(signStr)
	return sign, signStr
}
