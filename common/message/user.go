package message

// User 定义一个用户结构体
type User struct {
	//为了序列化和反序列化成功，我们要保证用用户信息的json字符串的key
	//和结构体字段对应的tag名字一致
	UserId     int    `json:"userId"`
	UserPwd    string `json:"userPwd"`
	UserName   string `json:"userName"`
	UserStatus int    `json:"userStatus"` //用户状态
	Sex        string `json:"sex"`
}
