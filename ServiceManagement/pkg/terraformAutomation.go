package service

import (
	"bufio"
	_ "encoding/json"
	"errors"
	"fmt"
	_ "io"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitTF(cloud_path string, uid int) error {
	var writer *bufio.Writer
	main_tf_file, err1 := os.Create("./tf/" + cloud_path + "/main.tf")
	var_tf_file, err2 := os.Create("./tf/" + cloud_path + "/variables.tf")
	_, err3 := os.Create("./tf/" + cloud_path + "/terraform.tfvars.json")
	if err1 != nil || err2 != nil || err3 != nil {
		// 创建失败
		return errors.New("create file failed!")
	} else {
		// 创建成功，开始编写默认配置
		// 写入模板字符串
		content_main := ""
		switch {
		case cloud_path == "aws":
			content_main = default_main_aws
			break
		case cloud_path == "alicloud":
			content_main = default_main_ali + "\n" + ali_pre_config
			// 写入阿里前置配置ali_pre_config done
			break
		case cloud_path == "hwcloud":
			content_main = default_main_hw + "\n" + hw_pre_config
			// 写入华为前置配置hw_pre_config done
			break
		default:
			fmt.Println("content_main err!")
		}
		content_var := default_var

		// 写main.tf
		writer = bufio.NewWriter(main_tf_file)
		if _, err := writer.WriteString(content_main); err == nil {
			// fmt.Println("[main]add to buff succ!", content_main)
			if err = writer.Flush(); err != nil {
				return errors.New("[main]writer flush err:" + err.Error())
			} else {
				fmt.Println("[main]writer flush success!")
			}
		}
		// 写variables.tf
		writer = bufio.NewWriter(var_tf_file)
		if _, err := writer.WriteString(content_var); err == nil {
			// fmt.Println("[var]add to buff succ!", content_var)
			if err = writer.Flush(); err != nil {
				return errors.New("[var]writer flush err:" + err.Error())
			} else {
				fmt.Println("[var]writer flush success!")
			}
		}
		// 写terraform.tfvars.json
		if err := UpdateTfVars(uid, cloud_path); err != nil {
			return err
		}
	}
	return nil
}

// 操作main.tf done
func UpdateMainTf(form InstanceConfig, tid int) string {
	db = ConnectDatabase()
	var (
		resource_type string
		image_id      string
		instance_type string
		instance_name string
		// config        InstanceConfig
	)
	// 将form内的字段填充至tf中，若使用模板则对应到tid
	switch {
	case tid != -1:
		// 载入模板
		row := db.QueryRow("select resource_type, image_id, instance_type from template where tid = ?", tid)
		err := row.Scan(&resource_type, &image_id, &instance_type)
		if err != nil {
			fmt.Println("query err1: ", err)
			return ""
		}
		// tpl := db.e
		break
	default:
		// 填充form
		resource_type = form.ResourceType
		image_id = form.ImageId
		instance_type = form.InstanceType
		// 自定义逻辑 done
		// config.ResourceType = form.ResourceType
		// config.Bandwidth = form.Bandwidth
		// config.NumberOf = form.NumberOf
		// config.LastTime = form.LastTime
		// config.ImageId = form.ImageId
		// config.InstanceType = form.InstanceType
		break
	}
	fmt.Println("aaa ", form)
	instance_name = form.InstanceName
	// config.OtherArgs = form.OtherArgs
	// config.RegionId = form.RegionId
	// if form.OtherArgs != (Argument{}) {
	// 	// 传入可选参数 tags 下…… todo
	// 	// name := form.OtherArgs.name
	// 	// desciption := form.OtherArgs.description
	// 	// tags_name := form.OtherArgs.tags_name
	// 	_ = form.OtherArgs.name
	// 	_ = form.OtherArgs.description
	// 	_ = form.OtherArgs.tags_name
	// }
	resource_str := "\nresource"
	fmt.Println("ccc ", form.InstanceName)
	switch {
	case resource_type == "aws":
		resource_str +=
			` "aws_instance" "` + instance_name + `" {
	ami = "` + image_id + `"
	instance_type = "` + instance_type + `"`
		if form.TagsName != "" {
			resource_str +=
				`
	tags = {
		Name = "` + form.TagsName + `"
	}
}`
		} else {
			resource_str +=
				`}`
		}
		break
	case resource_type == "alicloud":
		resource_str +=
			` "alicloud_instance" "` + instance_name + `" {
	instance_name = "` + form.InstanceName + `"
	image_id = "` + form.ImageId + `"
	instance_type = "` + form.InstanceType + `"
	vswitch_id = data.alicloud_instances.myali_instances.instances.0.vswitch_id
  	internet_max_bandwidth_out = 10
	security_groups = data.alicloud_instances.myali_instances.instances.0.security_groups
}`
		break
	case resource_type == "hwcloud":
		resource_str +=
			` "huaweicloud_compute_instance" "` + instance_name + `" {
	name = "` + instance_name + `"
	image_id = "` + image_id + `"
	flavor_id = "` + form.InstanceType + `"
	network {
		uuid = data.huaweicloud_vpc_subnet.net.id
	}`
		if form.TagsName != "" {
			resource_str +=
				`
	tags = {
		Name = "` + form.TagsName + `"
	}
}`
		} else {
			resource_str +=
				`}`
		}
		break
	}
	fmt.Println("resource_str: \n", resource_str)
	return resource_str
}

// 操作variables.tf done
func UpdateTfVars(uid int, resource_type string) error {
	var (
		a_key  string
		s_key  string
		region string
	)
	accesskey, _ := GetKey_db(uid)
	switch {
	case resource_type == "aws":
		a_key = accesskey.Aws_a_key
		s_key = accesskey.Aws_s_key
		region = accesskey.Aws_region
		break
	case resource_type == "alicloud":
		a_key = accesskey.Ali_a_key
		s_key = accesskey.Ali_s_key
		region = accesskey.Ali_region
		break
	case resource_type == "hwcloud":
		a_key = accesskey.Hw_a_key
		s_key = accesskey.Hw_s_key
		region = accesskey.Hw_region
		break
	default:
		return errors.New("[tfvars]writer GetKey_db err: " + resource_type)
	}

	// 解析数据库中读到的accesskey来写入tfvars
	jsonstr :=
		`{
		"access_key": "` + a_key + `",
		"secret_key": "` + s_key + `",
		"region": "` + region + `"
	}`
	// return jsonstr
	// 写入terraform.tfvars.json
	tfvars_json_file, err := os.OpenFile("./tf/"+resource_type+"/terraform.tfvars.json", os.O_RDWR, 0666)
	if err != nil {
		return errors.New("[tfvars]writer OpenFile err:" + err.Error() + "(" + resource_type + ")")
	}
	writer := bufio.NewWriter(tfvars_json_file)
	if _, err = writer.WriteString(jsonstr); err == nil {
		if err = writer.Flush(); err != nil {
			return errors.New("[tfvars]writer flush err:" + err.Error() + "(" + resource_type + ")")
		} else {
			fmt.Println("[tfvars]writer flush success!" + "(" + jsonstr + ")")
		}
	}
	return nil
}
