package main

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

var userList []string

type lotteryController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}
func main() {
	app := newApp()
	userList = []string{} //initialize userList

	app.Run(iris.Addr(":8080"))
}

func (c *lotteryController) Get() string {
	count := len(userList)
	return fmt.Sprintf("Users currently participating in the lottery: %d\n", count)
}
