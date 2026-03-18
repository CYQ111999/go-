package process2

import "fmt"

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// 因为UserMgr实例在服务器端只有一个
// 因为在很多地方都会使用到因此定义为全局变量
var (
	userMgr *UserMgr
)

// 完成对userMgr初始化的工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// AddOnlineUser 完成对onlineUsers的添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

// 删除
func (this *UserMgr) delOnlineUser(userId int) {
	delete(this.onlineUsers, userId)
}

// GetAllOnlineUser 返回当前所有在线的用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

// GetOnlineUserById 根据用户id返回对应的值
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	up, ok := this.onlineUsers[userId]
	if !ok {
		//说明要查找的这个用户当前不在线
		err = fmt.Errorf("用户%d 不存在", userId)
		return
	}
	return
}
