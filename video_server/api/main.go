package main

import (
	"./handler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)
// RESTful API
// Create = POST , Read = GET, Update = POST,  Delete = DELETE
// list all video for user_name , /users/:user_name/videos, Method=GET, HTTP-statusCode = 200, 400, 500



// 用户注册和登录handler
func RegisterHandler() *httprouter.Router{
	router := httprouter.New()
	router.POST("/user", handler.CreateUser)  //一个handler是一个goroutine，一个goroutine大概是一个4KB大小
	router.POST("/user/:user_name", handler.UserLogin)
	return router
}


// part1: listen->registerHandler->router->真正处理请求的Handler
// part2: handler->validation{ 1.request请求体验证  2.user在数据库层面是否存在 }-> 业务逻辑处理logic-> 封装reponse发回去

func main(){
	//r := RegisterHandler() //得到一个Router
	//log.Fatal(http.ListenAndServe(":8000", r))  // 监听ip:8000端口的请求，把请求转发给r处理

}







