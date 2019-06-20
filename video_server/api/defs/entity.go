package defs

// 给UserName打一个tag，tag名是user_name，tag的属性是json。用于序列化和反序列化
// 还有protobuf属性

// User的实体- data model
type UserCredential struct{
	UserName string `json:"user_name"`
	Password string `json:"Password"`
}


// 视频的data model实体
type VideoInfo struct{
	Vid string `json:"vid"`
	AuthorId int `json:"author_id"`
	VideoName string `json:"video_name"`
	DisplayCtime string `json:"display_ctime""`  // create_time在数据库层面创建，真正的入库时间
}


// 评论comment的实体
type Comment struct{
	CommentId string  // 评论ID
	AuthorName string  // 通过author_id查询users表得到login_name的
	VideoId string  // 视频ID
	Content string // 评论内容
}

// session实体
type Session struct{
	LoginName string   // 用户登录名
	TTL int64
}


