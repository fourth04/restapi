package controllers

import (
	"fmt"

	"github.com/kataras/iris/mvc"

	"dnsutils"
)

// UrlsafeController is our sample controller
// it handles GET: /urlsafe and GET: /urlsafe/{name}
type UrlsafeController struct {
	mvc.C
}

// Get is Get Method Controller
// Demos:
// curl -i http://localhost:8080/urlsafe?url=baidu.com
// curl -i http://localhost:8080/urlsafe?url=www.baidu.com
func (c *UrlsafeController) Get() {
	url := c.Ctx.URLParam("url")
	username, _, _ := c.Ctx.Request().BasicAuth()
	rv := make(map[string]interface{})
	rv["url"] = url
	loggerMessage := ""
	if result, err := dnsutils.UrlLibDetect(url); err != nil {
		rv["items"] = map[string]interface{}{
			"evilClass": 0,
			"evilType":  0,
			"level":     0,
			"urlType":   0,
			"url":       url,
		}
		rv["success"] = false
		rv["error"] = err
		loggerMessage = fmt.Sprintf("%s|%s|%v|%v|%v|%v|%v", url, username, false, 0, 0, 0, 0)
	} else {
		rv["items"] = map[string]interface{}{
			"evilClass": result.EvilClass,
			"evilType":  result.EvilType,
			"level":     result.Level,
			"urlType":   result.UrlType,
			"url":       result.Url,
		}
		rv["success"] = true
		loggerMessage = fmt.Sprintf("%s|%s|%v|%v|%v|%v|%v", url, username, true, result.UrlType, result.EvilClass, result.EvilType, result.Level)
	}
	c.Ctx.JSON(rv)
	c.Ctx.Values().Set("logger_message", loggerMessage)
}
