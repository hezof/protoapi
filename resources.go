/*
xml语法:
<?xml version="1.0" encoding="utf-8" ?>
<!DOCTYPE resources [

	<!ELEMENT resource (Code, Message, Status-Code)>
	<!ATTLIST accept-language CDATA "">
	<!ELEMENT code (#PCDATA)>
	<!ELEMENT name (#PCDATA)>
	<!ELEMENT message (#PCDATA)>
	<!ELEMENT status (#PCDATA)>

]>
<!-- accept-language使用 iso_language_code或iso_language_code-ISO_COUNTRY_CODE, 多值用逗号分割 -->
<resources accept-language="en,en-US,en-UK">

	<resource>
	    <!-- 必需: 错误代码 -->
	    <code>1001</code>
	    <!-- 可选: 错误名称 -->
		<name>test</name>
	    <!-- 可选: 错误消息 -->
	    <message>测试%v</message>
	    <!-- 可选: 状态码 -->
	    <status>403</status>
	</resource>

</resources>
*/
package protoapi

import (
	"bytes"
	"encoding/xml"
	"github.com/hezof/base"
	"github.com/hezof/log"
	"os"
	"path/filepath"
	"strings"
)

type resource struct {
	Status  uint32 `xml:"status"`
	Code    uint32 `xml:"code"`
	Name    string `xml:"name"`
	Message string `xml:"message"`
}

type resources struct {
	XMLName        xml.Name    `xml:"resources"`
	AcceptLanguage string      `xml:"accept-language,attr"`
	Resources      []*resource `xml:"resource"`
}

/*
注意: 为了性能, 没有锁机制! 不能并发执行InitResultBundle()与LoadResultBundle()
*/
var (
	allResMap = make(map[string]map[uint32]*resource)
	defResMap = make(map[uint32]*resource)
	hasResMap = false // 标记是否存在资源文件, 用于快速判断
)

func InitResourceBundle(resDir, defLang string) error {

	dir, _ := core.LocatePath(resDir)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		// 加载所有xml文件
		if info != nil && !info.IsDir() && strings.HasSuffix(path, ".xml") {

			log.Info("load resource file: %v", path)

			ls, bs, er := ReadResourceConfig(path)
			if er != nil {
				return er
			}
			for _, l := range ls {
				// 相同语言可能覆盖,以最后加载为准! 需要各xml的accept-language正交!
				all, ok := allResMap[l]
				if ok {
					for k, v := range bs {
						all[k] = v
					}
				} else {
					all = bs
				}
				allResMap[l] = all
			}
		}
		return nil
	})
	if err == nil {
		defResMap = allResMap[defLang]
		hasResMap = len(allResMap) > 0
	}
	return err
}

func ReadResourceConfig(path string) (langs []string, bundle map[uint32]*resource, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	rss := new(resources)
	err = xml.Unmarshal(data, rss)
	if err != nil {
		return
	}
	langs = strings.Split(rss.AcceptLanguage, ",")
	bundle = make(map[uint32]*resource, len(rss.Resources))
	for _, rs := range rss.Resources {
		bundle[rs.Code] = rs
	}
	return
}

func LoadResourceBundle(code uint32, languages ...string) (uint32, string, string, bool) {
	for _, l := range languages {
		if bds := allResMap[l]; bds != nil {
			if bd := bds[code]; bd != nil {
				return bd.Status, bd.Name, bd.Message, true
			}
		}
	}
	if bd := defResMap[code]; bd != nil {
		return bd.Status, bd.Name, bd.Message, true
	}
	return 0, "", "", false
}

// 根据Accept-Language快速获取(从左到右,不按q排序). 该方法性能优于parseAcceptLanguage!
func fastGetResMapByAcceptLanguage(acceptLanguage string) map[uint32]*resource {
	// 语法Accept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2
	var data = core.UnsafeBytes(acceptLanguage)
	var temp []byte
	var idx int
	for len(data) > 0 {
		idx = bytes.IndexByte(data, ',')
		if idx == -1 {
			temp = data
			data = nil // 下轮结束
		} else {
			temp = data[:idx]
			data = data[idx+1:]
		}
		idx = bytes.IndexByte(temp, ';')
		if idx > 0 {
			temp = temp[:idx]
		}
		if ret := allResMap[core.UnsafeString(temp)]; ret != nil {
			return ret
		}
	}
	return defResMap
}
