package internal

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"sync"
)

const (
	_BodyName       = "_body_"
	_RSign          = "_relate_sign_"
	_USign          = "__user_sign_"
	_OPENID         = "_open_id_"
	_AccessToken    = "_access_token_"
	_Extras         = "_extras_"
	_AcceptLanguage = "_Language_"
	_Brand          = "_Brand"
	_Project        = "_Project"
	_Product        = "_Product"
	_Client         = "_Client"
	_Env            = "_Env"
)

var contextPool = &sync.Pool{
	New: func() interface{} {
		return new(context)
	},
}

func NewContext(ctx *gin.Context) Context {
	context := contextPool.Get().(*context)
	context.ctx = ctx
	return context
}

func ReleaseContext(ctx Context) {
	c := ctx.(*context)
	c.ctx = nil
	contextPool.Put(c)
}

var _ Context = (*context)(nil)

type Context interface {
	init()
	RSign() string
	USign() string
	OpenId() string
	Extras() string
	AccessToken() string
	Language() string
	Brand() string
	Project() string
	Product() string
	Client() string
	Env() string
	SetRSign(RID string)
	SetUSign(UID string)
	SetOpenId(openId string)
	SetAccessToken(accessToken string)
	SetExtras(extras string)
	SetLanguage(language string)
	SetBrand(brand string)
	SetProject(project string)
	SetProduct(product string)
	SetClient(device string)
	SetEnv(Env string)
}

type context struct {
	ctx *gin.Context
}

func (c *context) init() {
	body, err := c.ctx.GetRawData()
	if err != nil {
		panic(err)
	}

	c.ctx.Set(_BodyName, body)
	c.ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
}

func (c *context) RSign() string {
	val, ok := c.ctx.Get(_RSign)
	if !ok {
		return ""
	}

	return fmt.Sprint(val)
}

func (c *context) USign() string {
	val, ok := c.ctx.Get(_USign)
	if !ok {
		return ""
	}
	return fmt.Sprint(val)
}

func (c *context) OpenId() string {
	val, ok := c.ctx.Get(_OPENID)
	if !ok {
		return ""
	}

	return fmt.Sprint(val)
}

func (c *context) Extras() string {
	val, ok := c.ctx.Get(_Extras)
	if !ok {
		return ""
	}

	return fmt.Sprint(val)
}

func (c *context) AccessToken() string {
	val, ok := c.ctx.Get(_AccessToken)
	if !ok {
		return ""
	}

	return fmt.Sprint(val)
}

func (c *context) Language() string {
	val, ok := c.ctx.Get(_AcceptLanguage)
	if !ok {
		return ""
	}

	return fmt.Sprint(val)
}

func (c *context) Brand() string {
	val, ok := c.ctx.Get(_Brand)
	if !ok {
		return ""
	}

	return fmt.Sprint(val)
}

func (c *context) Project() string {
	val, ok := c.ctx.Get(_Project)
	if !ok {
		return ""
	}

	return fmt.Sprint(val)
}

func (c *context) Product() string {
	val, ok := c.ctx.Get(_Product)
	if !ok {
		return ""
	}

	return fmt.Sprint(val)
}

func (c *context) Client() string {
	val, ok := c.ctx.Get(_Client)
	if !ok {
		return ""
	}

	return fmt.Sprint(val)
}

func (c *context) Env() string {
	val, ok := c.ctx.Get(_Env)
	if !ok {
		return ""
	}

	return fmt.Sprint(val)
}

func (c *context) SetRSign(rSign string) {
	c.ctx.Set(_RSign, rSign)
}

func (c *context) SetUSign(uSign string) {
	c.ctx.Set(_USign, uSign)
}

func (c *context) SetOpenId(openId string) {
	c.ctx.Set(_OPENID, openId)
}

func (c *context) SetAccessToken(accessToken string) {
	c.ctx.Set(_AccessToken, accessToken)
}

func (c *context) SetExtras(extras string) {
	c.ctx.Set(_Extras, extras)
}

func (c *context) SetLanguage(language string) {
	c.ctx.Set(_AcceptLanguage, language)
}

func (c *context) SetBrand(brand string) {
	c.ctx.Set(_Brand, brand)
}

func (c *context) SetProject(project string) {
	c.ctx.Set(_Project, project)
}

func (c *context) SetProduct(product string) {
	c.ctx.Set(_Product, product)
}

func (c *context) SetClient(client string) {
	c.ctx.Set(_Client, client)
}

func (c *context) SetEnv(Env string) {
	c.ctx.Set(_Env, Env)
}
