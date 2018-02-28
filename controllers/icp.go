package controllers

import (
	"fmt"

	"dnsutils"

	"github.com/go-redis/redis"

	"github.com/kataras/iris/mvc"
	"golang.org/x/net/publicsuffix"
)

// IcpController is our sample controller
// it handles GET: /icp and GET: /icp/{name}
type IcpController struct {
	mvc.C
	Client *redis.Client
}

// Get is a Get method controller
// Demos:
// curl -i http://localhost:8080/icp?icp=baidu.com
// curl -i http://localhost:8080/icp?icp=www.baidu.com
// 返回结果说明：当success返回为false时，说明，尝试了很多次查询，但是云端接口都为null
func (c *IcpController) Get() {
	// domain := c.Ctx.FormValue("domain")
	domain := c.Ctx.URLParam("domain")
	username, _, _ := c.Ctx.Request().BasicAuth()
	rv := make(map[string]interface{})
	rv["domain"] = domain
	tldPlusOne, _ := publicsuffix.EffectiveTLDPlusOne(domain)
	if tldPlusOne != "" {
		if result, err := c.Client.HGetAll(tldPlusOne).Result(); err != nil || len(result) == 0 {
			resultCloud, errCloud := dnsutils.IcpQuery(tldPlusOne, 20)
			if errCloud != nil {
				rv["success"] = false
				rv["items"] = map[string]string{
					"domain":           domain,
					"name":             "",
					"companyBeianCode": "",
					"siteBeianCode":    "",
					"status":           "1",
				}
				rv["error"] = err
			} else {
				domain := resultCloud["domain"].(string)
				items := resultCloud["items"].(map[string]interface{})
				siteBeianCode := items["siteBeianCode"].(string)
				if siteBeianCode != "" {
					c.Client.HMSet(domain, items)
				}
				resultCloud["success"] = true
				rv = resultCloud
			}
		} else {
			rv["success"] = true
			rv["items"] = result
		}
	} else {
		rv["success"] = "true"
		rv["items"] = map[string]string{
			"domain":           domain,
			"name":             "",
			"companyBeianCode": "",
			"siteBeianCode":    "",
			"status":           "1",
		}
		rv["error"] = "domain format error"
	}
	c.Ctx.JSON(rv)
	loggerMessage := fmt.Sprintf("%s|%s|%v", domain, username, rv["success"])
	c.Ctx.Values().Set("logger_message", loggerMessage)
}
