package dbops

// 对session进行操作
import(
	"../defs"
	"../utils"
	"database/sql"
	"log"
	"strconv"
	"sync"
)

// 把session插入数据库
func InsertSession(sid string, ttl int64, loginName string) error{

	// 获取session_id,使用UUID算法生成唯一值
	sid, err := utils.NewUUID()
	if err != nil{
		return err
	}

	//转换ttl类型
	ttlstr := strconv.FormatInt(ttl, 10)
	if err != nil{
		return err
	}

	//预编译SQL
	stmtIns, err := DBConn.Prepare("INSERT INTO sessions (session_id, TTL, login_name) VALUES (?,?,?)")
	if err != nil{
		return err
	}

	//执行SQL
	if _, err = stmtIns.Exec(sid, ttlstr, loginName); err != nil{
		return err
	}



	//顺利执行完毕，关闭连接资源
	defer stmtIns.Close()
	return nil
}


// 把session从DB中提出
func RetrieveSession(sid string) (*defs.Session, error){
	//定义一个空的Session对象，并拿到其引用
	s := &defs.Session{}

	// 预编译SQL
	stmtOut, err := DBConn.Prepare("SELECT TTL, login_name FROM sessions WHERE session_id = ?")
	if err != nil{
		return nil, err
	}

	var ttlstr string
	var loginName string

	//执行SQL，查询操作，返回结果集
	err = stmtOut.QueryRow(sid).Scan(&ttlstr, &loginName)    //与SQL的select TTL,login_name的顺序对应
	if err != nil && err != sql.ErrNoRows{
		return nil, err
	}

	var ttl int64
	ttl, err = strconv.ParseInt(ttlstr, 10, 64)  // base是进制数，bitSize是64位
	if err != nil{
		return nil, err
	}

	// 顺序执行完毕
	s.LoginName = loginName
	s.TTL = ttl

	defer stmtOut.Close()
	return s, nil
}


// 返回异步Map，key是session_id，value是session的实体
func RetrieveAllSessions()(*sync.Map, error){
	m := &sync.Map{}  // 申请一个系统的Map引用，已经存在的Map，不需要再重新生成

	stmtOut, err := DBConn.Prepare("SELECT session_id, TTL, login_name FROM sessions")
	if err != nil{
		return nil, err
	}

	// 执行SQL，得到所有的相关rows,再对rows进行处理
	rows, err := stmtOut.Query()
	if err != nil{
		return nil, err
	}

	//迭代每一行数据
	for rows.Next(){
		var sessionId string
		var ttlstr string
		var loginName string
		if err := rows.Scan(&sessionId, &ttlstr, &loginName); err != nil{
			log.Printf("RetrieveAllSession is error: %s", err)
			return nil, err
		}

		//存入Go自带的异步Map中
		if ttl, err := strconv.ParseInt(ttlstr, 10, 64); err == nil{
			s := &defs.Session{
				LoginName:loginName,
				TTL:ttl,
			}
			m.Store(sessionId, s)
		}else{
			log.Printf("RetrieveAllSession is error: %s", err)
			return nil, err
		}
	}

	defer stmtOut.Close()
	return m,nil
}

// 删除session
func DeleteSession(sid string) error {
	//预编译SQL
	stmtDel, err := DBConn.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil{
		return nil
	}

	//执行删除SQL
	if _, err = stmtDel.Exec(sid); err != nil{
		return err
	}

	defer stmtDel.Close()
	return nil
}


