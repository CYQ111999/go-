// model/user.go
package model

// User 定义一个用户结构体
type User struct {
	// 注意：json tag 要一致
	// 你原来的 tag 是 "user_id"，但代码中用的是 "userId"
	// 改成一致
	UserId   int    `json:"userId"` // 改为 userId
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}
