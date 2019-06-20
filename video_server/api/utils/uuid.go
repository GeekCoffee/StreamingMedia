package utils

import (
	"crypto/rand"
	"io"
	"fmt"
)

// uuid可为视频实体创建id，为comment创建comment_id也需要用
func NewUUID() (string, error){
	uuid := make([]byte, 16) // 默认一个uuid为16个字节
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil{
		return "", err
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	//version 4 (pseudo-random)
	uuid[6] = uuid[6]&^0xf0 | 0x40

	// Springf返回的是一个格式化后的string, Fprinf是把字节流写入w，Prinf是打印到显示屏
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}