package main

import (
	"fmt"
	"strings"

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

// HomePage
func (c *lotteryController) Get() string {
	count := len(userList)
	return fmt.Sprintf("Users currently participating in the lottery: %d\n", count)
}

// Importing User List
// POST http://localhost:8080/import
// params: user
func (c *lotteryController) PostImprt() string {
	strUsers := c.Ctx.FormValue("users")
	users := strings.Split(strUsers, ",")
	count1 := len(userList)
	for _, u : range users {
		u = strings.TrimSpace(u)
		if len(u) > 0 {
			userList = append(userList, u)
		}
	}
	count2 := len(userList)
	return fmt.Sprintf("Users participating in the lottery %d, Users successfully imported: %d\n", count2, (count2 - count1));
}

// Lottery Executing
// 
func (c *lotteryController) GetLucky() string {
	count := len(userList)
	if count > 1 {
		// Creating time_stamp for generating a random number
		seed := time.Now().UnixNano() 
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count))
		user := userList[index]
		userList = append(userList[0:index], userList[index + 1:]...)

		return fmt.Sprintf("The current winning user is: %s; Number of remaining users: %d\n", user, count - 1);
	} else if count == 1 {
		user := userList[0];
		return fmt.Sprintf("The current winning user is: %s; Number of remaining users: %d\n", user, count - 1);
	} else {
		return fmt.Sprintf("There are no user left in the room");
	}
}
