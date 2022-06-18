package service

import (
	"errors"
	"fmt"
	"strconv"
)

// P0 注册新用户 done
// --ConnectDatabase
func Signup() map[string]string {
	db = ConnectDatabase()
	usr, pwd, email := "", "", ""
	fmt.Println("开始注册！:3")
	fmt.Println("请输入账号：")
	fmt.Scanln(&usr)
	fmt.Println("请输入密码：")
	fmt.Scanln(&pwd)
	fmt.Println("请输入邮箱：")
	fmt.Scanln(&email)
	fmt.Println("注册中，请稍后……")
	r, err := db.Exec("insert into user(username, password, email, usergroup) values(?, ?, ?, ?);", usr, pwd, email, "user")
	if err != nil {
		fmt.Println("exec failed, ", err)
		return nil
	}
	id, errr := r.LastInsertId()
	if errr != nil {
		fmt.Println("exec failed, ", errr)
		return nil
	}
	fmt.Println("insert succ:", id, " 注册成功！")
	return map[string]string{
		"uid":      strconv.FormatInt(id, 10),
		"username": usr,
		"password": pwd,
	}
}

// P0 用户登录 done
// --ConnectDatabase
func Login() map[string]string {
	usr, pwd := "", ""
	fmt.Println("开始登录！:3")
	fmt.Println("请输入账号：")
	fmt.Scanln(&usr)
	fmt.Println("请输入密码：")
	fmt.Scanln(&pwd)

	var user User
	db = ConnectDatabase()
	rows := db.QueryRow("select uid, username, password from user where username = ? and password = ?", usr, pwd)
	err := rows.Scan(&user.Uid, &user.Username, &user.Password)
	if err != nil {
		// 查找失败
		fmt.Println("登录失败！:(")
		return map[string]string{
			"uid":      "",
			"username": "",
			"password": "",
		}
	} else {
		userInfo := map[string]string{
			"uid":      strconv.Itoa(user.Uid),
			"username": user.Username,
			"password": user.Password,
		}
		fmt.Println("登录成功！:)")
		return userInfo
	}
}

// P0 退出登录 done
func Logout() map[string]string {
	fmt.Println("退出登录成功，欢迎下次光临~:)")
	return map[string]string{
		"uid":      "",
		"username": "",
		"password": "",
	}
}

// P1 更新密钥 done
// --GetKey_io
func UpdateAccessKey(uid int, flag string) error {
	// 1. 已有用户修改密钥 2. 新用户添加密钥
	db = ConnectDatabase()
	switch {
	case flag == "update":
		fmt.Println("开始更新密钥！:3")
		key_pair := GetKey_io(uid)
		r, err := db.Exec("update accesskey set aws_a_key = ?, aws_s_key = ?, ali_a_key = ?, ali_s_key = ?, hw_a_key = ?, hw_s_key = ?, aws_region = ?, ali_region = ?, hw_region = ?, where uid = ?;",
			key_pair.Aws_a_key, key_pair.Aws_s_key, key_pair.Ali_a_key, key_pair.Ali_s_key, key_pair.Hw_a_key, key_pair.Hw_s_key, key_pair.Aws_region, key_pair.Ali_region, key_pair.Hw_region, uid)
		if err != nil {
			fmt.Println("exec failed, ", err)
			return errors.New(err.Error())
		}
		row_aff, err := r.RowsAffected()
		if err != nil {
			fmt.Println("exec failed, ", err)
			return errors.New(err.Error())
		}
		fmt.Print("update succ:", row_aff, "更新成功！")
		// fmt.Println("key_pair:", key_pair)
		break

	case flag == "add":
		fmt.Println("开始创建密钥！:3")
		key_pair := GetKey_io(uid)
		r, err := db.Exec("insert into accesskey(uid, aws_a_key, aws_s_key, ali_a_key, ali_s_key, hw_a_key, hw_s_key) values(?, ?, ?, ?, ?, ?, ?);", uid, key_pair.Aws_a_key, key_pair.Aws_s_key, key_pair.Ali_a_key, key_pair.Ali_s_key, key_pair.Hw_a_key, key_pair.Hw_s_key)
		if err != nil {
			fmt.Println("exec failed, ", err)
			return errors.New(err.Error())
		}
		id, err := r.LastInsertId()
		if err != nil {
			fmt.Println("exec failed, ", err)
			return errors.New(err.Error())
		}
		fmt.Println("insert succ:", id, " 密钥添加成功！")
		break
	}
	return nil
}

// P0 管理密钥 done
// --UpdateAccessKey
func ManageAccessKey(userInfo map[string]string) {
	if userInfo["uid"] == "" {
		fmt.Println("请先登录！")
		return
	}

	fmt.Println("密钥加载中……")
	var (
		accesskey AccessKey
		key       int
	)
	db = ConnectDatabase()
	uid, err := strconv.Atoi(userInfo["uid"])
	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}
	row := db.QueryRow("select ifnull(aws_a_key, '') aws_a_key, ifnull(aws_s_key, '') aws_s_key, ifnull(ali_a_key, '') ali_a_key, ifnull(ali_s_key, '') ali_s_key, ifnull(hw_a_key, '') hw_a_key, ifnull(hw_s_key, '') hw_s_key from accesskey where uid = ?", uid)
	err = row.Scan(&accesskey.Aws_a_key, &accesskey.Aws_s_key,
		&accesskey.Ali_a_key, &accesskey.Ali_s_key,
		&accesskey.Hw_a_key, &accesskey.Hw_s_key)
	if err != nil {
		fmt.Println(err)
		fmt.Println("没有找到当前用户的密钥，是否添加？ 1. 是 2.否")
		fmt.Scanln(&key)
		switch {
		case key == 1:
			UpdateAccessKey(uid, "add")
		case key == 2:
			return
		}
		return
	}
	fmt.Println("密钥加载成功！:)")
	fmt.Println(accesskey)
	fmt.Println("请选择操作:\n1. 修改\n2. 退出")
	fmt.Scanln(&key)
	switch {
	case key == 1:
		UpdateAccessKey(uid, "update")
		row := db.QueryRow("select ifnull(aws_a_key, '') aws_a_key, ifnull(aws_s_key, '') aws_s_key, ifnull(ali_a_key, '') ali_a_key, ifnull(ali_s_key, '') ali_s_key, ifnull(hw_a_key, '') hw_a_key, ifnull(hw_s_key, '') hw_s_key from accesskey where uid = ?", uid)
		err := row.Scan(&accesskey.Aws_a_key, &accesskey.Aws_s_key,
			&accesskey.Ali_a_key, &accesskey.Ali_s_key,
			&accesskey.Hw_a_key, &accesskey.Hw_s_key)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("密钥已更新：\n", accesskey)
	case key == 2:
		fmt.Println("欢迎下次再来~:)")
		return
	}
	return
}
