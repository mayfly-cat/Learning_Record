package internal

/*
新错误码定义：111223444
111 : 应用标识，表示错误属于哪个应用，三位数字。（付费会员：125）
22  : 功能域标识，表示错误属于应用中的哪个功能模块，两位数字。（21:jump；22：车机；23：基础服务）
3   : 错误类型，表示错误属于哪种类型，一位数字。（1-系统错误，2-业务错误，3-参数解析错误，4-依赖基础服务错误，9-其他）
444 : 错误编码，错误类型下的具体错误，三位数字。
*/

const (
	CodeSuccess          = 20000     //成功
	CodeNotFound         = 40004     //未找到
	CodeExecError        = 125002001 //执行错误
	CodeParamsError      = 125003001 //参数接收错误
	CodeParamsNotAllowed = 125003002 //参数不合法
	CodeTokenEmpty       = 125003003 //token为空
	CodeTokenError       = 125003004 //token不合法
	CodeTokenExpired     = 125003005 //token过期
	CodeServerErr        = 125009001 //服务错误
)
