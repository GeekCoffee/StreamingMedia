package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// 定义全局变量，包内可通用
var(
	DBConn *sql.DB  //数据库操作句柄
	Err error  //数据库操作时错误
)

// init方法默认在包内第一个被调用的函数
func init() {
	// open并不是真正连接到了数据库
	// 指定连接数据库的方式是TCP方式
	// := 赋值方式，只在当前方法内有效
	DBConn, _ = sql.Open("mysql", "root:abc5518988@tcp(localhost:3306)/streaming_media?charset=utf8")


	if err := DBConn.Ping(); err != nil{
		panic("连接数据库时出现错误... "+err.Error())
	}
}
