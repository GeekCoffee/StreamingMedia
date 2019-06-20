package dbops

import (
	"database/sql"
	"fmt"
	"time"
	"../defs"
	"../utils"
	"log"
)

//先调用conn.go中的dbConn全局变量，得到数据库句柄
func CreateUserCredential(loginName string, pwd string) error{
	// 预编译, Ins表示是数据流入DB操作
	stmtIns, err := DBConn.Prepare("INSERT INTO users (login_name, pwd) VALUES (?, ?)")
	if err != nil{
		log.Printf("CreateUser error: %s", err)
		return err
	}

	if _, err = stmtIns.Exec(loginName, pwd); err != nil{
		return err
	}  //把两个问号?的占位符给赋值
	defer stmtIns.Close()
	return nil
}

func GetUserCredential(loginName string) (string, error){
	// 预编译, Out表示数据从DB流出操作
	stmtOut, err := DBConn.Prepare("SELECT pwd FROM users WHERE login_name = ?")
	if err != nil{
		log.Printf("GetUser error: %s", err)
		return "", err
	}

	// 读取loginName的值所在的row，然后读取所在行的pwd值放入变量pwd中
	var pwd string
	// ErrNoRows是一个string变量，返回集中没有row的错误
	if err = stmtOut.QueryRow(loginName).Scan(&pwd); err != nil && err != sql.ErrNoRows{
		return "", err
	}
	defer stmtOut.Close()  //可以使用defer，但是会开辟额外的栈，消耗系统性能
	return pwd, nil
}

func DeleteUser(loginName string, pwd string) error{
	//预编译
	stmtDel, err := DBConn.Prepare("DELETE FROM  users WHERE login_name = ? AND pwd = ?")
	if err != nil{
		log.Printf("DeleteUser error: %s", err)
		return err
	}

	if _, err = stmtDel.Exec(loginName, pwd); err != nil{
		return err
	}
	defer stmtDel.Close() // defer的语句在函数结束后执行，或者在函数结束return前执行
	return nil
}
// 每修改完一次代码，要继续test一下

// 1 1  2 3 5 8 13 21 .....

// 返回VideoInfo实体
// aid = 用户id
// 在author_id下添加一个视频数据
func AddNewVideo(authorId int, videoName string)(*defs.VideoInfo, error){
	// create uuid, UUID是通用唯一识别码, 非人工指定，也非人工去识别,在一定范围内唯一地与某个实体绑定
	// UUID可由16个字符或者32个字符组成，由十六进制表示

	// 使用uuid算法生成vid
	vid, err := utils.NewUUID()
	if err != nil{
		return nil, err
	}

	t := time.Now()  //返回当前时间的time的struct
	ctime := t.Format("Jan 02 2006, 15:04:05")  //标准化时间格式后形成string返回
	fmt.Println("ctime = ", ctime)
	stmtIns, err := DBConn.Prepare(`INSERT INTO video_info(id, author_id, video_name, display_ctime) 
			VALUES (?, ?, ?, ?)`)
	if err != nil{
		return nil, err
	}

	_ ,  err = stmtIns.Exec(vid, authorId, videoName, ctime)
	if err != nil{
		return nil, err
	}

	videoEntity := &defs.VideoInfo{
		Vid:vid,
		AuthorId:authorId,
		VideoName:videoName,
		DisplayCtime:ctime,
	}

	defer stmtIns.Close()

	return videoEntity, nil
}

// 根据vid获取视频信息
func GetVideoInfo(vid string) (*defs.VideoInfo, error){
	stmtOut, err := DBConn.Prepare("SELECT author_id, video_name, display_ctime FROM video_info WHERE id = ?")

	var aid int  //用户id，int
	var dct string  //格式化的显示时间，string
	var name string //视频name，string

	//得到vid所在的这一行的数据，取这行的数据中的author_id、name、display_ctime
	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)

	//为err错误，但不是ErrNoRows空值错误
	if err != nil && err != sql.ErrNoRows{
		return nil, err
	}

	// 返回集为nil
	if err == sql.ErrNoRows {
		return nil, nil
	}

	defer stmtOut.Close()
	res := &defs.VideoInfo{Vid: vid, AuthorId: aid, VideoName: name, DisplayCtime: dct}
	return res, nil
}

// 通过vid删除视频信息
func DeleteVideoInfo(vid string) error{
	stmtDel, err := DBConn.Prepare("DELETE FROM video_info where id = ?")
	if err != nil{
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil{
		return err
	}
	defer stmtDel.Close()
	return nil
}


// 在某视频ID下添加comment
// from,to是时间段的区间表示
func AddNewComment(vid string, authorId int, content string) error{
	cid, err := utils.NewUUID()
	if err != nil{
		return err
	}

	stmtIns, err := DBConn.Prepare("INSERT INTO comments (id, video_id, author_id, content) VALUES (?,?,?,?)")
	if err != nil{
		return err
	}

	_, err = stmtIns.Exec(cid, vid, authorId, content)
	if err != nil{
		return err
	}
	defer stmtIns.Close()
	return nil
}


// 获取comment的list，list获取从from到to时间段的评论
func ListComments(vid string, from, to int) ([]*defs.Comment, error){
	stmtOut, err := DBConn.Prepare(` SELECT comments.id, users.login_name, comments.content FROM comments
		INNER JOIN users ON comments.author_id = users.id 
		WHERE comments.video_id = ? AND comments.create_time > FROM_UNIXTIME(?)
		 AND comments.create_time <= FROM_UNIXTIME(?)`)
	// FROM_UNIXTIME(arg), 参数是UNIX时间戳

	var commentList []*defs.Comment

	if err != nil{
		return commentList, err
	}

	rows, err := stmtOut.Query(vid, from, to)  //得到id、login_name、content组成的二维表结果集
	if err != nil{
		return commentList, err
	}

	//使用迭代器循环遍历结果集
	for rows.Next(){
		var cid, loginName, content string // 评论ID、评论人名称、评论内容

		// 使用Scan方法读取返回集中的对应列的数据，并放入变量中
		// 要对应comment.id, users.login_name, comments.content的顺序
		err := rows.Scan(&cid, &loginName, &content)
		if err != nil{
			return nil, err
		}
		c := &defs.Comment{
			CommentId:cid,
			AuthorName:loginName,
			VideoId:vid,
			Content:content,
		}
		// 把c实体追加到commentList中
		commentList = append(commentList, c)
	}

	defer stmtOut.Close()
	return commentList, nil

}

