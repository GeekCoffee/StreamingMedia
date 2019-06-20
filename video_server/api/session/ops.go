package session

import (
	"../dbops"
	"../defs"
	"../utils"
	"log"
	"sync"
	"time"
)

var sessionMap *sync.Map  // 安全线程操作，用于低水平程序的并发读写


func init(){
	sessionMap = &sync.Map{}
}

//从DB读取session数据到cache
func LoadSessionFromDB(){
	m, err := dbops.RetrieveAllSessions()
	if err != nil{
		return
	}

	//使用range函数进行key-value的遍历，即map的遍历
	//目录是把map中的key-value遍历出来存入sessionMap
	//传入一个匿名函数，函数是可以当做参数来使用的
	m.Range(func(k, v interface{}) bool{
		ss := v.(*defs.Session)    //进行一次类型转换，把v转换为*Session，即v是session对象的引用
		sessionMap.Store(k, ss)
		return true   // if return false, range will stop
	})

}


//用于生成session对象，然后分别存入缓存Map和数据库mysql中
func GenerateNewSessionId(loginName string) string {
	// 通过UUID算法生成session_id
	sessionId, err:= utils.NewUUID()
	if err != nil{
		return err.Error()
	}

	//创建session的时间
	creatTime := NowInMilli()   //得到当前的时间戳，精确到毫秒就行
	ttl := creatTime + 30 * 60 * 1000000  //未来的某个时间点，TTL的过期时间为30分钟，转变为毫秒的时间戳

	//New一个Session对象，并存入cache和DB
	s := &defs.Session{LoginName:loginName, TTL:ttl}
	sessionMap.Store(sessionId, s)
	if err := dbops.InsertSession(sessionId, ttl, loginName); err != nil{
		return err.Error()
	}
	return sessionId
}

//传入sid判断session是否过期
//返回true说明session过期，false说明session没有过期
func IsSessionExpired(sid string) (string, bool){

	currentTime := NowInMilli()   // 得到当前的时间点

	//做TTL的判断
	v, ok := sessionMap.Load(sid)
	s := v.(*defs.Session)
	if ok { //成功取出
		if s.TTL < currentTime { // 若当前时间点大与TTL，说明session已经过期
			//TODO delete session
			return "delete expired session already", true
		}else{
			//session没有过期，返回session体中的loginName即可
			deleteExpiredSession(sid)
			return s.LoginName, false
		}
	}else{
		return "sessionMap.Load(sid) error ", true  //取出错误，默认是session过期
	}
}

//从cache和DB层面删除session
func deleteExpiredSession(sid string){
	//从Map中删除session
	sessionMap.Delete(sid)

	//从DB中删除session
	if err := dbops.DeleteSession(sid); err != nil{
		log.Printf("DeleteSession(sid) occur error: %v", err.Error())
	}
}

//得到当前时间的Unix时间戳
func NowInMilli() int64{
	return time.Now().UnixNano() / 1000000
}