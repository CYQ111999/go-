package main

import (
	"LearnGo/chatroom/server/model"
	"fmt"
	"net"
)

// 处理和客户端的通讯
func process(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("close conn failed, err:%v\n", err)
			return
		}
	}(conn)
	processor := &Processor{
		Conn: conn,
	}
	err := processor.Process2()
	if err != nil {
		fmt.Println("客户端和服务器端通讯协程错误err=", err)
		return
	}
}

// 这里我们编写一个函数，完成对UserDao的初始化任务
func initUserDao() {
	model.MyUserDao = model.NewUserDao(pool)

}

func main() {
	//服务器启动时我们就初始化连接池
	initPool("localhost:6379", 16, 0, 300)
	initUserDao()
	fmt.Println("服务器在8889端口监听....")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			fmt.Printf("close listener failed, err:%v\n", err)
			return
		}
	}(listen)
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	//一旦监听成功，就等待用户连接
	for {
		fmt.Println("等待客户端来连接服务器....")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("net.Listen err=", err)
		}
		//一旦链接成功，就启动一个协程和客户端保持通讯
		go process(conn)
	}
}
