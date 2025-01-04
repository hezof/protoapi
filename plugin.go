package protoapi

import (
	"bytes"
	"strings"
)

/*
ServicePlugin 服务插件. 通过此扩展接口可以:
1. 动态操作服务集(添加/修改/删除)
2. 打印或检查相关的服务元数据
*/
type ServicePlugin func(all *[]*ServiceSetting)

/*
RequestPlugin 请求插件.通过此扩展接口可以:
1. proxy: 将RequestSetting.Handles替换为其他请求的Handles
2. cache: 将RequestSetting.Handles最内层(或者最后一环)使用cache机制包裹,实现请求cache机制
3. access: 在RequestSetting.Plugins最外层使用logger计算耗时并打印所需的内容信息
4. others: 基于Context.Next()机制可以实现的拦截逻辑
*/
type RequestPlugin func(all map[string]map[string]*RequestSetting)

/*
MessageValidatePluginProvider message校验插件提供者
*/
type MessageValidatePluginProvider func(args []string) MessageValidatePlugin

var globalMessageValidatePluginProvider = make(map[string]MessageValidatePluginProvider)

func SetMessageValidatePluginProvider(k string, p MessageValidatePluginProvider) {
	globalMessageValidatePluginProvider[k] = p
}

/*
FieldValidatePluginProvider field校验插件提供者
*/
type FieldValidatePluginProvider func(args []string) FieldValidatePlugin

var globalFieldValidatePluginProvider = make(map[string]FieldValidatePluginProvider)

func SetFieldValidatePluginProvider(k string, p FieldValidatePluginProvider) {
	globalFieldValidatePluginProvider[k] = p
}

/*
CompilePluginExpression 编译插件表达式, 其固定语法为: plugin(arg1,arg2,...), 参数之间使用逗号","分隔并忽略二侧空白.
其中(,)\是元字符,如果出现表在表达式内请使用转义字符"\"转义!, 相应规则为:
1. 最左侧的(与最右侧的), 中间为args部分
2. 从左到右扫描",", 判断左边是否"\", 如果不是则作为分隔符, 如果是则作为普通字符! 如果出现参数以"\"结尾的情况, 请在","前加个空格!
*/
func CompilePluginExpression(expr string) (name string, args []string) {
	// 没有(则没参数
	lp := strings.IndexByte(expr, '(')
	if lp < 0 {
		return expr, nil
	}

	// 没有)则没参数
	rp := strings.LastIndexByte(expr, ')')
	if rp < lp {
		return expr, nil
	}

	// 解析名称
	name = strings.TrimSpace(expr[0:lp])

	// 解析参数
	expr = strings.TrimSpace(expr[lp+1 : rp])
	if n := len(expr); n > 0 {
		sb := new(bytes.Buffer)
		for i := 0; i < n; i++ {
			switch ch := expr[i]; ch {
			case '\\':
				if i+1 < n && expr[i+1] == ',' {
					sb.WriteByte(',')
					i++
				} else {
					sb.WriteByte(ch)
				}
			case ',':
				args = append(args, strings.TrimSpace(sb.String()))
				sb.Reset()
			default:
				sb.WriteByte(ch)
			}
		}
		args = append(args, strings.TrimSpace(sb.String()))
	}
	return
}
