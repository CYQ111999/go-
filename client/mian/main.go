package main

import (
	"LearnGo/chatroom/client/process"
	"fmt"
	"os"
)

// 定义两个变量，一个表示用户id，一个表示用户密码
var userId int
var userPwd string
var userName string

func main() {
	//接受用户的选择
	var key int
	//判断是否还继续显示菜单
	//var loop = true
	for {
		fmt.Println("---------------欢迎登陆聊天系统---------------")
		fmt.Println("\t\t\t 1 登录聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择(1-3)")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户的id号:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的密码:")
			fmt.Scanf("%s\n", &userPwd)
			//完成登录
			//1.创建一个UserProcess的实例
			up := &process.UserProcess{}
			up.Login(userId, userPwd, userName)
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户的id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的密码:")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户的昵称:")
			fmt.Scanf("%s\n", &userName)
			//调用UserProcess, 完成注册请求
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("你的输入有误，请重新输入")
		}
	}

	//根据用户的输入显示新的提示信息
	//if key == 1 {
	//	//说明用户要登录
	//	//fmt.Println("请输入用户的id")
	//	//fmt.Scanf("%d\n", &userId)
	//	//fmt.Println("请输入用户的密码")
	//	//fmt.Scanf("%s\n", &userPwd)
	//	//先把登录函数写到另一个文件,比如先写到login
	//	//login(userId, userPwd)
	//	//if err != nil {
	//	//	fmt.Println("登陆失败")
	//	//} else {
	//	//	fmt.Println("登陆成功")
	//	//}
	//} else if key == 2 {
	//
	//}
}
