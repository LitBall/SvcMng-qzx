package main

/* vscode调试步骤：
1. 终端输入：dlv debug --headless --listen=:2345 --log --api-version=2
2. launch.json 添加：
{
            "name": "HAHAHA", // 配置名字，随便取
            "type": "go",
            "request": "launch",	// 这个地方存疑，待会讨论
            "mode": "remote",
            "remotePath": "${fileDirname}",
            "port": 2345,			// 调试服务端口
            "host": "127.0.0.1",	// 调试服务IP
            "program": "${fileDirname}",
            "apiVersion": 2,
        }
3. 开始调试
*/
import (
	api "ServiceManagement/api"
	service "ServiceManagement/pkg"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB
var userInfo map[string]string

func main() {
	// 初始化变量
	var key = 0 // 执行操作key
	var userInfo = map[string]string{
		"uid":      "1",
		"username": "gigi",
		"password": "qzx1225",
	}

	// goto-label: ToLogin
ToLogin:
	// 下为登录流程，已优化

	// fmt.Println("请选择操作：\n1. 登录\n2. 注册")
	// fmt.Scanln(&key)
	// Db = service.ConnectDatabase()
	// switch {
	// case key == 1:
	// 	// 用户登录
	// 	userInfo = service.Login()
	// 	if userInfo == nil {
	// 		goto ToLogin
	// 	} else {
	// 		fmt.Println("当前用户：", userInfo["username"])
	// 		break
	// 	}
	// case key == 2:
	// 	// 用户注册
	// 	userInfo = service.Signup()
	// 	if userInfo == nil {
	// 		goto ToLogin
	// 	} else {
	// 		fmt.Println("当前用户：", userInfo["username"])
	// 		break
	// 	}
	// }

	uid, err := strconv.Atoi(userInfo["uid"])
	if err != nil {
		fmt.Println("uid get error:", err)
		return
	}

	key = 0
	// test
	// _ = api.GetHwFlavorDetails("s6.small.1", "cn-south-1", "NGASOM7EFQJNV6Y6RG8R", "spr4yDSO4ezhKjCCHGn3pVhyRq1IB5b4E7MJnNtx")
	_ = api.GetAliFlavorDetails("ecs.n4.small", "cn-hangzhou")
	api.Start()
	// over
	fmt.Println("请选择操作：\n1. 创建实例\n2. 管理密钥\n3. 查看模板")
	fmt.Scanln(&key)
	switch {
	case key == 1:
		// 进入创建实例流程
		if userInfo == nil {
			goto ToLogin
		}
		if err := service.CreateNewInstance(uid); err != nil {
			fmt.Println("部署失败，请重试！T.T")
		}

	case key == 2:
		// 管理密钥
		if userInfo == nil {
			goto ToLogin
		}
		service.ManageAccessKey(userInfo)
	case key == 3:
		// 查看模板
		if userInfo == nil {
			goto ToLogin
		}

	}
}
