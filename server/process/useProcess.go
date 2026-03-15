package process

import (
	"LearnGo/chatroom/common/message"
	"LearnGo/chatroom/server/model"
	"LearnGo/chatroom/server/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//增加一个字段，表示Conn是哪个用户
	UserId int
}

// NotifyOthersOnlineUser 这里我们编写一个通知所有在线用户的方法
// userId 要通知其他在线用户
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {
	//遍历onlineUsers, 然后一个一个的发送
	for id, up := range userMgr.onlineUsers {
		if id == userId {
			continue
		}
		//开始通知
		up.NotifyMeOnline(userId)
	}
}

func (this *UserProcess) NotifyMeOnline(userId int) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	//将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//将序列化后的notifyUserSatatusMes赋值给mes.Data
	mes.Data = string(data)

	//对mes再次序列化，准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//发送，创建我们Transfer实例，发送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}
}

func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	// 1. 反序列化注册消息
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return
	}

	// 添加调试信息
	fmt.Printf("收到注册请求: UserId=%d, UserName=%s\n",
		registerMes.User.UserId, registerMes.User.UserName)

	// 2. 准备响应消息
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	// 3. 创建 model.User 对象
	user := &model.User{
		UserId:   registerMes.User.UserId,
		UserPwd:  registerMes.User.UserPwd,
		UserName: registerMes.User.UserName,
	}

	// 4. 调用DAO层注册
	err = model.MyUserDao.Register(user)
	if err != nil {
		fmt.Printf("注册失败: %v\n", err)
		if errors.Is(err, model.ERROR_USER_EXISTS) {
			registerResMes.Code = 505
			registerResMes.Error = "用户已存在"
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册失败: " + err.Error()
		}
	} else {
		registerResMes.Code = 200
		registerResMes.Error = ""
		fmt.Printf("用户 %s 注册成功\n", user.UserName)
	}

	// 5. 序列化响应
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 6. 构建完整消息
	resMes.Data = string(data)
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 7. 发送响应
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

// ServerProcessLogin 编写一个函数ServerProcessLogin函数，专门处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//1.先从mes中取出data，并反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	//1.先声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	//2.再声名一个LoginResMes,并完成赋值
	var loginResMes message.LoginResMes
	//1.使用model.MyUserDao 到redis去验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		loginResMes.Code = 500
		loginResMes.Error = "该用户不存在，请注册再使用...!!!!!"
	} else {
		loginResMes.Code = 200
		//这里因为用户已经登录成功。我们应该把登录成功的用法放到userMrg中
		//将登陆成功的用户userId 赋给this
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)
		//通知其他用户我上线了
		this.NotifyOthersOnlineUser(loginMes.UserId)
		//将当前在线的用户Id放入loginResMes.UsersId
		//遍历 userMgr.onlineUser
		for id := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println(user, "登陆成功")
	}
	//3.将loginResMes 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//4.将data赋值给resMes
	resMes.Data = string(data)
	//5.对resMes序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//6.发送data将其封装到write函数
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
