package dbops

// 测试数据库操作函数
// 测试api.go文件中的函数
// 注意程序运行逻辑：init(dblogin, truncate tables) -> run tests->clear data(truncate table)

import(
	"fmt"
	"testing"
	"strconv"
	"time"
)

// go的test中现在还不能函数中传递参数，所以用一个全局变量tempvid来传递vid
var tempvid string

func clearTables(){
	// truncate table_name, 这句SQL是迅速清空表内的所有资料，并且对有自增值字段的值，做重归零操作
	DBConn.Exec("truncate users")
	DBConn.Exec("truncate video_info")
	DBConn.Exec("truncate comments")
	DBConn.Exec("truncate sessions")
}

// 测试版主函数
func TestMain(m *testing.M) {
	clearTables()   //先做一次clearTables，清空所有表内资料
	m.Run()       // 调用一下TestUserWorkFlow
	clearTables() // 然后清空表，回归到最原始状态
}

// 定义测试用户工作流
func TestUserWorkFlow(t *testing.T){
	t.Run("CreateUser", TestCreateUserCredential)
	t.Run("GetUser", TestGetUserCredential)
	t.Run("DelUser", TestDeleteUser)
	t.Run("RegetUser", TestRegetUser)
}

// 命名规范，需要前面写一个Test大写
// 测试addUser用户函数
func TestCreateUserCredential(t *testing.T) {
	err := CreateUserCredential("williamchen", "123")
	if err != nil{
		t.Errorf("CreateUser error: %v", err)
	}
}

// 测试读取用户信息函数
func TestGetUserCredential(t *testing.T) {
	pwd, err := GetUserCredential("williamchen")
	if pwd != "123" || err != nil{
		t.Errorf("GetUser error: %v", err)
	}
}

// 测试删除用户函数
func TestDeleteUser(t *testing.T) {
	err := DeleteUser("williamchen", "123")
	if err != nil{
		t.Errorf("DeleteUser error: %v", err)
	}
}

// 测试删除用户函数是否生效
func TestRegetUser(t *testing.T){
	pwd, err := GetUserCredential("williamchen")
	if pwd != "" || err != nil{
		t.Errorf("RegetUser error: %v", err)
	}
}

// 视频实体Test用例
func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("PrepareUser", TestCreateUserCredential)  // 需要先有用户才能有视频，因为视频是用户上传的，属于依赖关系
	t.Run("AddVideo", TestAddVideoInfo)
	t.Run("GetVideo", TestGetVideoInfo)
	t.Run("DelVideo", TestDeleteVideoInfo)
	t.Run("RegetVideo", TestRegetVideoInfo)
}



func TestAddVideoInfo(t *testing.T){
	videoEntity, err := AddNewVideo(1, "my-video")  //用户id和视频name
	if err != nil{
		t.Errorf("Error of AddVideoInfo: %v", err)
	}

	// vid是使用uuid生成的，不是数据表的auto_increment自增的，但它是primary key，主键
	tempvid = videoEntity.Vid
}

func TestGetVideoInfo(t *testing.T){
	_, err := GetVideoInfo(tempvid)
	if err != nil{
		t.Errorf("Error of GetVideoInfo: %v", err)
	}
}

func TestDeleteVideoInfo(t *testing.T) {
	if err := DeleteVideoInfo(tempvid); err != nil{
		t.Errorf("testDeleteVideoInfo error: %v", err)
	}
}

func TestRegetVideoInfo(t *testing.T){
	res, err := GetVideoInfo(tempvid)
	if err != nil && res != nil{
		t.Errorf("error of testRegetVideoInfo: %v", err)
	}
}



// 设计评论Test的工作流
func TestCommentsFlow(t *testing.T){
	clearTables()   // truncate一下所有表
	t.Run("PrepareUser", TestCreateUserCredential) // 需要先有用户
	//t.Run("PrepareVideo", TestAddVideoInfo)  // 然后才有用户下的视频 ，使用变量方式实现vid
	t.Run("AddNewComment", TestAddNewComment) // 最后在视频下进行评论
	t.Run("ListComments", TestListComments) // 获取某个视频下一段时间内的评论
}

// 添加一条评论
func TestAddNewComment(t *testing.T) {
	var vid, authorId, content  = "1", 1, "this is a good video"
	if err := AddNewComment(vid, authorId, content); err != nil{
		t.Errorf("Error of TestAddNewComment: %v", err)
	}
}

// 根据时间段取出所评论的list
func TestListComments(t *testing.T) {
	vid := "1"
	from := 1514764800   //大概是当前时间的前两周的时间戳

	// Now()返回Time类，UnixNano()返回当前时间戳，int64
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000,10))
	list, err := ListComments(vid, from, to)
	if err != nil{
		t.Errorf("Error of TestListComments:%v", err)
	}

	// 循环打印输出到console看看
	// i从0开始
	for i, c := range list{
		fmt.Printf("================== %d   %v",i, c)
	}
}


