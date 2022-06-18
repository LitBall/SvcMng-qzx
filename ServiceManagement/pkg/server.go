package service

/*
- 核心逻辑实现：
	- 数据库操作
	- 业务功能
- 注意：
	- 小写字幕开头标识符无法导出被外部访问
	- 查询多行用 db.Query, 单行用 db.Select / db.QueryRow
*/

import (
	"bufio"
	"database/sql"
	_ "encoding/json"
	"errors"
	"fmt"
	_ "io"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// 连接数据库 done
func ConnectDatabase() *sqlx.DB {
	database, err := sqlx.Open("mysql", "root:1005@tcp(127.0.0.1:3306)/multicloud_sm")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		log.Fatal(err)
	}
	db = database
	// fmt.Println("DB import success!")
	return db
	// defer db.Close()
	// 学习defer todo
}

// 用户进入系统，开始创建服务，默认获取服务模板数据 template done
// --ConnectDatabase
func LoadTemplate() ([]Template, error) {
	// var template map[int]Template
	var template []Template
	db = ConnectDatabase()

	sql := "select * from template"
	rows, err := db.Query(sql)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	for rows.Next() {
		var tempTmplte Template
		errr := rows.Scan(&tempTmplte.Tid, &tempTmplte.ResourceType,
			&tempTmplte.Bandwidth, &tempTmplte.NumberOf,
			&tempTmplte.LastTime, &tempTmplte.ImageId, &tempTmplte.InstanceType)
		if errr != nil {
			return nil, errors.New(err.Error())
		}
		fmt.Println(tempTmplte)
		template = append(template, tempTmplte)
		// template[tempTmplte.Tid] = tempTmplte
	}

	if len(template) != 0 {
		fmt.Println("加载模板成功!:)")
		return template, nil
	} else {
		return nil, errors.New(err.Error())
	}
}

// P00 用户部署已生成的服务，POST请求，操作terraform进行云资源编排，部署成功的服务存储至resource
// --UpdateTfVars
// --UpdateMainTf
// --Command
func LaunchService(uid int, form InstanceConfig, tid int) error {
	/*
		1. 读 / 创建服务
			1.1 创建新tf / 读取已有tf
			1.2 写入配置
		2. 调用terraform命令行工具执行
		3. 根据结果写resource
	*/
	var (
		writer     *bufio.Writer
		cloud_path string
		config     InstanceConfig
	)
	// id := strconv.Itoa(uid)
	// tf path: ./tf/{cloud_path}
	db = ConnectDatabase()
	if tid != -1 {
		// 基于模板逻辑 done
		row := db.QueryRow("select resource_type, bandwidth, number_of, last_time, image_id, instance_type from template where tid = ?", tid)
		err := row.Scan(&config.ResourceType, &config.Bandwidth, &config.NumberOf, &config.LastTime, &config.ImageId, &config.InstanceType)
		if err != nil {
			fmt.Println("query err2: ", err)
			return errors.New(err.Error())
		}
	} else {
		// 自定义逻辑 done
		config.ResourceType = form.ResourceType
		config.Bandwidth = form.Bandwidth
		config.NumberOf = form.NumberOf
		config.LastTime = form.LastTime
		config.ImageId = form.ImageId
		config.InstanceType = form.InstanceType
		config.InstanceName = form.InstanceName
		config.Description = form.Description
		config.TagsName = form.TagsName
	}
	cloud_path = config.ResourceType
	main_tf_file, err := os.OpenFile("./tf/"+cloud_path+"/main.tf", os.O_RDWR, 0666)
	if err != nil {
		// 打开失败：找不到文件，新建相关文件
		err = InitTF(cloud_path, uid)
		if err != nil {
			return errors.New(err.Error())
		}
	}
	main_tf_file, err = os.OpenFile("./tf/"+cloud_path+"/main.tf", os.O_APPEND, 0666)
	// 打开成功，开始编排云资源。调用 terraform cli
	// 根据form中传参写进main.tf
	writer = bufio.NewWriter(main_tf_file)
	content_main := UpdateMainTf(form, tid)
	if _, err = writer.WriteString(content_main); err == nil {
		// fmt.Println("[main]add to buff succ!", content_main)
		if err = writer.Flush(); err != nil {
			return errors.New("[main]writer flush err:" + err.Error())
		} else {
			fmt.Println("[main]writer flush success!")
		}
	}
	// 自动更新当前密钥信息
	if err = UpdateTfVars(uid, cloud_path); err != nil {
		return err
	}
	// 开始调用terraform命令行工具 done
	// --Command
	if e1 := Init(cloud_path); e1 == nil {
		if e2 := Plan(); e2 == nil {
			if e3 := Apply(); e3 == nil {
				fmt.Printf("terraform apply success!")
			} else {
				return errors.New(e3.Error())
			}
		} else {
			return errors.New(e2.Error())
		}
	} else {
		return errors.New(e1.Error())
	}
	// 从resource中返回执行结果
	fmt.Println("aaa ", config.InstanceName)
	fmt.Println("bbb ", config.TagsName)
	r, err := db.Exec("insert into resource(resource_type, bandwidth, number_of, last_time, image_id, instance_type, instance_name) values(?, ?, ?, ?, ?, ?, ?);",
		config.ResourceType, config.Bandwidth, config.NumberOf, config.LastTime, config.ImageId, config.InstanceType, config.InstanceName)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return errors.New(err.Error())
	}
	id, err := r.LastInsertId()
	if err != nil {
		fmt.Println("exec failed, ", err)
		return errors.New(err.Error())
	}
	fmt.Println("\ninsert succ:", id, " 资源添加成功！")
	return nil
}

// P0 创建实例 todo
// --LoadTemplate
// --LaunchService
func CreateNewInstance(uid int) error {
	/*
		1. 使用/不使用模板，加载模板
		2. 部署服务
	*/
	var (
		template map[int]Template // 模板列表
		form     InstanceConfig   // 允许传入的参数
		tid      int              // 使用的模板id
		key      string
	)

	fmt.Println("开始创建实例！:3")
	fmt.Println("请问是否使用已有模板？（或是基于模板定制服务） y/n")
	fmt.Scanln(&key)
	if key == "y" {
		temp, err := LoadTemplate()
		if err != nil {
			return errors.New(err.Error())
		}
		for i := 0; i < len(temp); i++ {
			template[temp[i].Tid] = temp[i]
		}
		// 传入模板 done
		// todo 格式化打印模板信息
		fmt.Println("模板加载成功~", template)
		fmt.Println("请输入模板编号以选择模板：")
		fmt.Scanln(&tid)
		// 传入form
		form.ResourceType = template[tid].ResourceType
		form.NumberOf = strconv.Itoa(template[tid].NumberOf)
		form.LastTime = strconv.Itoa(template[tid].LastTime)
		// form.NumberOf = template[tid].NumberOf
		// form.LastTime = template[tid].LastTime
		fmt.Println("模板已载入~请问是否需要定制？是/否")
		fmt.Scanln(&key)
		if key == "是" {
			// 处理form
		}
	} else {
		// 传入form
		template = nil
		tid = -1
		fmt.Println("需要辛苦自己输入一下参数了噢~('·n·`)")
		fmt.Println("注：如需跳过输入，直接键入回车即可。")
		fmt.Println("请选择云服务器类型：aws/ali/huawei")
		fmt.Scanln(&form.ResourceType)

		fmt.Println("请选择带宽：")
		fmt.Scanln(&form.Bandwidth)

		fmt.Println("请选择实例数量：")
		fmt.Scanln(&form.NumberOf)

		fmt.Println("请输入租期：")
		fmt.Scanln(&form.LastTime)
	}
	// 传入参数准备部署
	fmt.Println("服务部署中……请稍后……><")
	if err := LaunchService(uid, form, tid); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// 封装，用于IO获取密钥 done
func GetKey_io(uid int) AccessKey {
	var key_pair AccessKey
	key_pair.Uid = uid
	fmt.Println("注：如需跳过输入，直接键入回车即可。")
	fmt.Println("请输入aws密钥ID：")
	fmt.Scanln(&key_pair.Aws_a_key)
	fmt.Println("请输入aws私有访问密钥：")
	fmt.Scanln(&key_pair.Aws_s_key)

	fmt.Println("请输入alicloud密钥ID：")
	fmt.Scanln(&key_pair.Ali_a_key)
	fmt.Println("请输入alicloud密钥ID：")
	fmt.Scanln(&key_pair.Ali_s_key)

	fmt.Println("请输入huaweicloud密钥ID：")
	fmt.Scanln(&key_pair.Hw_a_key)
	fmt.Println("请输入huaweicloud密钥ID：")
	fmt.Scanln(&key_pair.Hw_s_key)

	// fmt.Println("获取成功！key_pair:", key_pair)
	return key_pair
}

// 从数据库中读取密钥对 done
func GetKey_db(uid int) (AccessKey, error) {
	var key_pair AccessKey
	db := ConnectDatabase()
	fmt.Println("init key_pair:", key_pair)
	err := db.QueryRow("select uid, aws_a_key, aws_s_key, ali_a_key, ali_s_key, hw_a_key, hw_s_key, aws_region, ali_region, hw_region from accesskey where uid = ?", uid).Scan(&key_pair.Uid, &key_pair.Aws_a_key, &key_pair.Aws_s_key,
		&key_pair.Ali_a_key, &key_pair.Ali_s_key,
		&key_pair.Hw_a_key, &key_pair.Hw_s_key,
		&key_pair.Aws_region, &key_pair.Ali_region, &key_pair.Hw_region,
	)
	if err == sql.ErrNoRows {
		// return AccessKey{}, errors.New(err.Error())
		return AccessKey{}, sql.ErrNoRows
	}
	fmt.Println("get key_pair:", key_pair)
	return key_pair, nil
}

// P0 更新密钥 封装用于响应请求
func UpdateKey(para TfVars, uid int) error {
	fmt.Println("开始更新密钥！:3")
	db = ConnectDatabase()
	r, err := db.Exec("update accesskey set aws_a_key = ?, aws_s_key = ?, ali_a_key = ?, ali_s_key = ?, hw_a_key = ?, hw_s_key = ?, aws_region = ?, ali_region = ?, hw_region = ? where uid = ?;",
		para.Aws_a_key, para.Aws_s_key, para.Ali_a_key, para.Ali_s_key, para.Hw_a_key, para.Hw_s_key,
		para.Aws_region, para.Ali_region, para.Hw_region, uid)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return errors.New(err.Error())
	}
	row_aff, err := r.RowsAffected()
	if err != nil {
		fmt.Println("exec failed, ", err)
		return errors.New(err.Error())
	}
	if row_aff == 0 {
		fmt.Println("没有找到该用户密钥，无法更新！正在创建中……")
		res, err := db.Exec("insert into accesskey(uid, aws_a_key, aws_s_key, ali_a_key, ali_s_key, hw_a_key, hw_s_key, aws_region, ali_region, hw_region) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
			uid, para.Aws_a_key, para.Aws_s_key, para.Ali_a_key, para.Ali_s_key, para.Hw_a_key, para.Hw_s_key,
			para.Aws_region, para.Ali_region, para.Hw_region)
		if err != nil {
			fmt.Println("exec failed, ", err)
			return errors.New(err.Error())
		}
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Println("exec failed, ", err)
			return errors.New(err.Error())
		}
		fmt.Println("insert succ:", id, " 密钥添加成功！")
		return nil
	}
	fmt.Print("update succ, row:", row_aff, " 更新成功！")
	// 写terraform.tfvars.json
	cloud := [3]string{"aws", "alicloud", "hwcloud"}
	for i := 0; i < 3; i++ {
		cloud_path := cloud[i]
		if err := UpdateTfVars(uid, cloud_path); err != nil {
			return err
		}
	}

	return nil
}

// P0 创建实例基于模板 封装用于响应请求 done
func CreateInstanceBaseTpl(tid int, para InstanceConfig) (Resource, error) {
	var resource Resource
	// 传入模板，此时传入form为{}
	fmt.Println("服务部署中……请稍后……><")
	if err := LaunchService(1, para, tid); err != nil {
		return resource, errors.New(err.Error())
	}
	db = ConnectDatabase()
	row := db.QueryRow("select rid, resource_type, bandwidth, number_of, last_time, image_id, instance_type, instance_name, ifnull(region_id, '') from resource where rid=(select max(rid) from resource);")
	err := row.Scan(&resource.Rid, &resource.ResourceType, &resource.Bandwidth, &resource.NumberOf, &resource.LastTime,
		&resource.ImageId, &resource.InstanceType, &resource.InstanceName, &resource.RegionId)
	if err != nil {
		return resource, errors.New(err.Error())
	}
	return resource, nil
}

// P0 创建实例 封装用于响应请求
func CreateInstance(para InstanceConfig) (Resource, error) {
	var resource Resource
	// 传入form，此时传入tid为-1
	fmt.Println("服务部署中……请稍后……><")
	if err := LaunchService(1, para, -1); err != nil {
		return resource, errors.New(err.Error())
	}
	db = ConnectDatabase()
	row := db.QueryRow("select rid, resource_type, bandwidth, number_of, last_time, image_id, instance_type, instance_name, ifnull(region_id, '') from resource where rid=(select max(rid) from resource);")
	err := row.Scan(&resource.Rid, &resource.ResourceType, &resource.Bandwidth, &resource.NumberOf, &resource.LastTime,
		&resource.ImageId, &resource.InstanceType, &resource.InstanceName, &resource.RegionId)
	if err != nil {
		return resource, errors.New(err.Error())
	}

	return resource, nil
}
