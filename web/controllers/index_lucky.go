package controllers

import "github.com/Zawa-ll/raffle/comm"

// localhost:8080/lucky
func (c *IndexController) GetLucky() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	// 1 Authenticate the logged-in user
	loginuser := comm.GetLoginUser(c.Ctx.Request())
	if loginuser == nil || loginuser.Uid < 1 {
		rs["code"] = 101
		rs["msg"] = "Please log in first, then come to the lottery"
		return rs
	}
	ip := comm.ClientIP(c.Ctx.Request())
	api := &LuckyApi{}
	code, msg, gift := api.luckyDo(loginuser.Uid, loginuser.Username, ip)
	rs["code"] = code
	rs["msg"] = msg
	rs["gift"] = gift
	return rs
}
