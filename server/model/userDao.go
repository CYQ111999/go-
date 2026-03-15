package model

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

var (
	MyUserDao *UserDao
)

type UserDao struct {
	pool *redis.Pool
}

func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// 注意：这里返回 *User，不是 *message.User
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &User{} // 使用 model.User
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json unmarshal err:", err)
		return
	}
	return
}

func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

func (this *UserDao) Register(user *User) (err error) {
	conn := this.pool.Get()
	defer conn.Close()

	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("json marshal err:", err)
		return
	}

	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误", err)
		return
	}
	return
}
