package process

import (
	"LearnGo/chatroom/client/utils"
	"LearnGo/chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
	Conn net.Conn
	Tf   *utils.Transfer
}

func (this *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	defer conn.Close()
	// 2. 准备发送登录消息
	var mes message.Message
	mes.Type = message.RegisterMesType
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName
	// 3. 序列化registerMes
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		conn.Close()
		return
	}
	//把data赋值给mes.Data 字段
	mes.Data = string(data)
	//将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	tf := &utils.Transfer{
		Conn: conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息出错：err=", err)
		return
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功")
		os.Exit(0)
	} else {
		// 登录失败
		fmt.Println("登录失败:", registerResMes.Error)
		os.Exit(0)
	}
	return
}

// Login 写一个函数，完成登录
func (this *UserProcess) Login(userId int, userPwd string, userName string) (err error) {
	// 1. 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	// 2. 准备发送登录消息
	var mes message.Message
	mes.Type = message.LoginMesType
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	loginMes.UserName = userName
	// 3. 序列化loginMes
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		conn.Close()
		return
	}
	// 4. 构建完整消息
	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		conn.Close()
		return
	}
	// 5. 发送消息长度
	var pkgLen uint32 = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)

	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write err=", err)
		conn.Close()
		return
	}
	// 6. 发送消息体
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write err=", err)
		conn.Close()
		return
	}
	// 7. 创建Transfer并读取响应
	tf := &utils.Transfer{Conn: conn}
	// 8. 等待服务器响应
	resMes, err := tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg err=", err)
		conn.Close()
		return
	}
	// 9. 解析响应
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(resMes.Data), &loginResMes)
	if err != nil {
		fmt.Println("解析响应失败:", err)
		conn.Close()
		return
	}
	if loginResMes.Code == 200 {
		// 登录成功，保存连接
		this.Conn = conn
		this.Tf = tf
		fmt.Println("\n========== 登录成功！ ==========")
		// 启动消息接收goroutine
		go this.keepReading()

		//可以显示当前在线用户的列表，遍历loginResMes.UsersId
		fmt.Println("当前在线用户列表如下:")
		for _, v := range loginResMes.UsersId {
			//如果我们要求不显示自己在线
			if v == userId {
				continue
			}
			fmt.Println("用户id:\t", v)
		}
		go serverProcessMes(conn)
		// 显示菜单
		for {
			this.ShowMenu()
		}
	} else {
		// 登录失败
		fmt.Println("登录失败:", loginResMes.Error)
		conn.Close()
		return fmt.Errorf(loginResMes.Error)
	}
	return nil
}

// keepReading 持续读取服务器消息
func (this *UserProcess) keepReading() {
	fmt.Println("启动消息接收goroutine...")
	defer func() {
		if this.Conn != nil {
			this.Conn.Close()
			fmt.Println("连接已关闭")
		}
	}()
	for {
		if this.Tf == nil || this.Conn == nil {
			fmt.Println("连接已断开")
			break
		}
		mes, err := this.Tf.ReadPkg()
		if err != nil {
			fmt.Println("读取服务器消息失败:", err)
			break
		}
		fmt.Printf("收到服务器消息: 类型=%d\n", mes.Type)
		// 这里可以处理不同类型的服务器消息
	}
	fmt.Println("退出消息接收循环")
	os.Exit(0)
}
