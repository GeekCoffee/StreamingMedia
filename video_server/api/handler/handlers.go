package handler

import (
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params)  {
	io.WriteString(w, " Create a user 陈圣" )
}


func UserLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params){
	userName := p.ByName("user_name")
	io.WriteString(w, userName)
}

// test init()方法
//var I int
//
//func init(){
//	I = 999
//	fmt.Printf("i = %d\n", I)
//}





