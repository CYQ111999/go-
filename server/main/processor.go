package main

import (
	"LearnGo/chatroom/common/message"
	serverprocess "LearnGo/chatroom/server/process"
	"LearnGo/chatroom/server/utils"
	"fmt"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}

// 编写一个ServerProcessMess函数
// 功能：根据客户端发送的消息种类不同，决定调用哪个函数来处理
func (this *Processor) serverProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//处理登录
		up := &serverprocess.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)
		if err != nil {
			fmt.Println("server process login error:", err)
			return
		}
	case message.RegisterMesType:
		up := &serverprocess.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	//处理注册
	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}

// Process2 处理用户和客户端的通讯
func (this *Processor) Process2() (err error) {
	for {
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, e := tf.ReadPkg()
		if e == io.EOF {
			fmt.Println("客户端退出")
			return
		} else if e != nil {
			fmt.Println("readPkg err=", e)
			return
		}
		_ = this.serverProcessMes(&mes)
		err = nil
	}
}
